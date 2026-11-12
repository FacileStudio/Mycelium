package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/FacileStudio/porte"
	jwtpkg "github.com/FacileStudio/porte/oidc/jwt"
	apierrors "github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"
)

const (
	scopeDelegated = "delegated"

	delegatedRead  = "mycelium.read"
	delegatedWrite = "mycelium.write"
)

// delegatedRoutes is the whole of what a delegated token reaches, as
// "METHOD pattern" against the chi route pattern that matched. It is an
// allowlist on purpose: a route added later stays closed to delegated callers
// until someone writes it here. Only the artifact routes reach the common
// tree for a delegated identity, and only for an admin (artifactRoot); every
// other route resolves its tree through scopeRoot, whose commonAllowed refuses
// a delegated caller.
var delegatedRoutes = map[string]bool{
	"GET /api/spaces":            true,
	"POST /api/spaces":           true,
	"POST /api/memory/search":    true,
	"GET /api/sync/tree":         true,
	"GET /api/sync/files/*":      true,
	"PUT /api/sync/files/*":      true,
	"DELETE /api/sync/files/*":   true,
	"GET /api/artifacts":         true,
	"POST /api/artifacts":        true,
	"GET /api/artifacts/{id}":    true,
	"PUT /api/artifacts/{id}":    true,
	"DELETE /api/artifacts/{id}": true,
}

// delegatedVerifier verifies suite tokens minted by Registre's token exchange
// for OIDC_MACHINE_AUDIENCE. Discovery is lazy, as it is for the browser
// login, so the server still boots when the issuer is unreachable.
type delegatedVerifier struct {
	mu       sync.Mutex
	verifier *jwtpkg.Verifier
}

func (s *Server) delegatedTokens(ctx context.Context) (*jwtpkg.Verifier, error) {
	s.delegated.mu.Lock()
	defer s.delegated.mu.Unlock()
	if s.delegated.verifier != nil {
		return s.delegated.verifier, nil
	}
	verifier, err := jwtpkg.New(ctx, jwtpkg.Config{Issuer: s.OIDC.Issuer, Audience: s.OIDC.MachineAudience})
	if err != nil {
		return nil, err
	}
	s.delegated.verifier = verifier
	return verifier, nil
}

// delegationEnabled reports whether a bearer that is not one of Mycelium's own
// tokens may be tried as a delegated suite JWT.
func (s *Server) delegationEnabled(raw string) bool {
	return s.OIDC != nil && s.OIDC.MachineAudience != "" && strings.Count(raw, ".") == 2
}

// authDelegated authenticates a delegated suite JWT and applies its limits:
// the route must be on the allowlist, the method must be covered by the
// token's scopes, and the email claim names the caller. It never yields an
// admin, user or sync scope, so nothing downstream can mistake it for one.
func (s *Server) authDelegated(w http.ResponseWriter, r *http.Request, raw string, adminOnly bool) (Identity, bool) {
	verifier, err := s.delegatedTokens(r.Context())
	if err != nil {
		s.Log.Error("delegated auth: issuer discovery failed", slog.Any("error", err))
		httpjson.WriteError(w, apierrors.Unauthorized("unauthorized"))
		return Identity{}, false
	}
	claims, err := verifier.Verify(r.Context(), raw)
	if err != nil || claims.Subject == "" || strings.TrimSpace(claims.Email) == "" {
		httpjson.WriteError(w, apierrors.Unauthorized("unauthorized"))
		return Identity{}, false
	}
	if adminOnly || !delegatedRouteAllowed(r) || !delegatedScopeAllows(r.Method, claims.Scopes) {
		httpjson.WriteError(w, apierrors.Forbidden("forbidden"))
		return Identity{}, false
	}
	email := strings.TrimSpace(claims.Email)
	if err := s.ensureDelegatedUser(r.Context(), email, claims.Name, s.OIDC.Issuer, claims.Subject); err != nil {
		s.refuseDelegated(w, err)
		return Identity{}, false
	}
	return Identity{Email: email, Scope: scopeDelegated, Delegated: true, Actor: claims.Actor}, true
}

// refuseDelegated answers a delegated call whose account could not be
// ensured: unauthorized for an erased subject, internal when the deletion
// ledger could not be read.
func (s *Server) refuseDelegated(w http.ResponseWriter, err error) {
	if errors.Is(err, porte.ErrSubjectErased) {
		httpjson.WriteError(w, apierrors.Unauthorized("unauthorized"))
		return
	}
	s.Log.Error("delegated auth: account check failed", slog.Any("error", err))
	httpjson.WriteError(w, apierrors.Internal("internal error", err))
}

func delegatedRouteAllowed(r *http.Request) bool {
	rctx := chi.RouteContext(r.Context())
	if rctx == nil {
		return false
	}
	return delegatedRoutes[r.Method+" "+rctx.RoutePattern()]
}

func delegatedScopeAllows(method string, scopes []string) bool {
	if slices.Contains(scopes, delegatedWrite) {
		return true
	}
	return method == http.MethodGet && slices.Contains(scopes, delegatedRead)
}

// ensureDelegatedUser records a caller seen first through a delegated token,
// as the OIDC callback would, except that it never bootstraps an admin: the
// first account on an instance is made by a person signing in, not by a
// service acting for one. The token's subject is stored with the account, so
// Registre's deletion notice for that subject finds it, and an account with
// no picture yet is backfilled with the subject's picture URL instead of
// waiting for a browser sign-in. A subject Registre erased gets no account.
func (s *Server) ensureDelegatedUser(ctx context.Context, email, name, issuer, subject string) error {
	if err := s.refuseErased(ctx, issuer, subject); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	users := s.loadUsers()
	user, ok := users[email]
	if ok && user.Subject != "" && user.AvatarURL != "" {
		return nil
	}
	if !ok {
		user = User{Email: email, Name: name, CreatedAt: time.Now().UTC().Format(time.RFC3339)}
	}
	if user.Subject == "" {
		user.Issuer, user.Subject = issuer, subject
	}
	if user.AvatarURL == "" {
		user.AvatarURL = subjectPictureURL(issuer, subject)
	}
	users[email] = user
	if err := s.saveUsers(users); err != nil {
		s.Log.Error("users: save failed", slog.Any("error", err))
	}
	return nil
}
