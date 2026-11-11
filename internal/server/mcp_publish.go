package server

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/FacileStudio/Mycelium/internal/artifacts"
	"github.com/FacileStudio/Mycelium/internal/encre"
	mcpkit "github.com/FacileStudio/tronc/mcp"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// mcpMachine is the machine an artifact published over MCP is recorded
// from.
const mcpMachine = "mcp"

// errMCPTitle is the refusal for a publish without a title.
var errMCPTitle = errors.New("title is empty: an artifact is named, and replaced, by its title")

// publishArtifact is mycelium_publish_artifact: the API's checks (1 MiB cap,
// title length, markdown only), then encre's validation inside the artifact
// store, replacing the tree's artifact of the same title when there is one.
func (t *mcpTools) publishArtifact(_ context.Context, req *sdk.CallToolRequest, in mcpPublishInput) (*sdk.CallToolResult, mcpPublishOutput, error) {
	write := ArtifactWriteRequest{Title: strings.TrimSpace(in.Title), Content: in.Content}
	if write.Title == "" {
		return nil, mcpPublishOutput{}, errMCPTitle
	}
	if refused := validateArtifactWrite(write); refused != nil {
		return nil, mcpPublishOutput{}, errors.New(refused.Message)
	}
	tree, err := t.s.mcpOneTree(req, in.SpaceID)
	if err != nil {
		return nil, mcpPublishOutput{}, err
	}
	art, replaced, err := publishInTree(tree, write, time.Now())
	if err != nil {
		return nil, mcpPublishOutput{}, t.s.publishFailure(err)
	}
	return nil, mcpPublishOutput{mcpArtifact: t.summary(art, tree), Replaced: replaced}, nil
}

// publishInTree updates the newest artifact of the tree titled like write,
// or records a new one.
func publishInTree(tree mcpTree, write ArtifactWriteRequest, now time.Time) (artifacts.Artifact, bool, error) {
	list, _ := artifacts.List(tree.Root)
	for _, existing := range list {
		if strings.EqualFold(strings.TrimSpace(existing.Title), write.Title) {
			art, err := artifacts.Update(tree.Root, existing.ID, write.Title, []byte(write.Content), now)
			return art, true, err
		}
	}
	art, err := artifacts.AddContent(tree.Root, []byte(write.Content), "mcp.md", artifacts.Request{
		Title: write.Title, Machine: mcpMachine,
	}, now)
	return art, false, err
}

// publishFailure keeps the refusals a model can fix, encre's with its line
// and column, and hides the rest.
func (s *Server) publishFailure(err error) error {
	var invalid *encre.ValidationError
	switch {
	case errors.As(err, &invalid), errors.Is(err, artifacts.ErrTooLarge), errors.Is(err, artifacts.ErrNotMarkdown):
		return err
	}
	s.Log.Error("mcp: publishing an artifact failed", slog.Any("error", err))
	return mcpkit.ErrUnavailable
}
