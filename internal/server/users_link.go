package server

import (
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/FacileStudio/porte"
)

var (
	errUnverifiedEmail = errors.New("the issuer did not verify this email address")
	errSubjectConflict = errors.New("this account is already linked to another subject of the same issuer")
)

type oidcIdentity struct {
	Issuer        string
	Subject       string
	Email         string
	EmailVerified bool
	Name          string
	Picture       string
	Locale        string
}

func (i oidcIdentity) refresh(user User) User {
	if i.Name != "" {
		user.Name = i.Name
	}
	if i.Locale != "" {
		user.Locale = porte.NormalizeLocale(i.Locale)
	}
	user.AvatarURL = pictureURL(i.Picture)
	return user
}

// link records the issuer and subject on an account that was found through
// them, so a record that only carried the subject's picture holds the
// subject itself from then on.
func (i oidcIdentity) link(user User) User {
	user.Issuer, user.Subject = i.Issuer, i.Subject
	return user
}

// findBySubject finds the account a Registre subject maps to: a stored
// subject first, then the picture URL Registre publishes per subject, which is
// the only mark of the subject on a record provisioned by a delegated call or a
// snapshot before the subject was stored with it.
func findBySubject(users map[string]User, issuer, subject string) (string, bool) {
	for key, user := range users {
		if user.Subject != "" && user.Subject == subject && user.Issuer == issuer {
			return key, true
		}
	}
	picture := subjectPictureURL(issuer, subject)
	for key, user := range users {
		if picture != "" && user.Subject == "" && user.AvatarURL == picture {
			return key, true
		}
	}
	return "", false
}

func findByEmail(users map[string]User, email string) (string, bool) {
	if _, ok := users[email]; ok {
		return email, true
	}
	for key := range users {
		if strings.EqualFold(key, email) {
			return key, true
		}
	}
	return "", false
}

// resolveIdentity applies porte's linking rules to Mycelium's own user store.
// A stored subject always wins. Otherwise the identity links by email, and only
// when the issuer verified that email. An account that already holds another
// subject of the same issuer is never handed to a newcomer, so a recycled
// address cannot take over a live account.
func resolveIdentity(users map[string]User, ident oidcIdentity, now time.Time) (string, User, error) {
	if key, ok := findBySubject(users, ident.Issuer, ident.Subject); ok {
		return key, ident.refresh(ident.link(users[key])), nil
	}
	if !ident.EmailVerified {
		return "", User{}, errUnverifiedEmail
	}
	email := strings.ToLower(strings.TrimSpace(ident.Email))
	if key, ok := findByEmail(users, email); ok {
		user := users[key]
		if user.Subject != "" && user.Issuer == ident.Issuer {
			return "", User{}, errSubjectConflict
		}
		return key, ident.refresh(ident.link(user)), nil
	}
	user := User{
		Email: email, Admin: len(users) == 0, Issuer: ident.Issuer, Subject: ident.Subject,
		CreatedAt: now.UTC().Format(time.RFC3339),
	}
	return email, ident.refresh(user), nil
}

// linkIdentity stores the account an OIDC login resolves to. The first user
// ever seen becomes admin, mirroring the Nuage bootstrap convention.
func (s *Server) linkIdentity(ident oidcIdentity) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	users := s.loadUsers()
	key, user, err := resolveIdentity(users, ident, time.Now())
	if err != nil {
		return User{}, err
	}
	users[key] = user
	if err := s.saveUsers(users); err != nil {
		s.Log.Error("users: save failed", slog.Any("error", err))
	}
	return user, nil
}
