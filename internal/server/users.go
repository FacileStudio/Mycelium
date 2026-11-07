package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	apierrors "github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"
)

// User is an account on the server.
type User struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	Admin     bool   `json:"admin"`
	AvatarURL string `json:"avatar_url"`
	Locale    string `json:"locale,omitempty"`
	Issuer    string `json:"issuer,omitempty"`
	Subject   string `json:"subject,omitempty"`
	CreatedAt string `json:"created_at"`
}

func (s *Server) usersPath() string {
	return filepath.Join(s.DataDir, ".users.json")
}

func (s *Server) loadUsers() map[string]User {
	users, err := readStore[User](s.usersPath())
	if err != nil {
		s.Log.Error("users: corrupt store", slog.String("path", s.usersPath()), slog.Any("error", err))
	}
	return users
}

func (s *Server) saveUsers(users map[string]User) error {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(s.DataDir, 0o755); err != nil {
		return err
	}
	return writeFileAtomic(s.usersPath(), data, 0o600)
}

func (s *Server) usersList(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r)
	if id.Email == "" && id.Scope != scopeAdmin {
		httpjson.WriteError(w, apierrors.Forbidden("forbidden"))
		return
	}
	s.mu.RLock()
	users := s.loadUsers()
	s.mu.RUnlock()
	out := make([]User, 0, len(users))
	for _, u := range users {
		out = append(out, u)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Email < out[j].Email })
	httpjson.WriteJSON(w, http.StatusOK, out)
}

func (s *Server) userUpdate(w http.ResponseWriter, r *http.Request) {
	email := pathParam(r, "email")
	if email == "" {
		httpjson.WriteError(w, apierrors.Invalid("email required"))
		return
	}

	var req struct {
		Name  *string `json:"name"`
		Admin *bool   `json:"admin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.WriteError(w, apierrors.Invalid("bad request"))
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	users := s.loadUsers()
	user, ok := users[email]
	if !ok {
		httpjson.WriteError(w, apierrors.NotFound("user not found"))
		return
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Admin != nil {
		user.Admin = *req.Admin
	}
	users[email] = user
	if err := s.saveUsers(users); err != nil {
		s.Log.Error("users: save failed", slog.Any("error", err))
		httpjson.WriteError(w, apierrors.Internal("could not save user", err))
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, user)
}

func (s *Server) authMe(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r)
	if id.Email != "" {
		s.mu.RLock()
		user, ok := s.loadUsers()[id.Email]
		s.mu.RUnlock()
		if ok {
			httpjson.WriteJSON(w, http.StatusOK, map[string]any{"email": user.Email, "name": user.Name, "admin": user.Admin, "avatar_url": user.AvatarURL, "locale": user.Locale})
			return
		}
		httpjson.WriteJSON(w, http.StatusOK, map[string]any{"email": id.Email, "name": "", "admin": id.Scope == scopeAdmin})
		return
	}
	if id.Scope == scopeAdmin || id.Scope == scopeSync {
		httpjson.WriteJSON(w, http.StatusOK, map[string]any{"email": "", "name": "machine", "admin": id.Scope == scopeAdmin})
		return
	}
	httpjson.WriteError(w, apierrors.Forbidden("forbidden"))
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r)
	if id.TokenHash != "" {
		s.mu.Lock()
		delete(s.tokens, id.TokenHash)
		s.saveTokens()
		s.mu.Unlock()
	}
	w.WriteHeader(http.StatusNoContent)
}
