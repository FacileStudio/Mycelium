package server

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
)

type artifactFixture struct {
	srv     *Server
	h       http.Handler
	owner   string
	spaceID string
}

func newArtifactFixture(t *testing.T) artifactFixture {
	t.Helper()
	srv := New(t.TempDir(), "pw")
	h := srv.Handler()
	owner := sessionFor(t, srv, "yann@facile.studio", false)
	return artifactFixture{srv: srv, h: h, owner: owner, spaceID: createSpace(t, h, owner, "Facile")}
}

func (f artifactFixture) path(suffix string) string {
	return "/api/artifacts" + suffix + "?space_id=" + f.spaceID
}

func (f artifactFixture) create(t *testing.T, body string) ArtifactDetail {
	t.Helper()
	rec := spReq(t, f.h, apiCall{Method: "POST", Path: f.path(""), Token: f.owner, Body: body})
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /api/artifacts: %d %s", rec.Code, rec.Body.String())
	}
	var detail ArtifactDetail
	if err := json.Unmarshal(rec.Body.Bytes(), &detail); err != nil {
		t.Fatal(err)
	}
	return detail
}

func TestArtifactCreateLandsInTheSpace(t *testing.T) {
	f := newArtifactFixture(t)

	detail := f.create(t, `{"title":"Weekly plan","content":"# Plan\n\n- ship","machine":"bulle"}`)

	if detail.Title != "Weekly plan" || detail.Machine != "bulle" {
		t.Fatalf("unexpected detail: %+v", detail)
	}
	if !detail.Created.Equal(detail.Updated) || !strings.Contains(detail.Content, "- ship") {
		t.Fatalf("created and updated differ or content lost: %+v", detail)
	}
	matches, _ := filepath.Glob(filepath.Join(f.srv.spaceDir(f.spaceID), "artifacts", "*", "*", detail.ID+".md"))
	if len(matches) != 1 {
		t.Fatalf("artifact not written under the space tree: %v", matches)
	}
	common, _ := filepath.Glob(filepath.Join(f.srv.DataDir, "artifacts", "*", "*", "*.md"))
	if len(common) != 0 {
		t.Fatalf("artifact leaked into the common tree: %v", common)
	}
}

func TestArtifactUpdateKeepsIDAndCreated(t *testing.T) {
	f := newArtifactFixture(t)
	first := f.create(t, `{"title":"Weekly plan","content":"v1"}`)

	rec := spReq(t, f.h, apiCall{Method: "PUT", Path: f.path("/" + first.ID), Token: f.owner, Body: `{"content":"v2"}`})
	if rec.Code != http.StatusOK {
		t.Fatalf("PUT: %d %s", rec.Code, rec.Body.String())
	}
	var got ArtifactDetail
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.ID != first.ID || !got.Created.Equal(first.Created) || got.Updated.Before(first.Updated) {
		t.Fatalf("update changed identity: before %+v after %+v", first, got)
	}
	if got.Title != "Weekly plan" || !strings.HasSuffix(got.Content, "v2") {
		t.Fatalf("update lost the title or the body: %+v", got)
	}

	rec = spReq(t, f.h, apiCall{Method: "PUT", Path: f.path("/missing"), Token: f.owner, Body: `{"content":"x"}`})
	if rec.Code != http.StatusNotFound {
		t.Fatalf("PUT on a missing id: got %d, want 404", rec.Code)
	}
}

func TestArtifactWritesRefuseBadBodies(t *testing.T) {
	f := newArtifactFixture(t)
	cases := map[string]struct {
		body string
		want int
	}{
		"html document": {`{"title":"x","content":"<!DOCTYPE html><html><body>hi</body></html>"}`, http.StatusBadRequest},
		"script tag":    {`{"title":"x","content":"# Hi\n<script>alert(1)</script>"}`, http.StatusBadRequest},
		"empty content": {`{"title":"x","content":"  "}`, http.StatusBadRequest},
		"long title":    {`{"title":"` + strings.Repeat("a", 201) + `","content":"ok"}`, http.StatusBadRequest},
		"huge content":  {`{"content":"` + strings.Repeat("a", maxArtifactContent+1) + `"}`, http.StatusRequestEntityTooLarge},
	}
	for name, c := range cases {
		rec := spReq(t, f.h, apiCall{Method: "POST", Path: f.path(""), Token: f.owner, Body: c.body})
		if rec.Code != c.want {
			t.Errorf("%s: got %d, want %d", name, rec.Code, c.want)
		}
	}
}

func TestArtifactWritesGuardTheSpace(t *testing.T) {
	f := newArtifactFixture(t)
	outsider := sessionFor(t, f.srv, "noah@facile.studio", false)
	own := f.create(t, `{"title":"Private","content":"secret plan"}`)

	calls := []apiCall{
		{Method: "POST", Path: f.path(""), Token: outsider, Body: `{"title":"x","content":"y"}`},
		{Method: "PUT", Path: f.path("/" + own.ID), Token: outsider, Body: `{"content":"y"}`},
		{Method: "GET", Path: f.path("/" + own.ID), Token: outsider, Body: ""},
	}
	for _, c := range calls {
		if rec := spReq(t, f.h, c); rec.Code != http.StatusForbidden {
			t.Errorf("%s %s by a non-member: got %d, want 403", c.Method, c.Path, rec.Code)
		}
	}
	rec := spReq(t, f.h, apiCall{Method: "PUT", Path: f.path("/..%2Fx"), Token: f.owner, Body: `{"content":"y"}`})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("traversal id: got %d, want 400", rec.Code)
	}
}
