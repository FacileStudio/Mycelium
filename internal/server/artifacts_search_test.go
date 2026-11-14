package server

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/FacileStudio/Mycelium/internal/artifacts"
	"github.com/FacileStudio/tronc/mcp/mcptest"
)

type searchFixture struct {
	srv     *Server
	h       http.Handler
	admin   string
	member  string
	shared  string
	foreign string
}

func addArtifact(t *testing.T, root, body string, at time.Time) {
	t.Helper()
	if _, err := artifacts.AddContent(root, []byte(body), "a.md", artifacts.Request{Machine: "ruche"}, at); err != nil {
		t.Fatal(err)
	}
}

func newSearchFixture(t *testing.T) searchFixture {
	t.Helper()
	srv := New(t.TempDir(), "pw")
	f := searchFixture{srv: srv, shared: "shared01", foreign: "foreign01"}
	f.admin = sessionFor(t, srv, "owner@facile.studio", true)
	f.member = sessionFor(t, srv, "member@facile.studio", false)
	sessionFor(t, srv, "stranger@facile.studio", false)
	spaces := map[string]*Space{
		f.shared:  {ID: f.shared, Name: "Shared", Members: map[string]string{"owner@facile.studio": roleOwner, "member@facile.studio": roleMember}},
		f.foreign: {ID: f.foreign, Name: "Foreign", Members: map[string]string{"stranger@facile.studio": roleOwner}},
	}
	if err := srv.saveSpaces(spaces); err != nil {
		t.Fatal(err)
	}
	day := time.Date(2026, 11, 1, 9, 0, 0, 0, time.UTC)
	addArtifact(t, srv.DataDir, "# Personal plan\n\nThe encre tables scroll.", day)
	addArtifact(t, srv.spaceDir(f.shared), "# Encre release\n\nNotes.", day)
	addArtifact(t, srv.spaceDir(f.shared), "# Team notes\n\nencre ships on Friday.", day.Add(time.Hour))
	addArtifact(t, srv.spaceDir(f.foreign), "# Foreign encre\n\nSecret.", day)
	f.h = srv.Handler()
	return f
}

func (f searchFixture) search(t *testing.T, token string, query url.Values, want int) ArtifactSearchResponse {
	t.Helper()
	rec := spReq(t, f.h, apiCall{Method: "GET", Path: "/api/artifacts/search?" + query.Encode(), Token: token})
	if rec.Code != want {
		t.Fatalf("search %s answered %d, want %d: %s", query.Encode(), rec.Code, want, rec.Body.String())
	}
	var out ArtifactSearchResponse
	if want == http.StatusOK {
		if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
			t.Fatal(err)
		}
	}
	return out
}

func TestArtifactsSearchCoversOnlyReadableTrees(t *testing.T) {
	f := newSearchFixture(t)
	member := f.search(t, f.member, url.Values{"q": {"encre"}, "all": {"true"}}, http.StatusOK)
	if member.Total != 2 {
		t.Fatalf("a member sees the shared space only: %+v", member)
	}
	top := member.Results[0]
	if top.Title != "Encre release" || top.Match != artifacts.MatchTitle || top.SpaceID != f.shared || top.SpaceName != "Shared" {
		t.Fatalf("title match first, with its space: %+v", top)
	}
	if member.Results[1].Match != artifacts.MatchContent || member.Results[1].Snippet == "" {
		t.Fatalf("content match carries a snippet: %+v", member.Results[1])
	}
	admin := f.search(t, f.admin, url.Values{"q": {"encre"}, "all": {"true"}}, http.StatusOK)
	if admin.Total != 3 {
		t.Fatalf("the admin sees the common tree and the shared space: %+v", admin)
	}
	one := f.search(t, f.admin, url.Values{"q": {"encre"}}, http.StatusOK)
	if one.Total != 1 || one.Results[0].Title != "Personal plan" {
		t.Fatalf("no scope searches the default tree: %+v", one)
	}
	f.search(t, f.member, url.Values{"q": {"encre"}, "space_id": {f.foreign}}, http.StatusForbidden)
	f.search(t, f.member, url.Values{"q": {"encre"}}, http.StatusForbidden)
}

func TestArtifactsSearchRefusesBadInput(t *testing.T) {
	f := newSearchFixture(t)
	long := make([]byte, artifacts.MaxQueryRunes+1)
	for i := range long {
		long[i] = 'a'
	}
	for _, q := range []url.Values{{"q": {"  "}}, {"q": {string(long)}}, {"q": {"encre"}, "cursor": {"%%%"}}} {
		f.search(t, f.admin, q, http.StatusBadRequest)
	}
}

func TestArtifactsSearchPages(t *testing.T) {
	f := newSearchFixture(t)
	first := f.search(t, f.admin, url.Values{"q": {"encre"}, "all": {"true"}, "limit": {"2"}}, http.StatusOK)
	if len(first.Results) != 2 || first.NextCursor == "" || first.Total != 3 {
		t.Fatalf("first page: %+v", first)
	}
	next := f.search(t, f.admin, url.Values{"q": {"encre"}, "all": {"true"}, "limit": {"2"}, "cursor": {first.NextCursor}}, http.StatusOK)
	if len(next.Results) != 1 || next.NextCursor != "" {
		t.Fatalf("last page: %+v", next)
	}
	seen := map[string]bool{}
	for _, hit := range append(first.Results, next.Results...) {
		seen[hit.ID] = true
	}
	if len(seen) != 3 {
		t.Fatalf("pages overlap: %+v %+v", first, next)
	}
}

func TestMCPSearchArtifactsStaysInScope(t *testing.T) {
	f := newMCPFixture(t)
	day := time.Date(2026, 11, 1, 9, 0, 0, 0, time.UTC)
	addArtifact(t, f.srv.DataDir, "# Owner plan\n\nThe encre tables scroll.", day)
	addArtifact(t, f.srv.spaceDir(sharedSpace), "# Shared encre\n\nNotes.", day)
	addArtifact(t, f.srv.spaceDir(strangerSpace), "# Zed encre\n\nPrivate.", day)

	var owner mcpArtifactSearchOutput
	mcptest.Decode(t, mcptest.Call(t, f.session(t, ownerSubject), toolSearchArtifacts, map[string]any{"query": "encre"}), &owner)
	if owner.Total != 2 || owner.Results[0].Title != "Shared encre" || owner.Results[0].SpaceName != "Shared" {
		t.Fatalf("owner hits: %+v", owner)
	}
	if owner.Results[1].Match != artifacts.MatchContent || owner.Results[1].Snippet == "" || owner.Results[1].URL == "" {
		t.Fatalf("personal content hit: %+v", owner.Results[1])
	}
	var member mcpArtifactSearchOutput
	mcptest.Decode(t, mcptest.Call(t, f.session(t, memberSubject), toolSearchArtifacts, map[string]any{"query": "encre"}), &member)
	if member.Total != 1 || member.Results[0].SpaceID != sharedSpace {
		t.Fatalf("member hits: %+v", member)
	}
}
