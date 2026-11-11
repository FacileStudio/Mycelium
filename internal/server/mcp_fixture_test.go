package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"testing"

	"github.com/FacileStudio/Mycelium/internal/env"
	"github.com/FacileStudio/tronc/mcp/mcptest"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	ownerSubject    = "owner-sub"
	memberSubject   = "member-sub"
	strangerSubject = "stranger-sub"
	sharedSpace     = "shared01"
	strangerSpace   = "zed01"
	notesPage       = "---\ntitle: Notes\nupdated: 2026-01-01\n---\n\n# Notes\n\n### Kickoff %s\nkickoff detail\n"
)

// mcpFixture is a Mycelium server with /mcp on, an admin owner, a member of
// one shared space, and a stranger with a space of their own.
type mcpFixture struct {
	srv      *Server
	http     *httptest.Server
	issuer   *mcptest.Issuer
	resource string
}

// newMCPFixture starts the server and seeds the people, spaces and pages.
func newMCPFixture(t *testing.T) *mcpFixture {
	t.Helper()
	f := &mcpFixture{issuer: mcptest.NewIssuer(t), srv: New(t.TempDir(), "")}
	var handler http.Handler
	f.http = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { handler.ServeHTTP(w, r) }))
	t.Cleanup(f.http.Close)
	f.resource = f.http.URL + "/mcp"
	f.srv.OIDC = &env.OIDC{Issuer: f.issuer.URL, ClientID: "mycelium-client", MCPResource: f.resource}
	if err := f.srv.EnableMCP(); err != nil {
		t.Fatal(err)
	}
	f.seed(t)
	handler = f.srv.Handler()
	return f
}

// seed writes the accounts, the two spaces and a notes page in each tree.
func (f *mcpFixture) seed(t *testing.T) {
	t.Helper()
	users := map[string]User{
		"owner@facile.studio":    {Email: "owner@facile.studio", Admin: true, Issuer: f.issuer.URL, Subject: ownerSubject},
		"member@facile.studio":   {Email: "member@facile.studio", Issuer: f.issuer.URL, Subject: memberSubject},
		"stranger@facile.studio": {Email: "stranger@facile.studio", Issuer: f.issuer.URL, Subject: strangerSubject},
	}
	spaces := map[string]*Space{
		sharedSpace:   {ID: sharedSpace, Name: "Shared", Members: map[string]string{"owner@facile.studio": roleOwner, "member@facile.studio": roleMember}},
		strangerSpace: {ID: strangerSpace, Name: "Zed", Members: map[string]string{"stranger@facile.studio": roleOwner}},
	}
	if err := f.srv.saveUsers(users); err != nil {
		t.Fatal(err)
	}
	if err := f.srv.saveSpaces(spaces); err != nil {
		t.Fatal(err)
	}
	f.page(t, f.srv.DataDir, "personal")
	f.page(t, f.srv.spaceDir(sharedSpace), "shared")
	f.page(t, f.srv.spaceDir(strangerSpace), "stranger")
}

// page writes memory/notes.md and an index into a tree.
func (f *mcpFixture) page(t *testing.T, root, tag string) {
	t.Helper()
	dir := filepath.Join(root, "memory")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := []byte(fmt.Sprintf(notesPage, tag))
	if err := os.WriteFile(filepath.Join(dir, "notes.md"), body, 0o644); err != nil {
		t.Fatal(err)
	}
	index := []byte("## Pages\n\n- [Notes](notes.md): " + tag + " notes\n")
	if err := os.WriteFile(filepath.Join(dir, "index.md"), index, 0o644); err != nil {
		t.Fatal(err)
	}
}

// token mints an access token for subject bound to the resource.
func (f *mcpFixture) token(t testing.TB, subject string, scopes ...string) string {
	if len(scopes) == 0 {
		scopes = []string{"mycelium.read", "mycelium.write"}
	}
	return f.issuer.Token(t, f.resource, subject, scopes...)
}

// session connects as subject with both scopes.
func (f *mcpFixture) session(t *testing.T, subject string) *sdk.ClientSession {
	t.Helper()
	return mcptest.Connect(t, f.resource, f.token(t, subject), mcptest.ClientOptions{})
}

// TestMCPServe serves the fixture on MCP_SERVE until interrupted, for the
// official conformance suite. A request without a token gets the owner's.
func TestMCPServe(t *testing.T) {
	addr := os.Getenv("MCP_SERVE")
	if addr == "" {
		t.Skip("set MCP_SERVE=127.0.0.1:<port> to serve the endpoint")
	}
	f := newMCPFixture(t)
	token := f.token(t, ownerSubject)
	target, _ := url.Parse(f.http.URL)
	proxy := httputil.NewSingleHostReverseProxy(target)
	front := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			r.Header.Set("Authorization", "Bearer "+token)
		}
		r.Header.Del("Origin")
		r.URL.Path = "/mcp"
		proxy.ServeHTTP(w, r)
	})
	go func() { _ = http.ListenAndServe(addr, front) }()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
