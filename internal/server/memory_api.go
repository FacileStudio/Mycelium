package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/FacileStudio/Mycelium/internal/memory"
	apierrors "github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"
)

const (
	defaultHitLimit = 10
	maxHitLimit     = 50
	excerptLimit    = 240
)

// MemorySearchRequest is the POST /api/memory/search body. An empty SpaceID
// searches the common tree; anything else is checked against membership before
// it is used.
type MemorySearchRequest struct {
	Query   string `json:"query"`
	SpaceID string `json:"space_id"`
	Limit   int    `json:"limit"`
}

// MemoryHit is one ranked chunk of the wiki. Score is the BM25 score and is
// comparable only against the other hits of the same query.
type MemoryHit struct {
	Path    string  `json:"path"`
	Heading string  `json:"heading"`
	Line    int     `json:"line"`
	Score   float64 `json:"score"`
	Excerpt string  `json:"excerpt"`
	Date    string  `json:"date,omitempty"`
}

// MemorySearchResponse carries the ranked hits.
type MemorySearchResponse struct {
	Results []MemoryHit `json:"results"`
}

// memorySearchPost ranks the wiki's findings by BM25 over the markdown tree.
func (s *Server) memorySearchPost(w http.ResponseWriter, r *http.Request) {
	var req MemorySearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.WriteError(w, apierrors.Invalid("invalid body"))
		return
	}
	root, ok := s.memoryScope(w, r, req.SpaceID)
	if !ok {
		return
	}
	query := strings.TrimSpace(req.Query)
	if query == "" {
		httpjson.WriteJSON(w, http.StatusOK, MemorySearchResponse{Results: []MemoryHit{}})
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, MemorySearchResponse{
		Results: s.rankMemory(filepath.Join(root, "memory"), query, hitLimit(req.Limit)),
	})
}

// memoryScope resolves the tree to search. The space id arrives in the body
// rather than the query string, so it is put back where scopeRoot looks for it
// instead of restating the membership guard: one guard, one place.
func (s *Server) memoryScope(w http.ResponseWriter, r *http.Request, spaceID string) (string, bool) {
	if spaceID == "" {
		return s.scopeRoot(w, r)
	}
	scoped := r.Clone(r.Context())
	query := scoped.URL.Query()
	query.Set("space_id", spaceID)
	scoped.URL.RawQuery = query.Encode()
	return s.scopeRoot(w, scoped)
}

// chunkIndex is the chunk table a ranking is read back through: it gives each
// hit its heading and excerpt.
type chunkIndex struct {
	byKey  map[string]memory.Chunk
	byPath map[string][]memory.Chunk
}

func chunkIndexOf(dir string) chunkIndex {
	index := chunkIndex{byKey: map[string]memory.Chunk{}, byPath: map[string][]memory.Chunk{}}
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		rel, relErr := filepath.Rel(dir, path)
		if relErr != nil {
			return nil
		}
		for _, c := range memory.Chunks(rel, string(data)) {
			index.byKey[memory.ChunkKey(c)] = c
			index.byPath[rel] = append(index.byPath[rel], c)
		}
		return nil
	})
	return index
}

// keyAt names the chunk a matching line belongs to: the last chunk starting at
// or before it, since Chunks emits them in file order.
func (idx chunkIndex) keyAt(path string, line int) (string, bool) {
	list := idx.byPath[path]
	if len(list) == 0 {
		return "", false
	}
	best := list[0]
	for _, c := range list {
		if c.Line <= line {
			best = c
		}
	}
	return memory.ChunkKey(best), true
}

// rankMemory is the search behind the API and the MCP tool. A wiki that cannot
// be walked is logged and answers no hits rather than an error, because it is
// still not a reason to 500 a search.
func (s *Server) rankMemory(dir, query string, limit int) []MemoryHit {
	results, err := memory.SearchChunks(dir, query)
	if err != nil {
		s.Log.Error("memory search failed", slog.Any("error", err))
		return []MemoryHit{}
	}
	index := chunkIndexOf(dir)
	seen := make(map[string]bool, len(results))
	hits := []MemoryHit{}
	for _, res := range results {
		key, ok := index.keyAt(res.Path, res.Line)
		if !ok || seen[key] {
			continue
		}
		seen[key] = true
		c := index.byKey[key]
		hits = append(hits, MemoryHit{
			Path:    c.Path,
			Heading: c.Heading,
			Line:    c.Line,
			Score:   float64(res.Score) / 100,
			Excerpt: excerptOf(memory.StripBlockMeta(c.Body)),
			Date:    memory.DayString(c.Date()),
		})
		if len(hits) == limit {
			break
		}
	}
	sort.SliceStable(hits, func(i, j int) bool { return betterHit(hits[i], hits[j]) })
	return hits
}

// betterHit is a total order, so a ranking built from a map still comes out the
// same on every run.
func betterHit(a, b MemoryHit) bool {
	if a.Score != b.Score {
		return a.Score > b.Score
	}
	if a.Path != b.Path {
		return a.Path < b.Path
	}
	return a.Line < b.Line
}

func excerptOf(body string) string {
	flat := strings.Join(strings.Fields(body), " ")
	runes := []rune(flat)
	if len(runes) <= excerptLimit {
		return flat
	}
	return strings.TrimSpace(string(runes[:excerptLimit])) + "…"
}

func hitLimit(limit int) int {
	if limit <= 0 {
		return defaultHitLimit
	}
	if limit > maxHitLimit {
		return maxHitLimit
	}
	return limit
}
