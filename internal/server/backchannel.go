package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/oidc"
	"github.com/FacileStudio/porte/session"
	apierrors "github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"
	"github.com/go-chi/chi/v5"
)

// errBackchannelOnly is what porte gets if the back-channel kit ever reaches
// for a store it was not built to use: Mycelium signs people in through its
// own handlers, and this kit only serves the logout route.
var errBackchannelOnly = errors.New("mycelium: the porte kit only serves back-channel logout")

type backchannelRuntime struct {
	mu      sync.Mutex
	handler http.Handler
}

// backchannelLogout hands Registre's logout token to porte's own handler,
// which verifies it and revokes the subject's browser sessions here.
func (s *Server) backchannelLogout(w http.ResponseWriter, r *http.Request) {
	if s.OIDC == nil {
		httpjson.WriteError(w, apierrors.NotFound("no such endpoint"))
		return
	}
	handler, err := s.backchannelHandler(r.Context())
	if err != nil {
		s.Log.Error("backchannel: the identity provider is unreachable", slog.Any("error", err))
		httpjson.WriteError(w, apierrors.Unavailable("the identity provider is unreachable"))
		return
	}
	handler.ServeHTTP(w, r)
}

// backchannelHandler builds porte's OIDC kit on first use, as ensureOIDC
// does, so the server boots while Registre is down and a later logout
// retries discovery. Only the back-channel route of the kit is reachable.
func (s *Server) backchannelHandler(ctx context.Context) (http.Handler, error) {
	s.backchannel.mu.Lock()
	defer s.backchannel.mu.Unlock()
	if s.backchannel.handler != nil {
		return s.backchannel.handler, nil
	}
	cfg := s.backchannelConfig()
	manager, err := session.New(cfg, session.Deps{Sessions: browserSessions{s}, Logger: s.Log})
	if err != nil {
		return nil, err
	}
	discovery, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	kit, err := oidc.New(discovery, cfg, oidc.Deps{
		Users:      backchannelUsers{},
		Identities: mcpIdentities{s},
		Codes:      backchannelCodes{},
		Sessions:   manager,
		Logger:     s.Log,
	})
	if err != nil {
		return nil, err
	}
	kitRoutes := chi.NewRouter()
	kit.Mount(kitRoutes)
	only := chi.NewRouter()
	only.Post(porte.RouteBackchannelLogout, kitRoutes.ServeHTTP)
	s.backchannel.handler = only
	return only, nil
}

// backchannelConfig is the porte view of Mycelium's OIDC environment. The
// success URL is only read by the sign-in routes this kit never serves, so
// an unset one falls back to the redirect URL to pass porte's validation.
func (s *Server) backchannelConfig() porte.Config {
	success := s.OIDC.SuccessURL
	if success == "" {
		success = s.OIDC.RedirectURL
	}
	return porte.Config{
		Issuer:       s.OIDC.Issuer,
		ClientID:     s.OIDC.ClientID,
		ClientSecret: s.OIDC.ClientSecret,
		RedirectURL:  s.OIDC.RedirectURL,
		SuccessURL:   success,
	}
}

// backchannelUsers refuses: accounts are created by linkIdentity only.
type backchannelUsers struct{}

// UpsertFromOIDC refuses.
func (backchannelUsers) UpsertFromOIDC(context.Context, porte.Claims) (int64, error) {
	return 0, errBackchannelOnly
}

// backchannelCodes refuses: CLI login codes live in oidc_login_code.go.
type backchannelCodes struct{}

// Create refuses.
func (backchannelCodes) Create(context.Context, porte.LoginCode) error { return errBackchannelOnly }

// Consume refuses.
func (backchannelCodes) Consume(context.Context, string) (porte.LoginCode, error) {
	return porte.LoginCode{}, errBackchannelOnly
}

// DeleteExpired is a no-op.
func (backchannelCodes) DeleteExpired(context.Context, time.Time) (int64, error) { return 0, nil }
