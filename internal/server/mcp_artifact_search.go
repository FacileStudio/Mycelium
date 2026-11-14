package server

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/FacileStudio/Mycelium/internal/artifacts"
	mcpkit "github.com/FacileStudio/tronc/mcp"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// mcpArtifactSearchPage is how many hits one search answer holds.
const mcpArtifactSearchPage = 20

// errMCPArtifactQuery is the refusal for an empty or oversize search.
var errMCPArtifactQuery = errors.New("query must hold 1 to 200 characters: pass the words to look for")

// searchArtifacts is mycelium_search_artifacts: artifacts.Search over every
// tree in scope.
func (t *mcpTools) searchArtifacts(_ context.Context, req *sdk.CallToolRequest, in mcpArtifactSearchInput) (*sdk.CallToolResult, mcpArtifactSearchOutput, error) {
	query := strings.TrimSpace(in.Query)
	if query == "" || utf8.RuneCountInString(query) > artifacts.MaxQueryRunes {
		return nil, mcpArtifactSearchOutput{}, errMCPArtifactQuery
	}
	var at mcpArtifactCursor
	if err := mcpkit.DecodeCursor(in.Cursor, &at); err != nil || at.Offset < 0 {
		return nil, mcpArtifactSearchOutput{}, mcpkit.ErrCursor
	}
	_, trees, err := t.s.mcpScope(req, in.SpaceID)
	if err != nil {
		return nil, mcpArtifactSearchOutput{}, err
	}
	scoped := make([]artifacts.Tree, 0, len(trees))
	byID := map[string]mcpTree{}
	for _, tree := range trees {
		scoped = append(scoped, artifacts.Tree{Root: tree.Root, SpaceID: tree.SpaceID, SpaceName: tree.Name})
		byID[tree.SpaceID] = tree
	}
	hits, err := artifacts.Search(scoped, query)
	if err != nil {
		return nil, mcpArtifactSearchOutput{}, mcpkit.ErrUnavailable
	}
	out := mcpArtifactSearchOutput{Results: []mcpArtifactHit{}, Total: len(hits)}
	rest := hits[min(at.Offset, len(hits)):]
	if len(rest) > mcpArtifactSearchPage {
		rest = rest[:mcpArtifactSearchPage]
		out.NextCursor = mcpkit.EncodeCursor(mcpArtifactCursor{Offset: at.Offset + mcpArtifactSearchPage})
	}
	for _, hit := range rest {
		out.Results = append(out.Results, mcpArtifactHit{
			mcpArtifact: t.summary(hit.Artifact, byID[hit.SpaceID]), Match: hit.Match, Snippet: hit.Snippet,
		})
	}
	return nil, out, nil
}
