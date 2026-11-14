package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/FacileStudio/Mycelium/internal/config"
	"github.com/FacileStudio/Mycelium/internal/memory"
)

const memoryFixture = "---\ntitle: Ranking\n---\n\n### bm25 replaced the strict AND\n" +
	"Strict AND collapsed on reformulated queries.\n"

func TestMemorySearchReadsTheLocalTree(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "memory"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "memory", "ranking.md"), []byte(memoryFixture), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	t.Setenv("DATA_DIR", dir)

	results, err := memory.SearchChunks(config.MemoryDir(), "bm25 ranking")
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(results) == 0 || results[0].Path != "ranking.md" {
		t.Fatalf("results = %+v, want ranking.md first", results)
	}
}
