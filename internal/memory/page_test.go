package memory

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const wholePage = "---\ntitle: Alice\ntype: actor\n---\n\n# Alice\n\nAlice prefers short answers.\n"

func TestWritePageCreatesStampsIndexesAndLogs(t *testing.T) {
	dir := wikiWithPage(t)
	now := time.Date(2026, 11, 11, 9, 0, 0, 0, time.UTC)
	res, err := WritePage(dir, Page{Key: "people/alice", Content: wholePage, Summary: "who Alice is", Log: "people: add Alice"}, now)
	if err != nil {
		t.Fatal(err)
	}
	page := readWiki(t, dir, "people/alice.md")
	if !res.Created || !res.Indexed || res.Page != "people/alice.md" ||
		!strings.Contains(page, "updated: 2026-11-11") || !strings.Contains(page, "created: 2026-11-11") {
		t.Fatalf("new page not stamped: %+v\n%s", res, page)
	}
	if index := readWiki(t, dir, "index.md"); !strings.Contains(index, "## People\n\n- [Alice](people/alice.md): who Alice is") {
		t.Fatalf("index pointer missing:\n%s", index)
	}
	later := now.AddDate(0, 0, 3)
	res, err = WritePage(dir, Page{Key: "people/alice.md", Content: strings.Replace(wholePage, "short", "long", 1), Log: "people: Alice changed her mind"}, later)
	if err != nil {
		t.Fatal(err)
	}
	page = readWiki(t, dir, "people/alice.md")
	if res.Created || res.Indexed || !strings.Contains(page, "created: 2026-11-11") || !strings.Contains(page, "updated: 2026-11-14") ||
		!strings.Contains(page, "long answers") {
		t.Fatalf("replace must keep created and move updated: %+v\n%s", res, page)
	}
	history := readWiki(t, dir, "log.md")
	if !strings.Contains(history, "## [2026-11-11] create-page | people: add Alice") ||
		!strings.Contains(history, "## [2026-11-14] lint | people: Alice changed her mind") ||
		!strings.HasPrefix(history, "## [2026-08-20] ingest | something earlier\n") {
		t.Fatalf("log must only grow:\n%s", history)
	}
}

func TestWritePageRefusesWhatTheWikiForbids(t *testing.T) {
	dir := writeDir(t, map[string]string{"standards/go.md": "---\ntitle: Go\ntype: standard\n---\nrule\n",
		"moved.md": "---\ntitle: Moved\ntype: standard\n---\nrule\n"})
	cases := map[string]Page{
		"index.md":                {Key: "index.md", Content: wholePage, Log: "x"},
		"log.md":                  {Key: "log", Content: wholePage, Log: "x"},
		"normative page":          {Key: "standards/go.md", Content: wholePage, Log: "x"},
		"page declaring standard": {Key: "moved.md", Content: wholePage, Log: "x"},
		"new standard":            {Key: "x.md", Content: strings.Replace(wholePage, "actor", "standard", 1), Log: "x"},
		"no frontmatter":          {Key: "x.md", Content: "# Alice\n", Log: "x"},
		"no title":                {Key: "x.md", Content: "---\ntype: actor\n---\nbody\n", Log: "x"},
		"no log":                  {Key: "x.md", Content: wholePage},
		"French":                  {Key: "x.md", Content: wholePage + "\nElle préfère les réponses courtes et nous avons décidé de garder cela pour la suite.\n", Log: "x"},
		"outside the wiki":        {Key: "../rules/x.md", Content: wholePage, Log: "x"},
		"too large":               {Key: "x.md", Content: wholePage + strings.Repeat("a", MaxPageBytes), Log: "x"},
	}
	for name, page := range cases {
		if _, err := WritePage(dir, page, time.Now()); err == nil {
			t.Errorf("%s: want a refusal", name)
		}
	}
	if got := readWiki(t, dir, "standards/go.md"); !strings.Contains(got, "type: standard") {
		t.Fatalf("a standard was overwritten: %s", got)
	}
}

func TestDeletePageDropsTheIndexLineAndLogs(t *testing.T) {
	dir := wikiWithPage(t)
	if _, err := WritePage(dir, Page{Key: "people/alice.md", Content: wholePage, Log: "people: add Alice"}, time.Now()); err != nil {
		t.Fatal(err)
	}
	key, err := DeletePage(dir, "people/alice", "Alice asked to be forgotten", time.Now())
	if err != nil || key != "people/alice.md" {
		t.Fatalf("delete failed: %q %v", key, err)
	}
	if _, err := os.Stat(filepath.Join(dir, "memory", "people", "alice.md")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("the page is still there: %v", err)
	}
	if index := readWiki(t, dir, "index.md"); strings.Contains(index, "people/alice.md") || !strings.Contains(index, "tools/lefthook.md") {
		t.Fatalf("index must lose exactly the page's line:\n%s", index)
	}
	if history := readWiki(t, dir, "log.md"); !strings.Contains(history, "lint | deleted people/alice.md: Alice asked to be forgotten") {
		t.Fatalf("delete not logged:\n%s", history)
	}
	if _, err := DeletePage(dir, "people/alice", "again", time.Now()); !errors.Is(err, ErrPageMissing) {
		t.Fatalf("a missing page must be ErrPageMissing: %v", err)
	}
	for _, page := range []string{"log.md", "index.md", "standards/go.md"} {
		if _, err := DeletePage(dir, page, "why", time.Now()); err == nil || errors.Is(err, ErrPageMissing) {
			t.Errorf("%s must be refused before it is looked up: %v", page, err)
		}
	}
}
