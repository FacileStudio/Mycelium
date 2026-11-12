package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/FacileStudio/porte"
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
	return writeFileAtomic(f.path, state, 0o600)
}

func (s *Server) deletions() *directory.Deletions {
	s.erasure.once.Do(func() {
		receiver := &directory.Deletions{
			OnPersonDeleted: s.erasePerson,
			Holds:           s.holdsIdentity,
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

// holdsIdentity reports whether an account in .users.json maps to the
// subject again, found exactly as erasePerson finds it, so the receiver
// erases an account re-created after its notice was handled.
func (s *Server) holdsIdentity(_ context.Context, issuer, subject string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := findBySubject(s.loadUsers(), issuer, subject)
	return ok, nil
}

// refuseErased fails with porte.ErrSubjectErased when Registre erased the
// subject and this instance handled it. Every path that creates or links an
// account asks it first, since none of them goes through porte/oidc. An
// unreadable ledger is an error, so nothing is created on doubt.
func (s *Server) refuseErased(ctx context.Context, issuer, subject string) error {
	erased, err := s.deletions().Erased(ctx, issuer, subject)
	if err != nil {
		return err
	}
	if erased {
		return porte.ErrSubjectErased
	}
	return nil
}

func (s *Server) deletionsHandler(w http.ResponseWriter, r *http.Request) {
	s.deletions().Handler().ServeHTTP(w, r)
}

// erasePerson deletes the account behind a Registre identity and what it owns:
// its pending MCP actions, its sessions and tokens, its memberships, and the spaces only it used and
// Registre does not own. Space data shared with others stays. A subject with no
// account is a no-op, so a replayed notice changes nothing.
func (s *Server) erasePerson(ctx context.Context, issuer, subject string) error {
	if s.mcp != nil {
		if err := s.mcp.pending.Erase(ctx, strconv.FormatInt(subjectNumber(issuer, subject), 10)); err != nil {
			return err
		}
	}
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
		if best == "" || outranks(role, email, space.Members[best], best) {
			best = email
		}
	}
	if best != "" && space.RegistreID == "" {
		space.Members[best] = roleOwner
	}
}

func outranks(role, email, bestRole, best string) bool {
	if (role == roleAdmin) != (bestRole == roleAdmin) {
		return role == roleAdmin
	}
	return email < best
}
