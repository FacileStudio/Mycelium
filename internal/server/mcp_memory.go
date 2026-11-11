package server

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	mcpkit "github.com/FacileStudio/tronc/mcp"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// errMCPQuery is the refusal for an empty search.
var errMCPQuery = errors.New("query is empty: pass the words to look for")

// searchMemory is mycelium_search_memory: the API's hybrid search, run on
// every tree in scope and merged by fused score.
func (t *mcpTools) searchMemory(ctx context.Context, req *sdk.CallToolRequest, in mcpSearchInput) (*sdk.CallToolResult, mcpSearchOutput, error) {
	query := strings.TrimSpace(in.Query)
	if query == "" {
		return nil, mcpSearchOutput{}, errMCPQuery
	}
	_, trees, err := t.s.mcpScope(req, in.SpaceID)
	if err != nil {
		return nil, mcpSearchOutput{}, err
	}
	limit := hitLimit(in.Limit)
	out := mcpSearchOutput{Results: []mcpHit{}}
	for _, tree := range trees {
		hits, degraded := t.s.searchTree(ctx, tree, query, limit)
		out.Results = append(out.Results, hits...)
		out.Degraded = out.Degraded || degraded
	}
	sort.SliceStable(out.Results, func(i, j int) bool { return out.Results[i].Score > out.Results[j].Score })
	if len(out.Results) > limit {
		out.Results = out.Results[:limit]
	}
	return nil, out, nil
}

// searchTree is memorySearchPost's ranking over one tree.
func (s *Server) searchTree(ctx context.Context, tree mcpTree, query string, limit int) ([]mcpHit, bool) {
	dir := filepath.Join(tree.Root, "memory")
	index := chunkIndexOf(dir)
	lexical := s.lexicalRanking(dir, query, index)
	vector, degraded := s.vectorRanking(ctx, query, s.treeScope(tree.Root), index)
	hits := []mcpHit{}
	for _, hit := range fuse(index, [][]string{lexical, vector}, limit) {
		hits = append(hits, mcpHit{
			SpaceID: tree.SpaceID, SpaceName: tree.Name, Path: hit.Path, Heading: hit.Heading,
			Line: hit.Line, Date: hit.Date, Excerpt: hit.Excerpt, Score: hit.Score,
		})
	}
	return hits, degraded
}

// readMemoryPage is mycelium_read_memory_page.
func (t *mcpTools) readMemoryPage(_ context.Context, req *sdk.CallToolRequest, in mcpPageInput) (*sdk.CallToolResult, mcpPageOutput, error) {
	tree, err := t.s.mcpOneTree(req, in.SpaceID)
	if err != nil {
		return nil, mcpPageOutput{}, err
	}
	key, text, err := readMCPPage(tree, in.Path)
	if err != nil {
		return nil, mcpPageOutput{}, err
	}
	content, total, next := lineWindow(text, in.FromLine)
	return nil, mcpPageOutput{
		SpaceID: tree.SpaceID, SpaceName: tree.Name, Path: key, Content: content, TotalLines: total, NextLine: next,
	}, nil
}

// readMemoryIndex is mycelium_read_memory_index: the trees, and the index
// of the one named, or of the personal tree when there is one.
func (t *mcpTools) readMemoryIndex(_ context.Context, req *sdk.CallToolRequest, in mcpIndexInput) (*sdk.CallToolResult, mcpIndexOutput, error) {
	id, err := t.s.mcpCaller(req)
	if err != nil {
		return nil, mcpIndexOutput{}, err
	}
	out := mcpIndexOutput{Spaces: []mcpSpace{}}
	for _, tree := range t.s.mcpTrees(id) {
		out.Spaces = append(out.Spaces, mcpSpace{SpaceID: tree.SpaceID, Name: tree.Name, Role: tree.Role})
	}
	if strings.TrimSpace(in.SpaceID) == "" && !t.s.commonAllowed(id) {
		out.Note = "You have no personal tree here. Pass one of the spaces' space_id to read its index."
		if len(out.Spaces) == 0 {
			out.Note = "You belong to no space on this server yet: an admin of it adds you to one."
		}
		return nil, out, nil
	}
	tree, err := t.s.mcpTreeFor(id, in.SpaceID)
	if err != nil {
		return nil, mcpIndexOutput{}, err
	}
	out.SpaceID, out.SpaceName = tree.SpaceID, tree.Name
	out.Index, out.Truncated, out.Note = treeIndex(tree)
	return nil, out, nil
}

// treeIndex reads a tree's index.md, cut to the read budget, or says why
// there is none.
func treeIndex(tree mcpTree) (string, bool, string) {
	data, err := os.ReadFile(filepath.Join(tree.Root, "memory", "index.md"))
	if err != nil {
		return "", false, "This tree has no index.md yet."
	}
	index, cut := mcpkit.Clip(string(data), mcpTextBudget)
	return index, cut, ""
}
