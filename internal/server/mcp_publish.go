package server

import (
	"context"
	"errors"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/FacileStudio/Mycelium/internal/artifacts"
	"github.com/FacileStudio/Mycelium/internal/encre"
	mcpkit "github.com/FacileStudio/tronc/mcp"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// mcpMachine is the machine an artifact published over MCP is recorded
// from when no named client acts for the person.
const mcpMachine = "mcp"

// clientName is a Registre client id that is already a short name, such as
// bulle, rather than an opaque generated id.
var clientName = regexp.MustCompile(`^[a-z][a-z0-9-]{0,31}$`)

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
	id, err := t.s.mcpCaller(req)
	if err != nil {
		return nil, mcpPublishOutput{}, err
	}
	tree, err := t.s.mcpTreeFor(id, in.SpaceID)
	if err != nil {
		return nil, mcpPublishOutput{}, err
	}
	write.Machine = machineOf(id.Actor)
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
		Title: write.Title, Machine: write.Machine,
	}, now)
	return art, false, err
}

// machineOf names the machine an artifact is recorded from: the client
// acting for the person on a delegated token, by its Registre client id when
// that id is a short name, and mcp otherwise.
func machineOf(actor string) string {
	if clientName.MatchString(actor) {
		return actor
	}
	return mcpMachine
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
