package server

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/FacileStudio/tronc/mcp/mcptest"
)

func searchSpaces(t *testing.T, f *mcpFixture, subject string) []string {
	t.Helper()
	var out mcpSearchOutput
	mcptest.Decode(t, mcptest.Call(t, f.session(t, subject), toolSearchMemory, map[string]any{"query": "kickoff"}), &out)
	names := []string{}
	for _, hit := range out.Results {
		names = append(names, hit.SpaceName)
	}
	slices.Sort(names)
	return names
}

func TestMCPSeesWhatTheWebAppShows(t *testing.T) {
	f := newMCPFixture(t)
	if got := searchSpaces(t, f, ownerSubject); !slices.Equal(got, []string{personalTree, "Shared"}) {
		t.Fatalf("the admin owner must search their personal tree and their space, got %v", got)
	}
	if got := searchSpaces(t, f, memberSubject); !slices.Equal(got, []string{"Shared"}) {
		t.Fatalf("a member must search only their space, got %v", got)
	}
	member := f.session(t, memberSubject)
	page := mcptest.Call(t, member, toolReadMemoryPage, map[string]any{"path": "notes"})
	if !page.IsError || !strings.Contains(mcptest.Text(page), "no personal tree") {
		t.Fatalf("a non-admin must not read the common tree: %s", mcptest.Text(page))
	}
	foreign := mcptest.Call(t, member, toolReadMemoryPage, map[string]any{"path": "notes", "space_id": strangerSpace})
	if !foreign.IsError || !strings.Contains(mcptest.Text(foreign), "no space") {
		t.Fatalf("a member must not read a space they are not in: %s", mcptest.Text(foreign))
	}
	var index mcpIndexOutput
	mcptest.Decode(t, mcptest.Call(t, member, toolReadMemoryIndex, map[string]any{}), &index)
	if len(index.Spaces) != 1 || index.Spaces[0].SpaceID != sharedSpace || index.Index != "" {
		t.Fatalf("a member's index call lists their one space and no common index: %+v", index)
	}
	var shared mcpPageOutput
	mcptest.Decode(t, mcptest.Call(t, member, toolReadMemoryPage, map[string]any{"path": "notes.md", "space_id": sharedSpace}), &shared)
	if !strings.Contains(shared.Content, "Kickoff shared") || shared.NextLine != 0 {
		t.Fatalf("the member reads their space's page: %+v", shared)
	}
}

func TestMCPCreatesTheAccountOnFirstCallNeverAsAdmin(t *testing.T) {
	f := newMCPFixture(t)
	token := f.issuer.Mint(t, map[string]any{
		"aud": []string{f.resource}, "sub": "newcomer-sub", "scope": "mycelium.read",
		"email": "newcomer@facile.studio", "email_verified": true,
	})
	var index mcpIndexOutput
	mcptest.Decode(t, mcptest.Call(t, mcptest.Connect(t, f.resource, token, mcptest.ClientOptions{}), toolReadMemoryIndex, map[string]any{}), &index)
	user, ok := f.srv.loadUsers()["newcomer@facile.studio"]
	if !ok || user.Admin || user.Subject != "newcomer-sub" || len(index.Spaces) != 0 {
		t.Fatalf("first call must create a non-admin account with no tree: %+v %+v", user, index)
	}
	unverified := f.issuer.Mint(t, map[string]any{
		"aud": []string{f.resource}, "sub": "unverified-sub", "scope": "mycelium.read", "email": "u@facile.studio",
	})
	req, _ := http.NewRequest(http.MethodPost, f.resource, strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"tools/list"}`))
	req.Header.Set("Authorization", "Bearer "+unverified)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("a token without a verified email must not create an account: %d", resp.StatusCode)
	}
}

func TestMCPAddFindingFilesOnAnExistingPage(t *testing.T) {
	f := newMCPFixture(t)
	owner := f.session(t, ownerSubject)
	args := map[string]any{"page": "notes", "title": "Tokens rotate", "source": "direct observation",
		"body": "The refresh token rotates on every use.", "log": "notes: token rotation"}
	var out mcpFindingOutput
	mcptest.Decode(t, mcptest.Call(t, owner, toolAddFinding, args), &out)
	page, _ := os.ReadFile(filepath.Join(f.srv.DataDir, "memory", "notes.md"))
	history, _ := os.ReadFile(filepath.Join(f.srv.DataDir, "memory", "log.md"))
	if out.Page != "notes.md" || !strings.Contains(string(page), "### Tokens rotate") ||
		!strings.Contains(string(history), "notes: token rotation") {
		t.Fatalf("finding not filed: %+v\n%s\n%s", out, page, history)
	}
	args["page"] = "missing"
	if res := mcptest.Call(t, owner, toolAddFinding, args); !res.IsError || !strings.Contains(mcptest.Text(res), "no page missing.md") {
		t.Fatalf("a missing page must be refused by name: %s", mcptest.Text(res))
	}
	args["page"], args["body"] = "notes", "Le jeton est renouvelé à chaque utilisation, ce qui est très important pour nous."
	if res := mcptest.Call(t, owner, toolAddFinding, args); !res.IsError || !strings.Contains(mcptest.Text(res), "English-only") {
		t.Fatalf("French prose must be refused: %s", mcptest.Text(res))
	}
}

func TestMCPLeavesAnEmptyInstanceToItsOwner(t *testing.T) {
	srv := New(t.TempDir(), "")
	ident := oidcIdentity{Issuer: "https://sso.example", Subject: "s", Email: "a@example.com", EmailVerified: true}
	if err := srv.provisionIdentity(context.Background(), ident); err != errMCPNoOwner || len(srv.loadUsers()) != 0 {
		t.Fatalf("an empty instance must not get its first account from MCP: %v", err)
	}
}
