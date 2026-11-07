package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FacileStudio/Mycelium/internal/env"
)

func serve(srv *Server, method, path, body string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, httptest.NewRequest(method, path, strings.NewReader(body)))
	return rec
}

func TestPasswordLoginIsNotMountedUnderAnIssuer(t *testing.T) {
	srv := New(t.TempDir(), "pw")
	srv.OIDC = &env.OIDC{Issuer: "https://sso.example", ClientID: "c", ClientSecret: "s", RedirectURL: "https://x/cb"}
	if rec := serve(srv, http.MethodPost, "/api/auth/login", `{"password":"pw"}`); rec.Code != http.StatusNotFound {
		t.Fatalf("login under an issuer = %d, want 404", rec.Code)
	}
	if rec := serve(New(t.TempDir(), "pw"), http.MethodPost, "/api/auth/login", `{"password":"pw"}`); rec.Code != http.StatusOK {
		t.Fatalf("login without an issuer = %d, want 200", rec.Code)
	}
}

func TestVersionAndSecurityHeaders(t *testing.T) {
	srv := New(t.TempDir(), "pw")
	srv.Version = "1.2.3"
	rec := serve(srv, http.MethodGet, "/version", "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"1.2.3"`) {
		t.Fatalf("/version = %d %s", rec.Code, rec.Body.String())
	}
	if rec.Header().Get("X-Frame-Options") != "DENY" || rec.Header().Get("Content-Security-Policy-Report-Only") == "" {
		t.Fatalf("security headers missing: %v", rec.Header())
	}
}

func TestDeletionsRouteRejectsAnUnsignedNotice(t *testing.T) {
	srv := New(t.TempDir(), "pw")
	rec := serve(srv, http.MethodPost, "/api/registre/deletions", `{"sub":"x"}`)
	if rec.Code < 400 {
		t.Fatalf("an unsigned notice answered %d", rec.Code)
	}
}

func TestTokensKeepWorkingWithoutASession(t *testing.T) {
	srv := New(t.TempDir(), "pw")
	token, err := srv.mintToken("laptop", scopeSync, "")
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/api/sync/tree", "/api/artifacts"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("%s with a machine token = %d", path, rec.Code)
		}
	}
}
