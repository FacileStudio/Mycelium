package server

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	mcpkit "github.com/FacileStudio/tronc/mcp"
)

// pendingUser is the MCP user id of a confirm page request. Only a browser
// session of an account linked to Registre passes: a session:<email> entry
// that carries an expiry, since no API or machine token is minted with one.
func (s *Server) pendingUser(r *http.Request) (string, error) {
	id := identityFrom(r)
	if id.Delegated || id.Actor != "" {
		return "", mcpkit.ErrDelegated
	}
	s.mu.RLock()
	token, ok := s.tokens[id.TokenHash]
	user, known := s.loadUsers()[id.Email]
	s.mu.RUnlock()
	switch {
	case !ok || !strings.HasPrefix(token.Name, browserSessionPrefix) || !isLoginSession(token) || id.Email == "":
		return "", errBrowserOnly
	case !known || user.Subject == "":
		return "", errNoSubject
	}
	return strconv.FormatInt(subjectNumber(user.Issuer, user.Subject), 10), nil
}

// mcpPendingItem is one action on the confirm page, with what it would
// delete as it is now: a page by its path, or an artifact by its id and
// title.
type mcpPendingItem struct {
	mcpkit.PendingItem
	SpaceID   string `json:"space_id"`
	SpaceName string `json:"space_name"`
	Path      string `json:"path"`
	Title     string `json:"title,omitempty"`
	Log       string `json:"log,omitempty"`
}

// describe adds what an action targets to its item.
func (s *Server) describe(_ context.Context, item mcpkit.PendingItem, a mcpkit.Action) mcpPendingItem {
	out := mcpPendingItem{PendingItem: item}
	target, tree, err := s.actionTree(a)
	if err != nil {
		return out
	}
	out.SpaceID, out.SpaceName, out.Path, out.Log = tree.SpaceID, tree.Name, target.Path, target.Log
	if a.Kind == kindDeleteArtifact {
		if art, err := artifactIn(tree, target.Path); err == nil {
			out.Title = art.Title
		}
	}
	return out
}

// pendingPage is the confirm page of the running MCP server, nil when MCP
// is off. Both confirmed kinds are deletes, so a confirmed run answers
// "deleted".
func (s *Server) pendingPage() *mcpkit.PendingPage[mcpPendingItem] {
	if s.mcp == nil || s.mcp.confirm == nil {
		return nil
	}
	return &mcpkit.PendingPage[mcpPendingItem]{
		Confirmations: s.mcp.confirm, Actor: webActor, User: s.pendingUser, Describe: s.describe,
		Done: func(context.Context, mcpkit.Outcome) mcpkit.PendingResult {
			return mcpkit.PendingResult{Status: "deleted"}
		},
	}
}

// authed is auth as chi middleware, for routers mounted by a kit.
func (s *Server) authed(next http.Handler) http.Handler {
	return s.auth(false, next.ServeHTTP)
}
