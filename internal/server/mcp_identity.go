package server

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/oidc"
	jwtpkg "github.com/FacileStudio/porte/oidc/jwt"
)

// errMCPIdentityWrite is what porte gets if it ever tries to write through
// the MCP identity view: Mycelium's accounts live in .users.json, and only
// linkIdentity writes them.
var errMCPIdentityWrite = errors.New("mycelium: MCP identities are read through .users.json only")

// errMCPNoOwner refuses to provision the first account of an instance: the
// first sign-in makes its admin, and that must be a person in a browser.
var errMCPNoOwner = errors.New("this server has no account yet: its owner signs in to the web app first")

// mcpIdentities shows porte/mcpauth Mycelium's accounts as porte identities.
// An account is found by its Registre subject, exactly as a deletion notice
// finds it, and gets a number derived from that subject, since .users.json
// keys accounts by email and has no numeric id.
type mcpIdentities struct{ s *Server }

// Find resolves a subject to the account it is linked to, or ErrNotFound.
func (m mcpIdentities) Find(_ context.Context, provider, subject string) (porte.StoredIdentity, error) {
	m.s.mu.RLock()
	_, ok := findBySubject(m.s.loadUsers(), provider, subject)
	m.s.mu.RUnlock()
	if !ok {
		return porte.StoredIdentity{}, porte.ErrNotFound
	}
	return porte.StoredIdentity{UserID: subjectNumber(provider, subject), Provider: provider, Subject: subject}, nil
}

// Save refuses: mcpauth only reads identities.
func (mcpIdentities) Save(context.Context, porte.StoredIdentity) error { return errMCPIdentityWrite }

// MarkRolesSynced refuses: Mycelium caches no roles.
func (mcpIdentities) MarkRolesSynced(context.Context, string, string, time.Time) error {
	return errMCPIdentityWrite
}

// ListByUser refuses: nothing on the MCP path lists a person's identities.
func (mcpIdentities) ListByUser(context.Context, int64) ([]porte.StoredIdentity, error) {
	return nil, errMCPIdentityWrite
}

// subjectNumber is a stable, positive, non-zero number for a subject. It only
// keys the per-person rate limit and tells porte the account exists: tools
// find the account again by the subject itself.
func subjectNumber(issuer, subject string) int64 {
	sum := sha256.Sum256([]byte(issuer + "\x00" + subject))
	return int64(binary.BigEndian.Uint64(sum[:8])>>1) | 1
}

// mcpAccounts creates the account on a person's first MCP call, through the
// same linking rules as a browser sign-in, so the person need not open the
// web app first. Unlike a sign-in it never bootstraps the first admin: an
// instance's first account is made by a person in a browser, so every
// account it creates is a plain user.
type mcpAccounts struct{ s *Server }

// Config names the issuer the accounts are filed under.
func (m mcpAccounts) Config() porte.Config { return porte.Config{Issuer: m.s.OIDC.Issuer} }

// ProvisionBearer links the token's verified identity to an account.
func (m mcpAccounts) ProvisionBearer(ctx context.Context, claims jwtpkg.Claims) (porte.Identity, error) {
	verified := claims.Profile.EmailVerified != nil && *claims.Profile.EmailVerified
	if claims.Subject == "" || strings.TrimSpace(claims.Email) == "" || !verified {
		return porte.Identity{}, fmt.Errorf("%w: the token carries no verified email", oidc.ErrProvisionRefused)
	}
	ident := oidcIdentity{
		Issuer: m.s.OIDC.Issuer, Subject: claims.Subject, Email: claims.Email, EmailVerified: true,
		Name: claims.Name, Picture: claims.Profile.Picture, Locale: claims.Profile.Locale,
	}
	if err := m.s.provisionIdentity(ctx, ident); err != nil {
		return porte.Identity{}, fmt.Errorf("%w: %w", oidc.ErrProvisionRefused, err)
	}
	return porte.Identity{UserID: subjectNumber(ident.Issuer, ident.Subject), Subject: claims.Subject}, nil
}

// provisionIdentity stores the account an MCP token resolves to, as
// linkIdentity does for a sign-in, except that a new account is never admin
// and an instance without accounts is left to its owner's first sign-in. A
// subject Registre erased is refused, so a token minted before the erasure
// cannot bring the account back.
func (s *Server) provisionIdentity(ctx context.Context, ident oidcIdentity) error {
	if err := s.refuseErased(ctx, ident.Issuer, ident.Subject); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	users := s.loadUsers()
	if len(users) == 0 {
		return errMCPNoOwner
	}
	key, user, err := resolveIdentity(users, ident, time.Now())
	if err != nil {
		return err
	}
	users[key] = user
	return s.saveUsers(users)
}
