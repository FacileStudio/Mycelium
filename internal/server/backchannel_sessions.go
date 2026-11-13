package server

import (
	"context"
	"strings"
	"time"

	"github.com/FacileStudio/porte"
)

// browserSessionPrefix names the token a browser sign-in mints. CLI logins
// (cli:<email>), machine tokens and API tokens carry other names.
const browserSessionPrefix = "session:"

// browserSessions shows porte the browser sessions in tokens.json. porte
// numbers a person by subjectNumber, as mcpIdentities.Find answers, and
// DeleteByUser drops that person's session:<email> tokens. Nothing else is
// reachable from the back-channel route, so every other method refuses.
type browserSessions struct{ s *Server }

// DeleteByUser revokes every browser session of the account numbered userID.
func (b browserSessions) DeleteByUser(_ context.Context, userID int64) (int64, error) {
	s := b.s
	s.mu.Lock()
	defer s.mu.Unlock()
	emails := map[string]bool{}
	for key, user := range s.loadUsers() {
		if subject := storedSubject(s.OIDC.Issuer, user); subject != "" && subjectNumber(s.OIDC.Issuer, subject) == userID {
			emails[strings.ToLower(key)] = true
		}
	}
	var deleted int64
	for hash, token := range s.tokens {
		if strings.HasPrefix(token.Name, browserSessionPrefix) && emails[strings.ToLower(token.UserEmail)] {
			delete(s.tokens, hash)
			deleted++
		}
	}
	if deleted > 0 {
		s.saveTokens()
	}
	return deleted, nil
}

// Create refuses: browser sessions are minted by Mycelium's own callback.
func (browserSessions) Create(context.Context, porte.Session) (porte.Session, error) {
	return porte.Session{}, errBackchannelOnly
}

// Find answers not found: Mycelium authenticates its own tokens.
func (browserSessions) Find(context.Context, string) (porte.Session, error) {
	return porte.Session{}, porte.ErrNotFound
}

// Touch is a no-op: Mycelium records last use itself.
func (browserSessions) Touch(context.Context, string, time.Time) error { return nil }

// Delete refuses: POST /api/auth/logout ends a single session.
func (browserSessions) Delete(context.Context, string) error { return errBackchannelOnly }

// DeleteLogins refuses: Mycelium has no password change to revoke on.
func (browserSessions) DeleteLogins(context.Context, int64, int64) (int64, error) {
	return 0, errBackchannelOnly
}

// ListByUser refuses: the session list is Mycelium's tokens page.
func (browserSessions) ListByUser(context.Context, int64) ([]porte.Session, error) {
	return nil, errBackchannelOnly
}

// DeleteByID refuses for the same reason as ListByUser.
func (browserSessions) DeleteByID(context.Context, int64, int64) error { return errBackchannelOnly }

// DeleteExpired is a no-op: expired tokens are refused at authentication.
func (browserSessions) DeleteExpired(context.Context, time.Time) (int64, error) { return 0, nil }
