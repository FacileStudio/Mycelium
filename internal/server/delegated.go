package server

import (
	"context"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

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
// until someone writes it here. None of these routes touches the common tree
// for a delegated identity, because commonAllowed refuses it and every one of
// them resolves its tree through scopeRoot.
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
	s.ensureDelegatedUser(email, claims.Name, subjectPictureURL(s.OIDC.Issuer, claims.Subject))
	return Identity{Email: email, Scope: scopeDelegated, Delegated: true, Actor: claims.Actor}, true
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
// service acting for one. It is also the only place a stored account meets
// its Registre subject, so an account with no picture yet is backfilled here
// with the subject's picture URL instead of waiting for a browser sign-in.
func (s *Server) ensureDelegatedUser(email, name, picture string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	users := s.loadUsers()
	user, ok := users[email]
	if ok && (user.AvatarURL != "" || picture == "") {
		return
	}
	if !ok {
		user = User{Email: email, Name: name, CreatedAt: time.Now().UTC().Format(time.RFC3339)}
	}
	user.AvatarURL = picture
	users[email] = user
	if err := s.saveUsers(users); err != nil {
		s.Log.Error("users: save failed", slog.Any("error", err))
	}
}
