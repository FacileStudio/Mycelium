package artifacts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func writeSource(t *testing.T, dir, name, body string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}
	return path
}

func TestHTMLIsRefusedAndOldHTMLStaysUnlisted(t *testing.T) {
	dir := t.TempDir()
	now := time.Date(2026, 8, 26, 10, 0, 0, 0, time.UTC)
	src := writeSource(t, dir, "in.html", "<p>hi</p>")
	if _, err := Add(dir, Request{Source: src}, now); err != ErrNotMarkdown {
		t.Fatalf("an .html file: got %v, want ErrNotMarkdown", err)
	}
	if _, err := AddContent(dir, []byte("  <!DOCTYPE html><html></html>"), "stdin.md", Request{}, now); err != ErrNotMarkdown {
		t.Fatalf("an HTML document named .md: got %v, want ErrNotMarkdown", err)
	}
	month := filepath.Join(Dir(dir), "2026", "08")
	if err := os.MkdirAll(month, 0o755); err != nil {
		t.Fatal(err)
	}
	old := writeSource(t, month, "2026-08-01-old-page-abc123.html", "<html><title>Old page</title></html>")
	all, err := List(dir)
	if err != nil || len(all) != 0 {
		t.Fatalf("an old HTML file must not be listed, got %d (%v)", len(all), err)
	}
	if _, err := Find(dir, "2026-08-01-old-page-abc123"); err == nil {
		t.Fatal("an old HTML file must not be found")
	}
	if _, err := os.Stat(old); err != nil {
		t.Fatalf("the old HTML file was touched: %v", err)
	}
}

func TestAddRoundTripsItsMetadataMarkdown(t *testing.T) {
	dir := t.TempDir()
	src := writeSource(t, dir, "spec.md", "# DMS Architecture\n\nSome body text.")
	now := time.Date(2026, 8, 27, 10, 0, 0, 0, time.UTC)

	art, err := Add(dir, Request{Source: src, Machine: "lucy"}, now)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	if !strings.HasPrefix(art.ID, "2026-08-27-dms-architecture-") {
		t.Fatalf("id from markdown heading: got %q, want prefix 2026-08-27-dms-architecture-", art.ID)
	}

	back, err := Find(dir, art.ID)
	if err != nil {
		t.Fatalf("find exact: %v", err)
	}
	if back.Title != "DMS Architecture" || back.Machine != "lucy" {
		t.Fatalf("metadata lost: %+v", back)
	}

	raw, err := os.ReadFile(back.Path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if !strings.Contains(string(raw), "type: artifact") {
		t.Fatalf("frontmatter missing type: artifact in %s", string(raw))
	}
}

func TestReAddingReplacesInPlaceWithoutStackingFrontmatter(t *testing.T) {
	dir := t.TempDir()
	now := time.Date(2026, 8, 26, 10, 0, 0, 0, time.UTC)

	src := writeSource(t, dir, "a.md", "# Suite drift\n\nfirst")
	firstArt, err := Add(dir, Request{Source: src, Machine: "ruche"}, now)
	if err != nil {
		t.Fatalf("first add: %v", err)
	}
	stamped, err := os.ReadFile(firstArt.Path)
	if err != nil {
		t.Fatalf("read stamped: %v", err)
	}
	again := writeSource(t, dir, "b.md", string(stamped))
	if _, err := Add(dir, Request{Source: again, Title: firstArt.ID, Machine: "lucy"}, now); err != nil {
		t.Fatalf("second add: %v", err)
	}

	all, err := List(dir)
	if err != nil || len(all) != 1 {
		t.Fatalf("want exactly one artifact, got %d (%v)", len(all), err)
	}
	final, err := os.ReadFile(all[0].Path)
	if err != nil {
		t.Fatalf("read final: %v", err)
	}
	if n := strings.Count(string(final), "type: artifact"); n != 1 {
		t.Fatalf("frontmatter stacked: %d blocks", n)
	}
	if all[0].Machine != "lucy" {
		t.Fatalf("re-add kept the old machine: %q", all[0].Machine)
	}
}

func TestBodyIsStoredByteForByte(t *testing.T) {
	dir := t.TempDir()
	body := "# Deploy\n\n````callout tone=\"warning\" title=\"Before you deploy\"\nCheck `x`.\n\n```sh\nmake\n```\n````\n\n" +
		"```chart\n{\"kind\": \"bar\", \"data\": []}\n```\n\n~~~stat\nraw\n~~~\n"
	art, err := AddContent(dir, []byte(body), "a.md", Request{}, time.Date(2026, 8, 26, 10, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(art.Path)
	if !strings.HasSuffix(string(raw), "---\n\n"+body) {
		t.Fatalf("body changed on disk: %q", raw)
	}
}

func TestOldExpiredArtifactStillListsAndIsKept(t *testing.T) {
	dir := t.TempDir()
	month := filepath.Join(Dir(dir), "2025", "01")
	if err := os.MkdirAll(month, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\ntitle: \"Old plan\"\nmachine: \"ruche\"\ncreated: 2025-01-02T10:00:00Z\n" +
		"expires: 2025-02-01T10:00:00Z\ntype: artifact\n---\n\n# Old plan\n"
	path := writeSource(t, month, "2025-01-02-old-plan-abc123.md", body)

	all, err := List(dir)
	if err != nil || len(all) != 1 {
		t.Fatalf("want the old artifact listed, got %d (%v)", len(all), err)
	}
	created := time.Date(2025, 1, 2, 10, 0, 0, 0, time.UTC)
	if all[0].Title != "Old plan" || !all[0].Created.Equal(created) || !all[0].Updated.Equal(created) {
		t.Fatalf("old metadata misread: %+v", all[0])
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("an old artifact past its expires line was removed: %v", err)
	}
}

func TestUpdateKeepsIDAndCreated(t *testing.T) {
	dir := t.TempDir()
	created := time.Date(2026, 8, 27, 12, 0, 0, 0, time.UTC)
	art, err := AddContent(dir, []byte("# First\n\nv1"), "a.md", Request{Machine: "bulle"}, created)
	if err != nil {
		t.Fatal(err)
	}
	later := created.Add(48 * time.Hour)

	got, err := Update(dir, art.ID, "Renamed", []byte("v2 body"), later)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	back, err := Find(dir, art.ID)
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(back.Path)
	if got.ID != art.ID || back.Path != art.Path || !back.Created.Equal(created) || !back.Updated.Equal(later) {
		t.Fatalf("update moved the artifact: %+v", back)
	}
	if back.Title != "Renamed" || back.Machine != "bulle" || !strings.HasSuffix(string(raw), "v2 body") {
		t.Fatalf("update lost content or metadata: %+v %q", back, raw)
	}
	if _, err := Update(dir, "first", "", []byte("x"), later); err == nil {
		t.Fatal("update by a fuzzy slug must be refused")
	}
}

func TestExternalRefsNamesOnlyWhatCannotResolveFromDisk(t *testing.T) {
	page := "[style](style.css) ![inline](data:image/png;base64,AAA) [top](#top) " +
		"![cdn](https://cdn.example/x.png) ![logo](./logo.png)"

	refs := ExternalRefs([]byte(page))

	want := map[string]bool{"style.css": true, "./logo.png": true}
	if len(refs) != len(want) {
		t.Fatalf("got %v, want %v", refs, want)
	}
	for _, ref := range refs {
		if !want[ref] {
			t.Fatalf("flagged a reference that resolves: %q", ref)
		}
	}
}

func TestAddContentRecordsRawMarkdown(t *testing.T) {
	dir := t.TempDir()
	now := time.Date(2026, 8, 27, 12, 0, 0, 0, time.UTC)
	art, err := AddContent(dir, []byte("# Dynamic Report\n\nContent."), "report.md", Request{
		Machine: "lucy",
	}, now)
	if err != nil {
		t.Fatalf("add content: %v", err)
	}
	if art.Title != "Dynamic Report" {
		t.Fatalf("title: got %q, want Dynamic Report", art.Title)
	}
	if !strings.HasPrefix(art.ID, "2026-08-27-dynamic-report-") {
		t.Fatalf("id prefix: got %q", art.ID)
	}
	u := URL("https://mycelium.facile.studio", art.ID)
	wantURL := "https://mycelium.facile.studio/artifacts/" + art.ID
	if u != wantURL {
		t.Fatalf("URL: got %q, want %q", u, wantURL)
	}
}
