package mcpserver

import (
	"context"
	"errors"
	"strings"

	"github.com/FacileStudio/Mycelium/internal/config"
	"github.com/FacileStudio/Mycelium/internal/memory"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// defaultSearchLimit matches the --limit default of "mycelium memory search",
// so the two surfaces return the same thing for the same query.
const defaultSearchLimit = 20

// searchInput is what a model may ask for. Query is the only required field.
type searchInput struct {
	Query string `json:"query" jsonschema:"the words to look for across the wiki"`
	Limit int    `json:"limit,omitempty" jsonschema:"how many results to return at most; 0 means 20"`
}

// searchHit is one matching line, carrying enough to open the page it came from.
type searchHit struct {
	Path                 string `json:"path" jsonschema:"page path relative to the memory root"`
	Line                 int    `json:"line" jsonschema:"line number of the match within that page"`
	Content              string `json:"content" jsonschema:"the matching line"`
	Score                int    `json:"score" jsonschema:"relevance, higher is better"`
	Date                 string `json:"date,omitempty" jsonschema:"the finding's date when it has one"`
	ChangedSinceRatified bool   `json:"changed_since_ratified" jsonschema:"true when a normative page changed after a human last ratified it, which makes this result non-authoritative"`
}

// searchOutput reports the hits.
type searchOutput struct {
	Results []searchHit `json:"results"`
	Note    string      `json:"note,omitempty" jsonschema:"why no result is marked changed_since_ratified when that state is unknown"`
}

// searchMemory ranks the local markdown tree with the same BM25 search the
// CLI runs.
func searchMemory(_ context.Context, _ *mcp.CallToolRequest, in searchInput) (*mcp.CallToolResult, searchOutput, error) {
	query := strings.TrimSpace(in.Query)
	if query == "" {
		return nil, searchOutput{}, errors.New("search_memory needs a query")
	}
	limit := in.Limit
	if limit <= 0 {
		limit = defaultSearchLimit
	}
	results, err := memory.SearchChunks(config.MemoryDir(), query)
	if err != nil {
		return nil, searchOutput{}, err
	}
	if len(results) > limit {
		results = results[:limit]
	}
	out := searchOutput{}
	changed, err := memory.ChangedPages(config.DataDir())
	if err != nil {
		out.Note = "ratification state is unknown, so no result is marked: " + err.Error()
	}
	out.Results = hits(results, changed)
	return nil, out, nil
}

// hits converts search results and marks any page whose normative content
// changed since a human ratified it. That marker is not decoration: a page in
// that state is not authoritative, which is why the CLI prints it beside every
// result and why it belongs in the structured answer too.
func hits(results []memory.SearchResult, changed map[string]bool) []searchHit {
	out := make([]searchHit, 0, len(results))
	for _, r := range results {
		out = append(out, searchHit{
			Path: r.Path, Line: r.Line, Content: r.Content, Score: r.Score,
			Date: r.Date, ChangedSinceRatified: changed[r.Path],
		})
	}
	return out
}
