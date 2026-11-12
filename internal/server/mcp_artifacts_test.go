package server

import (
	"strings"
	"testing"
	"time"

	"github.com/FacileStudio/Mycelium/internal/artifacts"
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

func TestMCPArtifactToolsNameTheMachine(t *testing.T) {
	f := newMCPFixture(t)
	art, err := artifacts.AddContent(f.srv.DataDir, []byte("# Report\n\nbody"), "report.md",
		artifacts.Request{Title: "Report", Machine: "ruche"}, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	owner := f.session(t, ownerSubject)
	var list mcpArtifactListOutput
	mcptest.Decode(t, mcptest.Call(t, owner, toolListArtifacts, map[string]any{"query": "report"}), &list)
	var read mcpArtifactReadOutput
	mcptest.Decode(t, mcptest.Call(t, owner, toolReadArtifact, map[string]any{"id": art.ID}), &read)
	if list.Total != 1 || list.Artifacts[0].Machine != "ruche" || read.Machine != "ruche" || read.ID != art.ID {
		t.Fatalf("list and read must carry the machine: %+v %+v", list, read.mcpArtifact)
	}
	raw := mcptest.Text(mcptest.Call(t, owner, toolListArtifacts, map[string]any{"query": "report"}))
	if !strings.Contains(raw, `"machine":"ruche"`) {
		t.Fatalf("the wire key is machine: %s", raw)
	}
}

func TestMCPPublishRecordsTheActingClient(t *testing.T) {
	f := newMCPFixture(t)
	cases := []struct {
		actor, want string
	}{
		{"bulle", "bulle"},
		{"1IdqK3bA2LWI-zuWVLqFpnIlCTVu117F", "mcp"},
		{"", "mcp"},
	}
	for _, c := range cases {
		claims := map[string]any{"aud": []string{f.resource}, "sub": ownerSubject, "scope": "mycelium.read mycelium.write"}
		if c.actor != "" {
			claims["act"] = map[string]any{"sub": c.actor}
		}
		session := mcptest.Connect(t, f.resource, f.issuer.Mint(t, claims), mcptest.ClientOptions{})
		var out mcpPublishOutput
		mcptest.Decode(t, mcptest.Call(t, session, toolPublishArtifact, map[string]any{"title": "From " + c.actor, "content": "# Notes"}), &out)
		var read mcpArtifactReadOutput
		mcptest.Decode(t, mcptest.Call(t, session, toolReadArtifact, map[string]any{"id": out.ID}), &read)
		if read.Machine != c.want {
			t.Errorf("actor %q: want machine %q, got %q", c.actor, c.want, read.Machine)
		}
	}
}
