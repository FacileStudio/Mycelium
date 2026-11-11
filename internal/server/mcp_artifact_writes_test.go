package server

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/FacileStudio/Mycelium/internal/artifacts"
	"github.com/FacileStudio/Mycelium/internal/encre/encretest"
	"github.com/FacileStudio/tronc/mcp/mcptest"
)

// freshArtifact records a new artifact in root and returns its id.
func freshArtifact(t testing.TB, root string) string {
	t.Helper()
	art, err := artifacts.AddContent(root, []byte("# Plan\n\nship it"), "plan.md",
		artifacts.Request{Title: fmt.Sprintf("Plan %d", pageSeq.Add(1)), Machine: "test"}, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	return art.ID
}

func TestMCPConformanceConfirmsArtifactDeletes(t *testing.T) {
	f := newMCPFixture(t)
	mcptest.Run(t, mcptest.Config{
		URL:   f.resource,
		Token: func(t testing.TB, scopes ...string) string { return f.token(t, ownerSubject, scopes...) },
		Confirm: []mcptest.ConfirmCase{{
			Tool: toolDeleteArtifact,
			Setup: func(t testing.TB) (map[string]any, func() int) {
				id := freshArtifact(t, f.srv.DataDir)
				runs := func() int {
					if _, err := artifacts.Find(f.srv.DataDir, id); err != nil {
						return 1
					}
					return 0
				}
				return map[string]any{"id": id}, runs
			},
		}},
	})
}

func TestMCPUpdateArtifactByIDRenamesAndValidates(t *testing.T) {
	encretest.Install(t)
	f := newMCPFixture(t)
	owner := f.session(t, ownerSubject)
	id := freshArtifact(t, f.srv.DataDir)
	var out mcpPublishOutput
	mcptest.Decode(t, mcptest.Call(t, owner, toolUpdateArtifact, map[string]any{"id": id, "title": "Renamed", "content": "# Plan\n\nshipped"}), &out)
	art, _ := artifacts.Find(f.srv.DataDir, id)
	if out.ID != id || out.Title != "Renamed" || art.Title != "Renamed" {
		t.Fatalf("update must keep the id and rename: %+v %+v", out, art)
	}
	refusals := map[string]string{"<!doctype html><html></html>": "markdown only", "# Bad\n\nINVALID": "forbidden character"}
	for content, want := range refusals {
		if res := mcptest.Call(t, owner, toolUpdateArtifact, map[string]any{"id": id, "content": content}); !res.IsError || !strings.Contains(mcptest.Text(res), want) {
			t.Errorf("encre must still validate an update: want %q, got %s", want, mcptest.Text(res))
		}
	}
	member := f.session(t, memberSubject)
	if res := mcptest.Call(t, member, toolUpdateArtifact, map[string]any{"id": id, "content": "# x", "space_id": sharedSpace}); !res.IsError {
		t.Fatalf("an artifact of another tree must not be found: %s", mcptest.Text(res))
	}
}
