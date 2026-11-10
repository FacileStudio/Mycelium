package server

import (
	"context"
	"encoding/json"
	"maps"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/FacileStudio/Mycelium/internal/env"
	"github.com/FacileStudio/porte/directory"
)

func seedRegistreStores(t *testing.T, srv *Server) {
	t.Helper()
	users := map[string]User{
		"alice@x.io": {Email: "alice@x.io", Admin: true},
		"Bob@X.io":   {Email: "Bob@X.io"},
		"carol@x.io": {Email: "carol@x.io"},
	}
	spaces := map[string]*Space{
		"team":  {ID: "team", Name: "Team", Members: map[string]string{"alice@x.io": roleOwner, "Bob@X.io": roleMember, "carol@x.io": roleMember}},
		"notes": {ID: "notes", Name: "Notes", Members: map[string]string{"alice@x.io": roleOwner}},
	}
	if err := srv.saveUsers(users); err != nil {
		t.Fatal(err)
	}
	if err := srv.saveSpaces(spaces); err != nil {
		t.Fatal(err)
	}
}

func registreSnapshot() directory.Snapshot {
	team := []directory.Member{
		{Email: "alice@x.io", Role: directory.RoleOwner},
		{Subject: "b", Email: "bob@x.io", Role: directory.RoleAdmin},
		{Subject: "d", Email: "dave@x.io", Name: "Dave", Picture: "https://sso.example/avatars/d", Role: "editor"},
	}
	notes := []directory.Member{{Email: "alice@x.io", Role: directory.RoleOwner}}
	return directory.Snapshot{Version: 7, Spaces: []directory.Space{
		{ID: "r-team", Name: "team", Members: team},
		{ID: "r-notes", Name: "notes", Members: notes},
	}}
}

func readBoth(t *testing.T, srv *Server) string {
	t.Helper()
	spaces, _ := os.ReadFile(srv.spacesPath())
	users, _ := os.ReadFile(srv.usersPath())
	return string(spaces) + string(users)
}

func assertFirstApply(t *testing.T, srv *Server) {
	t.Helper()
	spaces, users := srv.loadSpaces(), srv.loadUsers()
	team := spaces["team"]
	want := map[string]string{"alice@x.io": roleOwner, "Bob@X.io": roleAdmin, "dave@x.io": roleMember}
	if team.RegistreID != "r-team" || team.Name != "team" || !maps.Equal(team.Members, want) {
		t.Fatalf("team = %+v, want linked by name with members %v", team, want)
	}
	if notes := spaces["notes"]; notes.RegistreID != "" || len(notes.Members) != 1 {
		t.Fatalf("a personal space was linked or changed: %+v", notes)
	}
	if len(spaces) != 3 {
		t.Fatalf("want a new space for r-notes beside the personal one, got %d spaces", len(spaces))
	}
	if dave := users["dave@x.io"]; dave.Admin || dave.Name != "Dave" || dave.AvatarURL != "https://sso.example/avatars/d" {
		t.Fatalf("dave provisioned as %+v", dave)
	}
	if _, dup := users["bob@x.io"]; dup {
		t.Fatal("bob was provisioned twice instead of matched case-insensitively")
	}
	for email, subject := range map[string]string{"Bob@X.io": "b", "dave@x.io": "d", "alice@x.io": ""} {
		if user := users[email]; user.Subject != subject || (subject != "" && user.Issuer != testIssuer) {
			t.Fatalf("%s holds subject %q of %q, want %q", email, user.Subject, user.Issuer, subject)
		}
	}
}

func TestApplySnapshotLinksProvisionsAndIsIdempotent(t *testing.T) {
	srv := New(t.TempDir(), "pw")
	srv.OIDC = &env.OIDC{Issuer: testIssuer, ClientID: "mycelium", ClientSecret: "secret", RedirectURL: "https://x/cb"}
	seedRegistreStores(t, srv)
	if err := srv.applySnapshot(context.Background(), registreSnapshot()); err != nil {
		t.Fatal(err)
	}
	assertFirstApply(t, srv)
	before := readBoth(t, srv)
	for _, snapshot := range []directory.Snapshot{registreSnapshot(), {Version: 8}} {
		if err := srv.applySnapshot(context.Background(), snapshot); err != nil {
			t.Fatal(err)
		}
	}
	if after := readBoth(t, srv); after != before {
		t.Fatalf("a repeated snapshot, or one missing a linked space, rewrote the stores:\n%s\n---\n%s", before, after)
	}
}

type fakeRegistre struct {
	mu       sync.Mutex
	calls    []string
	failures int
	auth     bool
}

func (f *fakeRegistre) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	user, pass, _ := r.BasicAuth()
	f.auth = user == "mycelium" && pass == "secret"
	f.calls = append(f.calls, r.Method+" "+r.URL.Path)
	if r.URL.Path == "/api/directory/import" && f.failures > 0 {
		f.failures--
		http.NotFound(w, r)
		return
	}
	_ = json.NewEncoder(w).Encode(registreSnapshot())
}

func (f *fakeRegistre) snapshot() ([]string, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return slices.Clone(f.calls), f.auth
}

func TestImportRetriesThenRunsBeforeTheFirstSync(t *testing.T) {
	fake := &fakeRegistre{failures: 2}
	issuer := httptest.NewServer(fake)
	defer issuer.Close()
	srv := New(t.TempDir(), "pw")
	srv.OIDC = &env.OIDC{Issuer: issuer.URL, ClientID: "mycelium", ClientSecret: "secret"}
	srv.registre.backoff = time.Millisecond
	seedRegistreStores(t, srv)
	nudge := httptest.NewRecorder()
	srv.Handler().ServeHTTP(nudge, httptest.NewRequest(http.MethodPost, directory.SyncPath, nil))
	if nudge.Code != http.StatusAccepted {
		t.Fatalf("nudge = %d, want 202 with no credentials", nudge.Code)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go srv.RunRegistreSync(ctx)
	deadline := time.Now().Add(5 * time.Second)
	for srv.loadSpaces()["team"].RegistreID == "" && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	calls, authed := fake.snapshot()
	importCall, snapshotCall := "POST /api/directory/import", "GET /api/directory/snapshot"
	if len(calls) < 4 || !slices.Equal(calls[:4], []string{importCall, importCall, importCall, snapshotCall}) {
		t.Fatalf("calls = %v, want three imports (two refused) before the first snapshot", calls)
	}
	if !authed {
		t.Fatal("registre was not called with the client credentials")
	}
}

func TestRegistreOwnsSharedSpaceChanges(t *testing.T) {
	srv := New(t.TempDir(), "pw")
	srv.OIDC = &env.OIDC{Issuer: "https://sso.example", ClientID: "mycelium", ClientSecret: "secret"}
	token := sessionFor(t, srv, "alice@x.io", false)
	seedRegistreStores(t, srv)
	if err := srv.applySnapshot(context.Background(), registreSnapshot()); err != nil {
		t.Fatal(err)
	}
	h := srv.Handler()
	refused := []apiCall{
		{Method: "PUT", Path: "/api/spaces/team", Body: `{"name":"Renamed"}`},
		{Method: "DELETE", Path: "/api/spaces/team"},
		{Method: "POST", Path: "/api/spaces/notes/members", Body: `{"email":"Bob@X.io"}`},
		{Method: "POST", Path: "/api/spaces/team/leave"},
	}
	for _, call := range refused {
		call.Token = token
		if rec := spReq(t, h, call); rec.Code != http.StatusConflict {
			t.Errorf("%s %s = %d, want 409", call.Method, call.Path, rec.Code)
		}
	}
	rename := spReq(t, h, apiCall{Method: "PUT", Path: "/api/spaces/notes", Token: token, Body: `{"name":"Mine"}`})
	if rename.Code != http.StatusOK {
		t.Fatalf("renaming a personal space = %d, want 200", rename.Code)
	}
}
