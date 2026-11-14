package server

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"strings"
	"time"

	"github.com/FacileStudio/Mycelium/internal/memory"
	mcpkit "github.com/FacileStudio/tronc/mcp"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// addFinding is mycelium_add_finding: memory.Add, the CLI's own filing, on
// the tree named.
func (t *mcpTools) addFinding(_ context.Context, req *sdk.CallToolRequest, in mcpFindingInput) (*sdk.CallToolResult, mcpFindingOutput, error) {
	if err := in.validate(); err != nil {
		return nil, mcpFindingOutput{}, err
	}
	tree, err := t.s.mcpOneTree(req, in.SpaceID)
	if err != nil {
		return nil, mcpFindingOutput{}, err
	}
	key, err := mcpPageKey(in.Page)
	if err != nil {
		return nil, mcpFindingOutput{}, err
	}
	finding := memory.Finding{
		Page: key, Title: in.Title, Source: in.Source, Body: in.Body, Log: in.Log, Language: tree.language(),
	}
	res, err := memory.Add(tree.Root, finding, time.Now(), nil)
	if err != nil {
		return nil, mcpFindingOutput{}, t.s.findingFailure(err, key, tree)
	}
	return nil, mcpFindingOutput{SpaceID: tree.SpaceID, Page: res.Page, Indexed: res.Indexed}, nil
}

// validate names the first missing field by its argument name.
func (in mcpFindingInput) validate() error {
	for _, field := range [][2]string{{in.Page, "page"}, {in.Title, "title"}, {in.Source, "source"},
		{in.Body, "body"}, {in.Log, "log"}} {
		if strings.TrimSpace(field[0]) == "" {
			return fmt.Errorf("%s is empty: a finding needs a page, a title, a source, a body and a log line", field[1])
		}
	}
	return nil
}

// findingFailure turns memory.Add's refusals into ones a model can act on.
// Its own messages (French prose, a page without frontmatter, a busy wiki)
// pass through; a file error never does, since it carries server paths.
func (s *Server) findingFailure(err error, key string, tree mcpTree) error {
	var pathErr *fs.PathError
	if !errors.As(err, &pathErr) {
		return err
	}
	if errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("there is no page %s in %s. Findings go on an existing page: find one with %s or %s",
			key, tree.Name, toolSearchMemory, toolReadMemoryIndex)
	}
	s.Log.Error("mcp: filing a finding failed", slog.Any("error", err))
	return mcpkit.ErrUnavailable
}
