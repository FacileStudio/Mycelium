package memory

import (
	"strings"
	"time"
)

// SearchChunks ranks `### finding` blocks rather than whole pages. A page is a
// mixed bag — muse.md carries a hundred unrelated findings — so scoring the
// block that actually answers the query beats scoring the page that happens to
// contain it. The text scored is the enriched chunk text, so the page's title
// and type count toward the match.
func SearchChunks(memoryPath, query string) ([]SearchResult, error) {
	terms := queryTerms(query)
	if len(terms) == 0 {
		return nil, nil
	}
	units, err := readChunkDocs(memoryPath)
	if err != nil {
		return nil, err
	}
	c := newCorpus(units)
	weights := c.weights(terms)

	scores, _ := c.rank(weights)
	var results []SearchResult
	for i, d := range c.docs {
		if scores[i] <= 0 {
			continue
		}
		results = append(results, SearchResult{
			Path:    d.page,
			Line:    d.line,
			Content: d.display,
			Score:   int(scores[i] * 100),
			Date:    d.date,
		})
	}
	sortResults(results)
	return results, nil
}

func readChunkDocs(memoryPath string) ([]doc, error) {
	pages, err := readDocs(memoryPath)
	if err != nil {
		return nil, err
	}
	var chunks []Chunk
	for _, p := range pages {
		chunks = append(chunks, Chunks(p.path, p.body)...)
	}
	replaced := supersededIDs(chunks)
	now := time.Now()

	units := make([]doc, 0, len(chunks))
	for _, c := range chunks {
		unit := newUnit(ChunkKey(c), c.Path, c.Line, c.Text())
		unit.links = append(unit.links, splitRelated(c.Related)...)
		unit.display = chunkDisplay(c)
		unit.weight = chunkWeight(c, replaced, now)
		unit.date = DayString(c.Date())
		units = append(units, unit)
	}
	return units, nil
}

// chunkDisplay is what a reader sees for a hit: the block's own heading and the
// first line of its body. The enrichment header is deliberately left out: it
// repeats the path already printed beside the line number.
//
// The excerpt is drawn from the same text that was indexed, struck-through
// spans removed. Several findings in this wiki open with the claim they go on
// to retract, and leading a hit with it hands a reader the one sentence the
// page exists to correct.
func chunkDisplay(c Chunk) string {
	body := excerpt(dropStruck(c.Body))
	if c.Heading == "" {
		return body
	}
	if body == "" {
		return c.Heading
	}
	return c.Heading + " — " + body
}

// StripBlockMeta drops the provenance lines a wiki finding opens with. They are
// identical on every block, so leading an excerpt with them spends a reader's
// attention — and an agent's tokens — on nothing that says why the block
// matched.
func StripBlockMeta(body string) string {
	var kept []string
	for _, line := range strings.Split(body, "\n") {
		if len(kept) == 0 && (strings.TrimSpace(line) == "" || isBlockMeta(strings.TrimSpace(line))) {
			continue
		}
		kept = append(kept, line)
	}
	return strings.TrimSpace(strings.Join(kept, "\n"))
}

// isBlockMeta matches the provenance lines every wiki finding opens with. They
// are the same on every block, so showing one as a hit's excerpt tells a reader
// nothing about why the block matched.
func isBlockMeta(line string) bool {
	for _, prefix := range []string{"**Date**", "**Source**", "**Sources**", "---"} {
		if strings.HasPrefix(line, prefix) {
			return true
		}
	}
	return false
}

func excerpt(body string) string {
	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !isBlockMeta(trimmed) {
			return trimmed
		}
	}
	return strings.TrimSpace(body)
}
