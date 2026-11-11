package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"github.com/FacileStudio/porte/mcpauth"
	mcpkit "github.com/FacileStudio/tronc/mcp"
	"github.com/modelcontextprotocol/go-sdk/auth"
)

// mcpScopePrefix names Mycelium's MCP scopes: mycelium.read and
// mycelium.write, the scopes delegated suite tokens already carry.
const mcpScopePrefix = "mycelium"

// errMCPVerifierDown is what a caller gets while Registre's keys cannot be
// fetched. It is not a verdict on the token.
var errMCPVerifierDown = errors.New("token verification is unavailable")

// mcpRuntime is the remote MCP server, nil when MCP_RESOURCE is unset.
type mcpRuntime struct {
	server   *mcpkit.Server
	verifier lazyVerifier
	resource string
}

// lazyVerifier builds porte's MCP verifier on first use and retries until
// Registre answers, as the browser login does, so the sync server still
// boots while the issuer is unreachable.
type lazyVerifier struct {
	mu     sync.Mutex
	build  func(context.Context) (auth.TokenVerifier, error)
	ready  auth.TokenVerifier
	logger *slog.Logger
}

// verify checks a token with the verifier, building it first if need be.
func (l *lazyVerifier) verify(ctx context.Context, token string, r *http.Request) (*auth.TokenInfo, error) {
	l.mu.Lock()
	if l.ready == nil {
		verifier, err := l.build(ctx)
		if err != nil {
			l.mu.Unlock()
			l.logger.Error("mcp: issuer discovery failed", slog.Any("error", err))
			return nil, errMCPVerifierDown
		}
		l.ready = verifier
	}
	verifier := l.ready
	l.mu.Unlock()
	return verifier(ctx, token, r)
}

// EnableMCP builds the remote MCP server at OIDC.MCPResource, and does
// nothing when OIDC or the resource is unset. Call it before Handler. A
// person's first MCP call creates their account, as a first sign-in would.
func (s *Server) EnableMCP() error {
	if s.OIDC == nil || s.OIDC.MCPResource == "" {
		return nil
	}
	rt := &mcpRuntime{resource: s.OIDC.MCPResource}
	rt.verifier = lazyVerifier{logger: s.Log, build: func(ctx context.Context) (auth.TokenVerifier, error) {
		return mcpauth.Verifier(ctx, mcpauth.Config{
			Issuer: s.OIDC.Issuer, Resource: s.OIDC.MCPResource, ScopePrefix: mcpScopePrefix,
			Identities: mcpIdentities{s: s}, Provision: true, Accounts: mcpAccounts{s: s}, Logger: s.Log,
		})
	}}
	tools := &mcpTools{s: s}
	app, err := mcpkit.NewApp(mcpkit.AppConfig{
		Resource: s.OIDC.MCPResource, Issuer: s.OIDC.Issuer, Version: s.Version,
		Origins: s.CORSAllowedOrigins, Verifier: rt.verifier.verify, Logger: s.Log,
	}, mcpkit.AppSpec{Name: mcpScopePrefix, Title: "Mycelium", Instructions: mcpInstructions, Tools: tools.list()})
	if err != nil {
		return err
	}
	tools.appURL, rt.server = app.URL, app.Server
	s.mcp = rt
	s.Log.Info("mcp server on", slog.String("resource", rt.resource))
	return nil
}
