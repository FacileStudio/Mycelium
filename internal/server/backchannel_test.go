package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/FacileStudio/Mycelium/internal/env"
	"github.com/FacileStudio/porte"
)

func logoutToken(t *testing.T, issuer suiteIssuer, aud, subject string) string {
	t.Helper()
	now := time.Now()
	return signJWT(t, issuer.key, map[string]any{
		"iss": issuer.url, "aud": aud, "sub": subject, "jti": subject + "-logout",
		"iat": now.Unix(), "exp": now.Add(time.Minute).Unix(),
		"events": map[string]any{"http://schemas.openid.net/event/backchannel-logout": map[string]any{}},
	})
}

func postLogout(h http.Handler, token string) *httptest.ResponseRecorder {
	form := url.Values{"logout_token": {token}}
	req := httptest.NewRequest(http.MethodPost, "/api"+porte.RouteBackchannelLogout, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func meStatus(h http.Handler, token string) int {
	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code
}

// seedSignedIn stores alice and bob with their subject, carol with only
// Registre's picture URL, and returns a token for each login it names.
func seedSignedIn(t *testing.T, srv *Server, issuer string) map[string]string {
	t.Helper()
	users := map[string]User{
		"alice@facile.studio": {Email: "alice@facile.studio", Admin: true, Issuer: issuer, Subject: "alice"},
		"bob@facile.studio":   {Email: "bob@facile.studio", Issuer: issuer, Subject: "bob"},
		"carol@facile.studio": {Email: "carol@facile.studio", AvatarURL: subjectPictureURL(issuer, "carol")},
	}
	if err := srv.saveUsers(users); err != nil {
		t.Fatal(err)
	}
	tokens := map[string]string{}
	for _, name := range []string{"session:alice@facile.studio", "cli:alice@facile.studio", "session:bob@facile.studio", "session:carol@facile.studio"} {
		token, err := srv.mintNamedSession(name, name[strings.Index(name, ":")+1:], scopeUser)
		if err != nil {
			t.Fatal(err)
		}
		tokens[name] = token
	}
	return tokens
}

func expectMe(t *testing.T, h http.Handler, tokens map[string]string, want map[string]int) {
	t.Helper()
	for name, status := range want {
		if code := meStatus(h, tokens[name]); code != status {
			t.Errorf("%s answers %d, want %d", name, code, status)
		}
	}
}

// Signing out of Registre signs the person out of Mycelium: porte verifies
// the logout token and every browser session of that subject goes, whether
// the account stores its subject or only Registre's picture URL. Somebody
// else's sessions, a token for another audience and the person's CLI login
// are left alone.
func TestBackchannelLogoutRevokesOnlyTheSubjectsBrowserSessions(t *testing.T) {
	issuer := newSuiteIssuer(t)
	srv := New(t.TempDir(), "")
	srv.OIDC = &env.OIDC{Issuer: issuer.url, ClientID: "mycelium", ClientSecret: "secret", RedirectURL: "https://x/cb"}
	tokens := seedSignedIn(t, srv, issuer.url)
	h := srv.Handler()
	before := map[string]int{}
	for name := range tokens {
		before[name] = http.StatusOK
	}
	expectMe(t, h, tokens, before)

	if rec := postLogout(h, logoutToken(t, issuer, "another-app", "bob")); rec.Code != http.StatusUnauthorized {
		t.Fatalf("a logout token for another audience answered %d: %s", rec.Code, rec.Body)
	}
	for _, subject := range []string{"alice", "carol"} {
		if rec := postLogout(h, logoutToken(t, issuer, "mycelium", subject)); rec.Code != http.StatusOK {
			t.Fatalf("the logout of %s answered %d: %s", subject, rec.Code, rec.Body)
		}
	}

	expectMe(t, h, tokens, map[string]int{
		"session:alice@facile.studio": http.StatusUnauthorized,
		"session:carol@facile.studio": http.StatusUnauthorized,
		"session:bob@facile.studio":   http.StatusOK,
		"cli:alice@facile.studio":     http.StatusOK,
	})
}

func TestBackchannelLogoutIsAbsentWithoutSSO(t *testing.T) {
	srv := New(t.TempDir(), "")
	if rec := postLogout(srv.Handler(), "x"); rec.Code != http.StatusNotFound {
		t.Fatalf("without SSO the route answered %d, want 404", rec.Code)
	}
}
