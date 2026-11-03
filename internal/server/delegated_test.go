package server

import (
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
)

const (
	delegatedEmail = "yann@facile.studio"
	otherEmail     = "noah@facile.studio"
)

// suiteIssuer is a stand-in for Registre: discovery, a JWKS, and the key it
// signs exchanged tokens with.
type suiteIssuer struct {
	url string
	key *rsa.PrivateKey
}

func newSuiteIssuer(t *testing.T) suiteIssuer {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"issuer": server.URL, "jwks_uri": server.URL + "/jwks"})
	})
	jwk := map[string]string{
		"kty": "RSA", "alg": "RS256", "use": "sig", "kid": "stub",
		"n": base64.RawURLEncoding.EncodeToString(key.N.Bytes()), "e": "AQAB",
	}
	mux.HandleFunc("/jwks", func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"keys": []map[string]string{jwk}})
	})
	return suiteIssuer{url: server.URL, key: key}
}

func (i suiteIssuer) token(t *testing.T, aud, scope string) string {
	t.Helper()
	now := time.Now()
	return signJWT(t, i.key, map[string]any{
		"iss": i.url, "aud": aud, "sub": "registre-subject", "email": delegatedEmail,
		"scope": scope, "act": map[string]string{"sub": "bulle"},
		"iat": now.Unix(), "exp": now.Add(5 * time.Minute).Unix(),
	})
}

type delegatedFixture struct {
	srv     *Server
	h       http.Handler
	issuer  suiteIssuer
	rw      string
	ro      string
	own     string
	foreign string
}

func newDelegatedFixture(t *testing.T, audience string) delegatedFixture {
	t.Helper()
	issuer := newSuiteIssuer(t)
	srv := New(t.TempDir(), "")
	srv.OIDC = &env.OIDC{Issuer: issuer.url, ClientID: "web", MachineAudience: audience}
	h := srv.Handler()
	owner := sessionFor(t, srv, delegatedEmail, true)
	other := sessionFor(t, srv, otherEmail, false)
	f := delegatedFixture{
		srv: srv, h: h, issuer: issuer,
		rw:      issuer.token(t, "mycelium", "mycelium.read mycelium.write"),
		ro:      issuer.token(t, "mycelium", "mycelium.read"),
		own:     createSpace(t, h, owner, "Own"),
		foreign: createSpace(t, h, other, "Other"),
	}
	seedFile(t, filepath.Join(srv.DataDir, "memory", "secret.md"), "# Owner\n\nzanzibar common secret\n")
	seedFile(t, filepath.Join(srv.spaceDir(f.own), "memory", "note.md"), "# Note\n\nzanzibar space note\n")
	return f
}

func seedFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func expect(t *testing.T, f delegatedFixture, c apiCall, want int) *httptest.ResponseRecorder {
	t.Helper()
	rec := spReq(t, f.h, c)
	if rec.Code != want {
		t.Fatalf("%s %s: got %d, want %d: %s", c.Method, c.Path, rec.Code, want, rec.Body.String())
	}
	return rec
}

func TestDelegatedReadsOwnSpace(t *testing.T) {
	f := newDelegatedFixture(t, "mycelium")
	rec := expect(t, f, apiCall{Method: "GET", Path: "/api/sync/files/memory/note.md?space_id=" + f.own, Token: f.ro}, http.StatusOK)
	if !strings.Contains(rec.Body.String(), "space note") {
		t.Fatalf("wrong file served: %s", rec.Body.String())
	}
	expect(t, f, apiCall{Method: "GET", Path: "/api/sync/tree?space_id=" + f.own, Token: f.ro}, http.StatusOK)
	rec = expect(t, f, apiCall{Method: "GET", Path: "/api/spaces", Token: f.ro}, http.StatusOK)
	if !strings.Contains(rec.Body.String(), f.own) || strings.Contains(rec.Body.String(), f.foreign) {
		t.Fatalf("space list must hold only the caller's spaces: %s", rec.Body.String())
	}
}

func TestDelegatedWritesOwnSpace(t *testing.T) {
	f := newDelegatedFixture(t, "mycelium")
	expect(t, f, apiCall{Method: "PUT", Path: "/api/sync/files/memory/new.md?space_id=" + f.own, Token: f.rw, Body: "hello"}, http.StatusNoContent)
	expect(t, f, apiCall{Method: "DELETE", Path: "/api/sync/files/memory/new.md?space_id=" + f.own, Token: f.rw}, http.StatusNoContent)
	rec := expect(t, f, apiCall{Method: "POST", Path: "/api/spaces", Token: f.rw, Body: `{"name":"Bulle"}`}, http.StatusCreated)
	if !strings.Contains(rec.Body.String(), `"role":"owner"`) {
		t.Fatalf("a delegated creator must own the space: %s", rec.Body.String())
	}
}

func TestDelegatedReadScopeCannotWrite(t *testing.T) {
	f := newDelegatedFixture(t, "mycelium")
	expect(t, f, apiCall{Method: "PUT", Path: "/api/sync/files/memory/new.md?space_id=" + f.own, Token: f.ro, Body: "hello"}, http.StatusForbidden)
	expect(t, f, apiCall{Method: "POST", Path: "/api/memory/search", Token: f.ro, Body: `{"query":"zanzibar","space_id":"` + f.own + `"}`}, http.StatusForbidden)
}

func TestDelegatedRefusesOtherUsersSpace(t *testing.T) {
	f := newDelegatedFixture(t, "mycelium")
	expect(t, f, apiCall{Method: "GET", Path: "/api/sync/tree?space_id=" + f.foreign, Token: f.rw}, http.StatusForbidden)
	expect(t, f, apiCall{Method: "PUT", Path: "/api/sync/files/memory/x.md?space_id=" + f.foreign, Token: f.rw, Body: "x"}, http.StatusForbidden)
	expect(t, f, apiCall{Method: "POST", Path: "/api/memory/search", Token: f.rw, Body: `{"query":"zanzibar","space_id":"` + f.foreign + `"}`}, http.StatusForbidden)
}

func TestDelegatedNeverReachesCommonTree(t *testing.T) {
	f := newDelegatedFixture(t, "mycelium")
	expect(t, f, apiCall{Method: "POST", Path: "/api/memory/search", Token: f.rw, Body: `{"query":"zanzibar"}`}, http.StatusForbidden)
	expect(t, f, apiCall{Method: "GET", Path: "/api/sync/tree", Token: f.rw}, http.StatusForbidden)
	expect(t, f, apiCall{Method: "GET", Path: "/api/sync/files/memory/secret.md", Token: f.rw}, http.StatusForbidden)
	expect(t, f, apiCall{Method: "PUT", Path: "/api/sync/files/memory/secret.md", Token: f.rw, Body: "x"}, http.StatusForbidden)
	rec := expect(t, f, apiCall{Method: "POST", Path: "/api/memory/search", Token: f.rw, Body: `{"query":"zanzibar","space_id":"` + f.own + `"}`}, http.StatusOK)
	if !strings.Contains(rec.Body.String(), "note.md") || strings.Contains(rec.Body.String(), "secret") {
		t.Fatalf("a space search must hold only the space's pages: %s", rec.Body.String())
	}
}

func TestDelegatedRefusesAdminRoute(t *testing.T) {
	f := newDelegatedFixture(t, "mycelium")
	expect(t, f, apiCall{Method: "GET", Path: "/api/settings", Token: f.rw}, http.StatusForbidden)
	expect(t, f, apiCall{Method: "POST", Path: "/api/tokens", Token: f.rw, Body: `{"name":"x"}`}, http.StatusForbidden)
}

func TestDelegatedRefusesRouteOffAllowlist(t *testing.T) {
	f := newDelegatedFixture(t, "mycelium")
	expect(t, f, apiCall{Method: "GET", Path: "/api/rules?space_id=" + f.own, Token: f.rw}, http.StatusForbidden)
	expect(t, f, apiCall{Method: "GET", Path: "/api/auth/me", Token: f.rw}, http.StatusForbidden)
}

func TestDelegatedRefusesWrongAudience(t *testing.T) {
	f := newDelegatedFixture(t, "mycelium")
	wrong := f.issuer.token(t, "nuage", "mycelium.read mycelium.write")
	expect(t, f, apiCall{Method: "GET", Path: "/api/spaces", Token: wrong}, http.StatusUnauthorized)
}

func TestDelegationOffKeepsOpaqueTokens(t *testing.T) {
	f := newDelegatedFixture(t, "")
	expect(t, f, apiCall{Method: "GET", Path: "/api/spaces", Token: f.rw}, http.StatusUnauthorized)
	session := sessionFor(t, f.srv, otherEmail, false)
	expect(t, f, apiCall{Method: "GET", Path: "/api/spaces", Token: session}, http.StatusOK)
}
