package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/FacileStudio/Mycelium/internal/artifacts"
)

func seedArtifact(t *testing.T, dataDir string) string {
	t.Helper()
	src := filepath.Join(t.TempDir(), "audit.md")
	if err := os.WriteFile(src, []byte("# Audit 2026\n\nFindings."), 0o644); err != nil {
		t.Fatal(err)
	}
	art, err := artifacts.Add(dataDir, artifacts.Request{
		Source:  src,
		Machine: "lucy",
	}, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	return art.ID
}

func TestArtifactsListAPI(t *testing.T) {
	dataDir := t.TempDir()
	srv := New(dataDir, "password")
	handler := srv.Handler()
	token := loginAs(t, handler, "password", "lucy")
	id := seedArtifact(t, dataDir)

	req := httptest.NewRequest("GET", "/api/artifacts", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /api/artifacts: got %d, want 200", w.Code)
	}
	var list []ArtifactSummary
	if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].ID != id {
		t.Fatalf("expected %s in list, got %+v", id, list)
	}
}

func TestOldHTMLFilesAreNeitherListedNorServed(t *testing.T) {
	dataDir := t.TempDir()
	handler := New(dataDir, "password").Handler()
	token := loginAs(t, handler, "password", "lucy")
	stray := filepath.Join(artifacts.Dir(dataDir), "2026", "08", "2026-08-01-old-page-abc123.html")
	if err := os.MkdirAll(filepath.Dir(stray), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(stray, []byte("<html><title>Old page</title></html>"), 0o644); err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{"/api/artifacts", "/api/artifacts/2026-08-01-old-page-abc123"} {
		req := httptest.NewRequest("GET", path, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if path == "/api/artifacts" && strings.TrimSpace(w.Body.String()) != "[]" {
			t.Fatalf("an old HTML file was listed: %s", w.Body.String())
		}
		if path != "/api/artifacts" && w.Code != http.StatusNotFound {
			t.Fatalf("an old HTML file was served: got %d", w.Code)
		}
	}
	if _, err := os.Stat(stray); err != nil {
		t.Fatalf("the old HTML file was touched: %v", err)
	}
}

func TestArtifactsGetAPI(t *testing.T) {
	dataDir := t.TempDir()
	srv := New(dataDir, "password")
	handler := srv.Handler()
	token := loginAs(t, handler, "password", "lucy")
	id := seedArtifact(t, dataDir)

	req := httptest.NewRequest("GET", "/api/artifacts/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /api/artifacts/%s: got %d, want 200", id, w.Code)
	}
	var detail ArtifactDetail
	if err := json.Unmarshal(w.Body.Bytes(), &detail); err != nil {
		t.Fatal(err)
	}
	if detail.ID != id || detail.Title != "Audit 2026" || detail.Machine != "lucy" {
		t.Fatalf("unexpected detail: %+v", detail)
	}
}

func TestArtifactsDeleteAPI(t *testing.T) {
	dataDir := t.TempDir()
	srv := New(dataDir, "password")
	handler := srv.Handler()
	token := loginAs(t, handler, "password", "lucy")
	id := seedArtifact(t, dataDir)

	req := httptest.NewRequest("DELETE", "/api/artifacts/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("DELETE /api/artifacts/%s: got %d, want 204", id, w.Code)
	}

	req = httptest.NewRequest("GET", "/api/artifacts/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("GET /api/artifacts/%s after delete: got %d, want 404", id, w.Code)
	}
}
