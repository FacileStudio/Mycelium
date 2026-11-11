package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/FacileStudio/tronc/mcp/mcptest"
)

const alicePage = "---\ntitle: Alice\ntype: actor\n---\n\n# Alice\n\nAlice prefers short answers.\n"

var pageSeq atomic.Int64

// freshPage writes a new page in root and returns its key.
func freshPage(t testing.TB, root string) string {
	t.Helper()
	key := fmt.Sprintf("people/p%d.md", pageSeq.Add(1))
	path := filepath.Join(root, "memory", filepath.FromSlash(key))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(alicePage), 0o644); err != nil {
		t.Fatal(err)
	}
	return key
}

// gone reports whether a page of root was deleted.
func gone(root, key string) bool {
	_, err := os.Stat(filepath.Join(root, "memory", filepath.FromSlash(key)))
	return errors.Is(err, os.ErrNotExist)
}

func TestMCPConformanceConfirmsEveryDelete(t *testing.T) {
	f := newMCPFixture(t)
	mcptest.Run(t, mcptest.Config{
		URL:   f.resource,
		Token: func(t testing.TB, scopes ...string) string { return f.token(t, ownerSubject, scopes...) },
		Confirm: []mcptest.ConfirmCase{{
			Tool: toolDeleteMemoryPage,
			Setup: func(t testing.TB) (map[string]any, func() int) {
				key := freshPage(t, f.srv.DataDir)
				runs := func() int {
					if gone(f.srv.DataDir, key) {
						return 1
					}
					return 0
				}
				return map[string]any{"path": key, "log": "people: the page is obsolete"}, runs
			},
		}},
	})
}

func TestMCPWritePageCreatesAndReplacesInTheCallersTree(t *testing.T) {
	f := newMCPFixture(t)
	member := f.session(t, memberSubject)
	args := map[string]any{"path": "people/alice", "content": alicePage, "summary": "who Alice is", "log": "people: add Alice", "space_id": sharedSpace}
	var out mcpWritePageOutput
	mcptest.Decode(t, mcptest.Call(t, member, toolWriteMemoryPage, args), &out)
	root := f.srv.spaceDir(sharedSpace)
	page, _ := os.ReadFile(filepath.Join(root, "memory", "people", "alice.md"))
	index, _ := os.ReadFile(filepath.Join(root, "memory", "index.md"))
	history, _ := os.ReadFile(filepath.Join(root, "memory", "log.md"))
	if !out.Created || !out.Indexed || out.Path != "people/alice.md" || !strings.Contains(string(page), "updated: ") ||
		!strings.Contains(string(index), "(people/alice.md): who Alice is") || !strings.Contains(string(history), "create-page | people: add Alice") {
		t.Fatalf("page not created with its bookkeeping: %+v\n%s\n%s\n%s", out, page, index, history)
	}
	args["content"] = strings.Replace(alicePage, "short", "long", 1)
	mcptest.Decode(t, mcptest.Call(t, member, toolWriteMemoryPage, args), &out)
	if out.Created || out.Indexed {
		t.Fatalf("a second write replaces: %+v", out)
	}
	delete(args, "space_id")
	if res := mcptest.Call(t, member, toolWriteMemoryPage, args); !res.IsError || !strings.Contains(mcptest.Text(res), "no personal tree") {
		t.Fatalf("a non-admin must not write the common tree: %s", mcptest.Text(res))
	}
	args["space_id"] = strangerSpace
	if res := mcptest.Call(t, member, toolWriteMemoryPage, args); !res.IsError || !strings.Contains(mcptest.Text(res), "no space") {
		t.Fatalf("a member must not write a space they are not in: %s", mcptest.Text(res))
	}
	var pages mcpListPagesOutput
	mcptest.Decode(t, mcptest.Call(t, member, toolListMemoryPages, map[string]any{"space_id": sharedSpace}), &pages)
	if pages.Total != 2 || pages.Pages[0].Path != "notes.md" || pages.Pages[1].Title != "Alice" {
		t.Fatalf("the list holds the pages, not index.md or log.md: %+v", pages)
	}
}

func TestMCPPageToolsRefuseLogIndexAndStandards(t *testing.T) {
	f := newMCPFixture(t)
	standard := filepath.Join(f.srv.DataDir, "memory", "standards", "go.md")
	_ = os.MkdirAll(filepath.Dir(standard), 0o755)
	_ = os.WriteFile(standard, []byte("---\ntitle: Go\ntype: standard\n---\nrule\n"), 0o644)
	owner := f.session(t, ownerSubject)
	refusals := map[string]string{"log.md": "append-only", "index.md": "maintained", "standards/go.md": "normative"}
	for path, want := range refusals {
		write := mcptest.Call(t, owner, toolWriteMemoryPage, map[string]any{"path": path, "content": alicePage, "log": "x"})
		del := mcptest.Call(t, owner, toolDeleteMemoryPage, map[string]any{"path": path, "log": "x"})
		if !write.IsError || !strings.Contains(mcptest.Text(write), want) {
			t.Errorf("write %s: want a refusal naming %q, got %s", path, want, mcptest.Text(write))
		}
		if !del.IsError || !strings.Contains(mcptest.Text(del), want) {
			t.Errorf("delete %s: want a refusal naming %q, got %s", path, want, mcptest.Text(del))
		}
	}
	noFront := mcptest.Call(t, owner, toolWriteMemoryPage, map[string]any{"path": "x", "content": "# Alice\n", "log": "x"})
	if !noFront.IsError || !strings.Contains(mcptest.Text(noFront), "frontmatter") {
		t.Errorf("a page without frontmatter must be refused: %s", mcptest.Text(noFront))
	}
	missing := mcptest.Call(t, owner, toolDeleteMemoryPage, map[string]any{"path": "nowhere.md", "log": "x"})
	if !missing.IsError || !strings.Contains(mcptest.Text(missing), "no page nowhere.md") {
		t.Errorf("a missing page is refused before anyone is asked: %s", mcptest.Text(missing))
	}
	if body, _ := os.ReadFile(standard); !strings.Contains(string(body), "type: standard") {
		t.Fatalf("the standard changed: %s", body)
	}
}

func TestMCPCreateSpaceIsIdempotentByName(t *testing.T) {
	f := newMCPFixture(t)
	member := f.session(t, memberSubject)
	var first, second mcpCreateSpaceOutput
	mcptest.Decode(t, mcptest.Call(t, member, toolCreateSpace, map[string]any{"name": "Bulle", "description": "memory"}), &first)
	mcptest.Decode(t, mcptest.Call(t, member, toolCreateSpace, map[string]any{"name": "bulle"}), &second)
	space := f.srv.loadSpaces()[first.SpaceID]
	if !first.Created || second.Created || second.SpaceID != first.SpaceID || space == nil ||
		space.Members["member@facile.studio"] != roleOwner || len(space.Members) != 1 {
		t.Fatalf("want one space owned by the caller: %+v %+v %+v", first, second, space)
	}
	var out mcpWritePageOutput
	mcptest.Decode(t, mcptest.Call(t, member, toolWriteMemoryPage, map[string]any{
		"path": "profile.md", "content": alicePage, "log": "profile: start", "space_id": first.SpaceID,
	}), &out)
	if !out.Created {
		t.Fatalf("the new space takes pages: %+v", out)
	}
}

// pendingCall calls a confirm page route with a bearer token.
func pendingCall(t *testing.T, f *mcpFixture, method, path, token string) (int, map[string]any) {
	t.Helper()
	req, _ := http.NewRequest(method, f.http.URL+"/api/mcp/pending"+path, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	body := map[string]any{}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	return resp.StatusCode, body
}

func TestMCPDeleteWaitsForTheBrowserSession(t *testing.T) {
	f := newMCPFixture(t)
	key := freshPage(t, f.srv.DataDir)
	var out mcpDeletePageOutput
	mcptest.Decode(t, mcptest.Call(t, f.session(t, ownerSubject), toolDeleteMemoryPage, map[string]any{"path": key, "log": "people: obsolete"}), &out)
	link, _ := url.Parse(out.ConfirmURL)
	id := link.Query().Get("action")
	if out.Status != "pending_confirmation" || link.Path != "/pending" || id == "" || gone(f.srv.DataDir, key) {
		t.Fatalf("want a confirm link and the page kept: %+v", out)
	}
	cli, _ := f.srv.mintNamedSession("cli:owner@facile.studio", "owner@facile.studio", scopeAdmin)
	if status, _ := pendingCall(t, f, http.MethodPost, "/"+id+"/confirm", cli); status != http.StatusForbidden || gone(f.srv.DataDir, key) {
		t.Fatalf("a CLI token must not confirm: %d", status)
	}
	stranger, _ := f.srv.mintSessionToken("stranger@facile.studio", scopeUser)
	if status, _ := pendingCall(t, f, http.MethodPost, "/"+id+"/confirm", stranger); status != http.StatusNotFound || gone(f.srv.DataDir, key) {
		t.Fatalf("someone else's session must not confirm: %d", status)
	}
	browser, _ := f.srv.mintSessionToken("owner@facile.studio", scopeAdmin)
	status, list := pendingCall(t, f, http.MethodGet, "", browser)
	actions, _ := list["actions"].([]any)
	if status != http.StatusOK || len(actions) != 1 || actions[0].(map[string]any)["path"] != key {
		t.Fatalf("the owner's page lists the action: %d %v", status, list)
	}
	if status, body := pendingCall(t, f, http.MethodPost, "/"+id+"/confirm", browser); status != http.StatusOK || body["status"] != "deleted" || !gone(f.srv.DataDir, key) {
		t.Fatalf("the owner's confirm deletes the page: %d %v", status, body)
	}
	if status, _ := pendingCall(t, f, http.MethodPost, "/"+id+"/confirm", browser); status != http.StatusNotFound {
		t.Fatalf("a confirmation is single use: %d", status)
	}
}

func TestMCPErasureDropsPendingActions(t *testing.T) {
	f := newMCPFixture(t)
	key := freshPage(t, f.srv.spaceDir(sharedSpace))
	mcptest.Call(t, f.session(t, memberSubject), toolDeleteMemoryPage, map[string]any{"path": key, "log": "x", "space_id": sharedSpace})
	userID := fmt.Sprint(subjectNumber(f.issuer.URL, memberSubject))
	if open, _ := f.srv.mcp.pending.Open(t.Context(), userID); len(open) != 1 {
		t.Fatalf("want one waiting action, got %d", len(open))
	}
	if err := f.srv.erasePerson(t.Context(), f.issuer.URL, memberSubject); err != nil {
		t.Fatal(err)
	}
	if open, _ := f.srv.mcp.pending.Open(t.Context(), userID); len(open) != 0 {
		t.Fatalf("erasure must drop the person's actions, %d left", len(open))
	}
}
