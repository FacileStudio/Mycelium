package consolidate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var writeNow = time.Date(2026, 8, 24, 10, 0, 0, 0, time.UTC)

func golden(t *testing.T, name string) (string, string) {
	t.Helper()
	before, err := os.ReadFile(filepath.Join("testdata", "write", name+".before.md"))
	if err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(filepath.Join("testdata", "write", name+".after.md"))
	if err != nil {
		t.Fatal(err)
	}
	return string(before), string(after)
}

func tempPage(t *testing.T, rel, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestCreateGolden(t *testing.T) {
	before, want := golden(t, "create_existing")
	dir := tempPage(t, "projects/daemon.md", before)
	w := &Writer{MemoryPath: dir}
	c := Candidate{
		Text:        "The fix was to run sessions scan before each sync in the daemon so the shards are sealed before the files ride the sync.",
		EpisodeRefs: []string{"pi@2026-08-24T10:00:00Z"},
		Kind:        KindErrorFix,
	}
	path, err := w.Create(Decision{Outcome: OutcomeCreate, Match: &Match{Path: "projects/daemon.md", Similarity: 0.3}}, c, writeNow)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != want {
		t.Fatalf("create mismatch\n--- got ---\n%s\n--- want ---\n%s", data, want)
	}
}

func TestCreateNewPageGolden(t *testing.T) {
	dir := t.TempDir()
	w := &Writer{MemoryPath: dir}
	c := Candidate{
		Text:        "The fix was to pass --strip-trailing-cr to the CLI diff command or every line shows as changed.",
		EpisodeRefs: []string{"pi@2026-08-24T10:00:00Z"},
	}
	path, err := w.Create(Decision{Outcome: OutcomeCreate}, c, writeNow)
	if err != nil {
		t.Fatal(err)
	}
	wantPath := filepath.Join(dir, "bugs", "the-fix-was-to-pass-strip-trailing-cr-to-the-cli.md")
	if path != wantPath {
		t.Fatalf("path = %s, want %s", path, wantPath)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := `---
title: The fix was to pass --strip-trailing-cr to the CLI diff command or
type: bugs
sources:
  - consolidation
related: []
confidence: low
created: 2026-08-24
updated: 2026-08-24
---

# The fix was to pass --strip-trailing-cr to the CLI diff command or

### The fix was to pass --strip-trailing-cr to the CLI diff command or
**Date**: 2026-08-24
**Source**: consolidation, pi@2026-08-24T10:00:00Z
The fix was to pass --strip-trailing-cr to the CLI diff command or every line shows as changed.
`
	if string(data) != want {
		t.Fatalf("new page mismatch\n--- got ---\n%s\n--- want ---\n%s", data, want)
	}
}

func TestWriteAtomicNoTempLeftover(t *testing.T) {
	dir := tempPage(t, "bugs/x.md", "### one\nbody\n")
	w := &Writer{MemoryPath: dir}
	c := Candidate{Text: "a finding body that is long enough to stand alone in the wiki", EpisodeRefs: []string{"pi@2026-08-24T10:00:00Z"}}
	if _, err := w.Create(Decision{}, c, writeNow); err != nil {
		t.Fatal(err)
	}
	matches, _ := filepath.Glob(filepath.Join(dir, "**/*.tmp"))
	if len(matches) != 0 {
		t.Fatalf("temp files left behind: %v", matches)
	}
}

func TestShortTitle(t *testing.T) {
	tests := []struct{ name, in, want string }{
		{name: "short line kept whole", in: "the daemon drops stale locks", want: "the daemon drops stale locks"},
		{name: "multiline uses first line", in: "first line matters\nsecond line", want: "first line matters"},
		{name: "markdown noise stripped", in: "- `--flag` breaks the build", want: "--flag` breaks the build"},
		{
			name: "long line cut on word boundary",
			in:   "the register handler attaches the caller password to any existing account that had no password identity at all",
			want: "the register handler attaches the caller password to any existing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shortTitle(tt.in); got != tt.want {
				t.Fatalf("shortTitle = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestClassifyDir(t *testing.T) {
	tests := []struct{ name, in, want string }{
		{name: "failure words go to bugs", in: "connection refused until the retry loop was fixed", want: "bugs"},
		{name: "tool vocabulary goes to tools", in: "the CLI flag --all also cleans scripts", want: "tools"},
		{name: "person words go to people", in: "Alex is an external collaborator on GFConseil", want: "people"},
		{name: "everything else is projects", in: "the daemon seals session blocks hourly", want: "projects"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := classifyDir(tt.in); got != tt.want {
				t.Fatalf("classifyDir = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSlug(t *testing.T) {
	tests := []struct{ name, in, want string }{
		{name: "kebab cases and strips punctuation", in: "Hello, World! Again", want: "hello-world-again"},
		{name: "caps length", in: strings.Repeat("word ", 30), want: strings.Repeat("word-", 9) + "wor"},
		{name: "no leading dashes", in: "...leading", want: "leading"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slug(tt.in); got != tt.want {
				t.Fatalf("slug = %q, want %q", got, tt.want)
			}
		})
	}
}
