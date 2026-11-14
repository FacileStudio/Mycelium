package memory

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/FacileStudio/Mycelium/internal/config"
)

type ranker func(dir, query string) ([]SearchResult, error)

// wikiDir points the retrieval measurements at this machine's real wiki. It
// resolves through config so the XDG layout is honoured, and reports false
// when there is nothing to measure against.
func wikiDir() (string, bool) {
	dir := config.MemoryDir()
	if _, err := os.Stat(dir); err != nil {
		return "", false
	}
	return dir, true
}

func measure(t *testing.T, dir string, cases []EvalCase, rank ranker) (float64, float64) {
	t.Helper()
	hits, reciprocal := 0, 0.0
	for _, c := range cases {
		results, err := rank(dir, c.Query)
		if err != nil {
			t.Fatalf("%q: %v", c.Query, err)
		}
		wanted := map[string]bool{}
		for _, w := range c.Expect {
			wanted[filepath.FromSlash(w)] = true
		}
		seen, position := map[string]bool{}, 0
		for _, r := range results {
			if seen[r.Path] {
				continue
			}
			seen[r.Path] = true
			position++
			if wanted[r.Path] {
				if position <= evalK {
					hits++
				}
				reciprocal += 1.0 / float64(position)
				break
			}
		}
	}
	n := float64(len(cases))
	return float64(hits) / n, reciprocal / n
}

// TestChunkLevelBeatsPageLevel is the measurement that decides whether the
// chunker earns its place over page-level search.
func TestChunkLevelBeatsPageLevel(t *testing.T) {
	dir, ok := wikiDir()
	if !ok {
		t.Skip("no wiki on this machine")
	}
	cases := loadGolden(t)

	pageRecall, pageMRR := measure(t, dir, cases, Search)
	chunkRecall, chunkMRR := measure(t, dir, cases, SearchChunks)

	t.Logf("page-level   recall@%d = %.3f   MRR = %.3f", evalK, pageRecall, pageMRR)
	t.Logf("chunk-level  recall@%d = %.3f   MRR = %.3f", evalK, chunkRecall, chunkMRR)
	t.Logf("delta        recall %+.3f   MRR %+.3f", chunkRecall-pageRecall, chunkMRR-pageMRR)
}

// TestFuseRRFOrdersByRankNotScore proves fusion cannot be swayed by the scale
// of either half: a chunk ranked first by both wins over one ranked first by
// only one, whatever the raw numbers say.
