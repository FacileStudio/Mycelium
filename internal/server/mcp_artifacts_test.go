package server

import (
	"strings"
	"testing"

	"github.com/FacileStudio/Mycelium/internal/encre/encretest"
	"github.com/FacileStudio/tronc/mcp/mcptest"
)

func TestMCPPublishReplacesByTitle(t *testing.T) {
	f := newMCPFixture(t)
	owner := f.session(t, ownerSubject)
	var first, second mcpPublishOutput
	mcptest.Decode(t, mcptest.Call(t, owner, toolPublishArtifact, map[string]any{"title": "Weekly", "content": "# Weekly\n\none"}), &first)
	mcptest.Decode(t, mcptest.Call(t, owner, toolPublishArtifact, map[string]any{"title": "weekly", "content": "# Weekly\n\ntwo"}), &second)
	if first.Replaced || !second.Replaced || first.ID != second.ID || !strings.HasSuffix(first.URL, "/artifacts/"+first.ID) {
		t.Fatalf("the same title must replace in place: %+v %+v", first, second)
	}
	var list mcpArtifactListOutput
	mcptest.Decode(t, mcptest.Call(t, owner, toolListArtifacts, map[string]any{"query": "week"}), &list)
	if list.Total != 1 || list.Artifacts[0].ID != first.ID {
		t.Fatalf("want one artifact, got %+v", list)
	}
	var read mcpArtifactReadOutput
	mcptest.Decode(t, mcptest.Call(t, owner, toolReadArtifact, map[string]any{"id": first.ID}), &read)
	if !strings.Contains(read.Content, "two") || strings.Contains(read.Content, "one") {
		t.Fatalf("read must return the replacement: %q", read.Content)
	}
	member := f.session(t, memberSubject)
	if res := mcptest.Call(t, member, toolReadArtifact, map[string]any{"id": first.ID}); !res.IsError {
		t.Fatalf("a member must not read the owner's personal artifact: %s", mcptest.Text(res))
	}
	var shared mcpPublishOutput
	mcptest.Decode(t, mcptest.Call(t, member, toolPublishArtifact, map[string]any{"title": "Plan", "content": "# Plan", "space_id": sharedSpace}), &shared)
	var inSpace mcpArtifactListOutput
	mcptest.Decode(t, mcptest.Call(t, owner, toolListArtifacts, map[string]any{"space_id": sharedSpace}), &inSpace)
	if inSpace.Total != 1 || inSpace.Artifacts[0].ID != shared.ID || inSpace.Artifacts[0].URL != "" {
		t.Fatalf("the owner sees the member's space artifact, without a link: %+v", inSpace)
	}
}

func TestMCPPublishRefusesWithAPosition(t *testing.T) {
	encretest.Install(t)
	f := newMCPFixture(t)
	owner := f.session(t, ownerSubject)
	cases := map[string]string{
		"# Deep\n\nNESTED":             "invalid document: line 11, column 1: custom blocks nest more than 8 levels deep",
		strings.Repeat("a", 1<<20+1):   "content exceeds 1 MiB",
		"<!doctype html><html></html>": "markdown only",
		"# Bad\n\nINVALID":             "forbidden character",
	}
	for content, want := range cases {
		res := mcptest.Call(t, owner, toolPublishArtifact, map[string]any{"title": "Refused", "content": content})
		if !res.IsError || !strings.Contains(mcptest.Text(res), want) {
			t.Errorf("want a refusal carrying %q, got %s", want, mcptest.Text(res))
		}
	}
	var list mcpArtifactListOutput
	mcptest.Decode(t, mcptest.Call(t, owner, toolListArtifacts, map[string]any{}), &list)
	if list.Total != 0 {
		t.Fatalf("a refused publish must write nothing: %+v", list)
	}
}

func TestLineWindowContinues(t *testing.T) {
	line := strings.Repeat("x", 999) + "\n"
	text := strings.Repeat(line, 100)
	first, total, next := lineWindow(text, 0)
	rest, _, after := lineWindow(text, next)
	if total != 100 || next == 0 || after != 0 || first+rest != text {
		t.Fatalf("windows must cover the text once: total %d next %d after %d", total, next, after)
	}
}
