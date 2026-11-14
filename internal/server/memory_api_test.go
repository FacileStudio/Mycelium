package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

const wikiPage = `---
title: Sync gotchas
type: tool
---

The sync client reconciles against a local base manifest.

### Conflict backups
A genuine edit-vs-edit conflict keeps a conflict file rather than losing an edit.

### Debounce window
The worker coalesces a burst of writes into one pass per page.
`

// wikiServer builds a server whose common tree holds one wiki page.
func wikiServer(t *testing.T) *Server {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "memory", "tools"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	page := filepath.Join(dir, "memory", "tools", "sync.md")
	if err := os.WriteFile(page, []byte(wikiPage), 0o644); err != nil {
		t.Fatalf("write page: %v", err)
	}
	return New(dir, "")
}

func searchPost(t *testing.T, h http.Handler, body map[string]any) (int, MemorySearchResponse) {
	t.Helper()
	raw, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/memory/search", bytes.NewReader(raw))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	var res MemorySearchResponse
	if rec.Code == http.StatusOK {
		if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
			t.Fatalf("decode: %v", err)
		}
	}
	return rec.Code, res
}

func TestMemorySearchRanksFindings(t *testing.T) {
	s := wikiServer(t)
	code, res := searchPost(t, s.Handler(), map[string]any{"query": "conflict backup"})
	if code != http.StatusOK {
		t.Fatalf("status %d, want 200", code)
	}
	if len(res.Results) == 0 {
		t.Fatal("search returned nothing")
	}
	top := res.Results[0]
	if top.Path != filepath.Join("tools", "sync.md") || top.Heading != "Conflict backups" {
		t.Errorf("top hit = %+v, want the Conflict backups finding", top)
	}
	if top.Excerpt == "" || top.Score <= 0 {
		t.Errorf("hit carries no excerpt or score: %+v", top)
	}
}

func TestMemorySearchHonoursTheLimit(t *testing.T) {
	s := wikiServer(t)
	_, res := searchPost(t, s.Handler(), map[string]any{"query": "the", "limit": 1})
	if len(res.Results) != 1 {
		t.Fatalf("got %d hits, want 1", len(res.Results))
	}
}

func TestMemorySearchRejectsUnknownSpace(t *testing.T) {
	s := wikiServer(t)
	code, _ := searchPost(t, s.Handler(), map[string]any{"query": "conflict", "space_id": "nope"})
	if code != http.StatusForbidden {
		t.Fatalf("status %d, want 403", code)
	}
}
