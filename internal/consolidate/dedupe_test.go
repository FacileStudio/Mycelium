package consolidate

import (
	"os"
	"path/filepath"
	"testing"
)

const dedupeFixturePage = `---
title: Mycelium codebase
type: project
---

### The daemon regenerates agent configs every tick unless stamped
**Date**: 2026-08-01
**Source**: direct observation, internal/daemon/daemon.go
The daemon ticks every minute and regenerating agent configs is write-heavy, so
the install step is gated by a stamp file checked every five minutes.
`

func dedupeWiki(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "projects"), 0o755); err != nil {
		t.Fatal(err)
	}
	page := filepath.Join(dir, "projects", "mycelium.md")
	if err := os.WriteFile(page, []byte(dedupeFixturePage), 0o600); err != nil {
		t.Fatal(err)
	}
	return dir
}

type decideCase struct {
	name      string
	emptyWiki bool
	candidate Candidate
	want      Outcome
	wantMatch bool
}

func decideCaseWiki(t *testing.T, emptyWiki bool) string {
	t.Helper()
	if emptyWiki {
		return t.TempDir()
	}
	return dedupeWiki(t)
}

func checkDecision(t *testing.T, got Decision, tt decideCase) {
	t.Helper()
	if got.Outcome != tt.want {
		t.Fatalf("outcome = %v, want %v (similarity %.2f)", got.Outcome, tt.want, got.Similarity)
	}
	if (got.Match != nil) != tt.wantMatch {
		t.Fatalf("match present = %v, want %v", got.Match != nil, tt.wantMatch)
	}
	if tt.wantMatch && got.Match.Path != "projects/mycelium.md" {
		t.Fatalf("match path = %q", got.Match.Path)
	}
}

var (
	nearDuplicate = Candidate{
		Text:        "The daemon ticks every minute and regenerating agent configs is write-heavy, so the install step is gated by a stamp file checked every five minutes.",
		EpisodeRefs: []string{"pi@2026-08-24T10:00:00Z"},
		Kind:        KindGotcha,
	}
	unrelated = Candidate{
		Text:        "Zebra quantum marmalade escalates seventeen umbrellas beneath the orchard",
		EpisodeRefs: []string{"pi@2026-08-24T10:00:00Z"},
	}
	contradicting = Candidate{
		Text: "The daemon ticks every minute and regenerating agent configs is write-heavy, so the install step is gated by a stamp file checked every five minutes." +
			" That is no longer true: the fix was to drop the stamp gate because regeneration became cheap and idempotent.",
		EpisodeRefs: []string{"pi@2026-08-24T10:00:00Z"},
	}
)

func TestDecideOutcomes(t *testing.T) {
	tests := []decideCase{
		{
			name:      "create: empty wiki has nothing to consolidate against",
			emptyWiki: true,
			candidate: nearDuplicate,
			want:      OutcomeCreate,
			wantMatch: false,
		},
		{
			name:      "create: no strong lexical match",
			candidate: unrelated,
			want:      OutcomeCreate,
			wantMatch: false,
		},
		{
			name:      "noop: near-duplicate",
			candidate: nearDuplicate,
			want:      OutcomeNoop,
			wantMatch: true,
		},
		{
			name:      "noop: a contradiction is never struck automatically",
			candidate: contradicting,
			want:      OutcomeNoop,
			wantMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Deduper{MemoryPath: decideCaseWiki(t, tt.emptyWiki)}
			got, err := d.Decide(tt.candidate)
			if err != nil {
				t.Fatalf("Decide: %v", err)
			}
			checkDecision(t, got, tt)
		})
	}
}
