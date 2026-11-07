package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/FacileStudio/porte/directory"
)

const deletionStateFile = ".deletions.json"

type deletionRuntime struct {
	once     sync.Once
	receiver *directory.Deletions
}

type fileDeletionState struct {
	path string
}

func (f fileDeletionState) Load(context.Context) ([]byte, error) {
	data, err := os.ReadFile(f.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	return data, err
}

func (f fileDeletionState) Save(_ context.Context, state []byte) error {
	if err := os.MkdirAll(filepath.Dir(f.path), 0o755); err != nil {
		return err
	}
	tmp := f.path + ".tmp"
	if err := os.WriteFile(tmp, state, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, f.path)
}

func (s *Server) deletions() *directory.Deletions {
	s.erasure.once.Do(func() {
		receiver := &directory.Deletions{
			OnPersonDeleted: s.erasePerson,
			State:           fileDeletionState{path: filepath.Join(s.DataDir, deletionStateFile)},
			Logger:          s.Log,
		}
		if s.OIDC != nil {
			receiver.Client = s.directorySyncer().Client
		}
		s.erasure.receiver = receiver
	})
	return s.erasure.receiver
}

// RunDeletions polls Registre for erased accounts until ctx ends. Mycelium
// runs as one container, so this is the single replica that may call it.
func (s *Server) RunDeletions(ctx context.Context) {
	if s.OIDC == nil {
		return
	}
	s.deletions().Run(ctx)
}

func (s *Server) deletionsHandler(w http.ResponseWriter, r *http.Request) {
	s.deletions().Handler().ServeHTTP(w, r)
}

// erasePerson deletes the account behind a Registre identity and what it owns:
// its sessions and tokens, its memberships, and the spaces only it used and
// Registre does not own. Space data shared with others stays. A subject with no
// account is a no-op, so a replayed notice changes nothing.
func (s *Server) erasePerson(_ context.Context, issuer, subject string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	users := s.loadUsers()
	email, ok := findBySubject(users, issuer, subject)
	if !ok {
		return nil
	}
	if err := s.dropMemberships(email); err != nil {
		return err
	}
	for hash, token := range s.tokens {
		if token.UserEmail == email {
			delete(s.tokens, hash)
		}
	}
	s.saveTokens()
	delete(users, email)
	if err := s.saveUsers(users); err != nil {
		return err
	}
	s.Log.Info("registre: erased a deleted person's account")
	return nil
}

func (s *Server) dropMemberships(email string) error {
	spaces := s.loadSpaces()
	changed := false
	for id, space := range spaces {
		if _, member := space.Members[email]; !member {
			continue
		}
		changed = true
		delete(space.Members, email)
		if len(space.Members) == 0 && space.RegistreID == "" {
			delete(spaces, id)
			if err := os.RemoveAll(s.spaceDir(id)); err != nil {
				s.Log.Error("registre: could not remove an erased person's space", slog.Any("error", err))
				return err
			}
			continue
		}
		promoteOwner(space)
	}
	if !changed {
		return nil
	}
	return s.saveSpaces(spaces)
}

func promoteOwner(space *Space) {
	best := ""
	for email, role := range space.Members {
		if role == roleOwner {
			return
		}
		if best == "" || rank(role) > rank(space.Members[best]) || (rank(role) == rank(space.Members[best]) && email < best) {
			best = email
		}
	}
	if best != "" && space.RegistreID == "" {
		space.Members[best] = roleOwner
	}
}

func rank(role string) int {
	if role == roleAdmin {
		return 1
	}
	return 0
}
