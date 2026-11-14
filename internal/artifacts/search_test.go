package artifacts

import (
	"strings"
	"testing"
	"time"
)

func addFixture(t *testing.T, dir, body string, at time.Time) Artifact {
	t.Helper()
	art, err := AddContent(dir, []byte(body), "a.md", Request{Machine: "test"}, at)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	return art
}

func searchFixture(t *testing.T) (string, map[string]Artifact) {
	t.Helper()
	dir := t.TempDir()
	day := time.Date(2026, 11, 1, 9, 0, 0, 0, time.UTC)
	return dir, map[string]Artifact{
		"titled":  addFixture(t, dir, "# Encre tables\n\nHow tables render.", day),
		"newer":   addFixture(t, dir, "# Deploy notes\n\nThe encre tables now scroll sideways.", day.Add(48*time.Hour)),
		"older":   addFixture(t, dir, "# Old notes\n\nencre tables wrapped their text back then.", day.Add(24*time.Hour)),
		"partial": addFixture(t, dir, "# Encre only\n\nNothing about the other word.", day),
	}
}

func TestSearchRanksTitleMatchesBeforeContentMatches(t *testing.T) {
	dir, arts := searchFixture(t)
	hits, err := Search([]Tree{{Root: dir, SpaceID: "s1", SpaceName: "Team"}}, "Encre TABLES")
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	var got []string
	for _, h := range hits {
		got = append(got, h.ID)
	}
	want := []string{arts["titled"].ID, arts["newer"].ID, arts["older"].ID}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("order = %v, want %v", got, want)
	}
	if hits[0].Match != MatchTitle || hits[1].Match != MatchContent {
		t.Errorf("match = %q, %q", hits[0].Match, hits[1].Match)
	}
	if hits[0].SpaceID != "s1" || hits[0].SpaceName != "Team" {
		t.Errorf("tree not carried: %+v", hits[0])
	}
}

func TestSearchNeverMatchesFrontmatter(t *testing.T) {
	dir, _ := searchFixture(t)
	hits, err := Search([]Tree{{Root: dir}}, "machine")
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(hits) != 0 {
		t.Fatalf("frontmatter matched: %+v", hits)
	}
}

func TestSearchSnippetSitsAroundTheFirstWord(t *testing.T) {
	dir := t.TempDir()
	body := "# Long\n\n" + strings.Repeat("filler ", 60) + "the needle sits here " + strings.Repeat("tail ", 60)
	addFixture(t, dir, body, time.Date(2026, 11, 1, 9, 0, 0, 0, time.UTC))
	hits, err := Search([]Tree{{Root: dir}}, "needle")
	if err != nil || len(hits) != 1 {
		t.Fatalf("hits = %+v, err = %v", hits, err)
	}
	snip := hits[0].Snippet
	if !strings.Contains(snip, "needle") || !strings.HasPrefix(snip, "…") || !strings.HasSuffix(snip, "…") {
		t.Fatalf("snippet = %q", snip)
	}
	if n := len([]rune(snip)); n > snippetRunes+2 {
		t.Fatalf("snippet is %d runes", n)
	}
}

func TestSearchEmptyQueryFindsNothing(t *testing.T) {
	dir, _ := searchFixture(t)
	if hits, _ := Search([]Tree{{Root: dir}}, "   "); len(hits) != 0 {
		t.Fatalf("hits = %+v", hits)
	}
}

func TestSearchSnippetHidesMarkdownSyntax(t *testing.T) {
	dir := t.TempDir()
	addFixture(t, dir, "# Notes\n\n> **Wasm** builds `fast`\n\n| a | b |\n|---|---|\n| wasm | ==ok== |\n\n- [ ] ship wasm\n1. wasm again\n", time.Now())
	hits, err := Search([]Tree{{Root: dir}}, "wasm")
	if err != nil || len(hits) != 1 {
		t.Fatalf("hits %v err %v", hits, err)
	}
	if got, want := hits[0].Snippet, "Notes Wasm builds fast a b wasm ok ship wasm wasm again"; got != want {
		t.Fatalf("snippet %q, want %q", got, want)
	}
}
