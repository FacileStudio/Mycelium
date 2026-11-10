package server

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/FacileStudio/Mycelium/internal/env"
	"github.com/FacileStudio/porte/directory"
)

const erasedSubject = "46055c5e-acc1-4b60-b3d1-30dff29932bf"

// signNotice signs the user.deleted security event token Registre pushes to
// an app's deletion URI, with the issuer's stub key.
func signNotice(t *testing.T, key *rsa.PrivateKey, issuer, aud, subject string) string {
	t.Helper()
	encode := func(value any) string {
		raw, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		return base64.RawURLEncoding.EncodeToString(raw)
	}
	now := time.Now()
	claims := map[string]any{
		"iss": issuer, "aud": []string{aud}, "sub": subject, "jti": "n1",
		"iat": now.Unix(), "exp": now.Add(10 * time.Minute).Unix(),
		"events": map[string]any{directory.UserDeletedEvent: map[string]int{"erase_within_days": 30}},
	}
	signing := encode(map[string]string{"alg": "RS256", "typ": "secevent+jwt", "kid": "stub"}) + "." + encode(claims)
	digest := crypto.SHA256.New()
	digest.Write([]byte(signing))
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest.Sum(nil))
	if err != nil {
		t.Fatal(err)
	}
	return signing + "." + base64.RawURLEncoding.EncodeToString(signature)
}

func postNotice(h http.Handler, notice string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/api/registre/deletions", strings.NewReader(notice))
	req.Header.Set("Content-Type", "application/secevent+jwt")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

// The account is stored the way a delegated call provisioned it before the
// subject was stored: keyed by email, with only the subject's picture URL.
func seedErasedPerson(t *testing.T, srv *Server, issuer string) string {
	t.Helper()
	email := "gone@facile.studio"
	users := map[string]User{
		email:                {Email: email, AvatarURL: issuer + "/avatars/" + erasedSubject, CreatedAt: "2026-11-10T13:34:47Z"},
		"keep@facile.studio": {Email: "keep@facile.studio", Admin: true, Issuer: issuer, Subject: "keep"},
	}
	if err := srv.saveUsers(users); err != nil {
		t.Fatal(err)
	}
	spaces := map[string]*Space{
		"personal": {ID: "personal", Name: "Bulle", Members: map[string]string{email: roleOwner}},
		"shared":   {ID: "shared", Name: "Team", Members: map[string]string{email: roleAdmin, "keep@facile.studio": roleOwner}, RegistreID: "r1"},
	}
	if err := srv.saveSpaces(spaces); err != nil {
		t.Fatal(err)
	}
	seedFile(t, filepath.Join(srv.spaceDir("personal"), "artifacts", "2026", "11", "2026-11-10-note-abc123.md"), "# Note\n")
	if _, err := srv.mintToken("cli:"+email, scopeUser, email); err != nil {
		t.Fatal(err)
	}
	return email
}

func TestSignedNoticeErasesTheAccountItsSpaceAndItsTokens(t *testing.T) {
	issuer := newSuiteIssuer(t)
	srv := New(t.TempDir(), "")
	srv.OIDC = &env.OIDC{Issuer: issuer.url, ClientID: "mycelium", ClientSecret: "secret", RedirectURL: "https://x/cb"}
	email := seedErasedPerson(t, srv, issuer.url)
	h := srv.Handler()

	if rec := postNotice(h, signNotice(t, issuer.key, issuer.url, "other-app", erasedSubject)); rec.Code != http.StatusBadRequest {
		t.Fatalf("a notice for another audience answered %d", rec.Code)
	}
	for range 2 {
		if rec := postNotice(h, signNotice(t, issuer.key, issuer.url, "mycelium", erasedSubject)); rec.Code != http.StatusAccepted {
			t.Fatalf("the notice answered %d: %s", rec.Code, rec.Body.String())
		}
	}
	assertErased(t, srv, email)
}

func assertErased(t *testing.T, srv *Server, email string) {
	t.Helper()
	users := srv.loadUsers()
	if _, ok := users[email]; ok || len(users) != 1 {
		t.Fatalf("users after the notice: %v", users)
	}
	spaces := srv.loadSpaces()
	if _, ok := spaces["personal"]; ok {
		t.Error("the personal space survived")
	}
	if _, err := os.Stat(srv.spaceDir("personal")); err == nil {
		t.Error("the personal space's files, artifacts included, survived")
	}
	if members := spaces["shared"].Members; len(members) != 1 || members["keep@facile.studio"] != roleOwner {
		t.Errorf("shared space after the notice: %v", members)
	}
	if n := countnamed(srv, "cli:"+email); n != 0 {
		t.Errorf("%d token(s) of the erased person survived", n)
	}
	if _, err := os.Stat(filepath.Join(srv.DataDir, deletionStateFile)); err != nil {
		t.Errorf("the handled subject was not recorded: %v", err)
	}
}
