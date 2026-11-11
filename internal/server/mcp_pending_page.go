package server

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	apierrors "github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"
	mcpkit "github.com/FacileStudio/tronc/mcp"
)

// pendingUser is the MCP user id of a confirm page request. Only a browser
// session of an account linked to Registre passes.
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
	case !ok || !strings.HasPrefix(token.Name, "session:") || id.Email == "":
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

// pendingPage is the confirm page of the running MCP server.
func (s *Server) pendingPage() (*mcpkit.PendingPage[mcpPendingItem], error) {
	if s.mcp == nil || s.mcp.confirm == nil {
		return nil, apierrors.NotFound("MCP is off on this server")
	}
	return &mcpkit.PendingPage[mcpPendingItem]{
		Confirmations: s.mcp.confirm, Actor: webActor, User: s.pendingUser, Describe: s.describe,
	}, nil
}

// pendingList answers the signed-in person's waiting actions.
func (s *Server) pendingList(w http.ResponseWriter, r *http.Request) {
	page, err := s.pendingPage()
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	userID, err := s.pendingUser(r)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	list, err := page.List(r.Context(), userID)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, list)
}

// pendingConfirm runs one waiting action once.
func (s *Server) pendingConfirm(w http.ResponseWriter, r *http.Request) {
	s.pendingResolve(w, r, func(c *mcpkit.Confirmations, who mcpkit.Caller, id string) (string, error) {
		_, err := c.Confirm(r.Context(), who, id)
		return "deleted", err
	})
}

// pendingDismiss closes one waiting action without running it.
func (s *Server) pendingDismiss(w http.ResponseWriter, r *http.Request) {
	s.pendingResolve(w, r, func(c *mcpkit.Confirmations, who mcpkit.Caller, id string) (string, error) {
		return mcpkit.OutcomeDismissed, c.Dismiss(r.Context(), who, id)
	})
}

// pendingResolve checks the person and hands the action id to resolve.
func (s *Server) pendingResolve(w http.ResponseWriter, r *http.Request, resolve func(*mcpkit.Confirmations, mcpkit.Caller, string) (string, error)) {
	if _, err := s.pendingPage(); err != nil {
		httpjson.WriteError(w, err)
		return
	}
	userID, err := s.pendingUser(r)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	status, err := resolve(s.mcp.confirm, mcpkit.Caller{UserID: userID, Actor: webActor}, pathParam(r, "id"))
	var known *apierrors.Error
	if err != nil && !errors.As(err, &known) {
		err = apierrors.Internal("internal error", err)
	}
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, mcpkit.PendingResult{Status: status})
}
