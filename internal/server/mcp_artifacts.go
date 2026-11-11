package server

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/FacileStudio/Mycelium/internal/artifacts"
	mcpkit "github.com/FacileStudio/tronc/mcp"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// mcpArtifactPage is how many artifacts one list answer holds.
const mcpArtifactPage = 50

// errMCPArtifactID is the refusal for an id that cannot name an artifact.
var errMCPArtifactID = errors.New("id is not an artifact id: take one from mycelium_list_artifacts")

// mcpArtifactCursor is the position a list cursor carries.
type mcpArtifactCursor struct {
	Offset int `json:"o"`
}

// listArtifacts is mycelium_list_artifacts.
func (t *mcpTools) listArtifacts(_ context.Context, req *sdk.CallToolRequest, in mcpArtifactListInput) (*sdk.CallToolResult, mcpArtifactListOutput, error) {
	var at mcpArtifactCursor
	if err := mcpkit.DecodeCursor(in.Cursor, &at); err != nil || at.Offset < 0 {
		return nil, mcpArtifactListOutput{}, mcpkit.ErrCursor
	}
	_, trees, err := t.s.mcpScope(req, in.SpaceID)
	if err != nil {
		return nil, mcpArtifactListOutput{}, err
	}
	words := strings.Fields(strings.ToLower(in.Query))
	all := []mcpArtifact{}
	for _, tree := range trees {
		list, _ := artifacts.List(tree.Root)
		for _, art := range list {
			if titleHolds(art.Title, words) {
				all = append(all, t.summary(art, tree))
			}
		}
	}
	sort.SliceStable(all, func(i, j int) bool { return all[i].Created > all[j].Created })
	out := mcpArtifactListOutput{Total: len(all), Artifacts: all[min(at.Offset, len(all)):]}
	if len(out.Artifacts) > mcpArtifactPage {
		out.Artifacts = out.Artifacts[:mcpArtifactPage]
		out.NextCursor = mcpkit.EncodeCursor(mcpArtifactCursor{Offset: at.Offset + mcpArtifactPage})
	}
	return nil, out, nil
}

// titleHolds reports whether a title holds every word.
func titleHolds(title string, words []string) bool {
	lower := strings.ToLower(title)
	for _, word := range words {
		if !strings.Contains(lower, word) {
			return false
		}
	}
	return true
}

// summary renders one artifact of a tree. Only the personal tree's
// artifacts get a link: the web app opens a space's artifact only while
// that space is the selected one.
func (t *mcpTools) summary(art artifacts.Artifact, tree mcpTree) mcpArtifact {
	out := mcpArtifact{
		ID: art.ID, Title: art.Title, SpaceID: tree.SpaceID, SpaceName: tree.Name,
		Created: art.Created.UTC().Format(time.RFC3339), Updated: art.Updated.UTC().Format(time.RFC3339),
	}
	if tree.SpaceID == "" {
		out.URL = artifacts.URL(t.appURL, art.ID)
	}
	return out
}

// readArtifact is mycelium_read_artifact. An exact id wins over the fuzzy
// match artifacts.Find also accepts, in whichever tree it is.
func (t *mcpTools) readArtifact(_ context.Context, req *sdk.CallToolRequest, in mcpArtifactReadInput) (*sdk.CallToolResult, mcpArtifactReadOutput, error) {
	id, ok := safeName(strings.TrimSpace(in.ID))
	if !ok {
		return nil, mcpArtifactReadOutput{}, errMCPArtifactID
	}
	_, trees, err := t.s.mcpScope(req, in.SpaceID)
	if err != nil {
		return nil, mcpArtifactReadOutput{}, err
	}
	art, tree, found := findMCPArtifact(trees, id)
	if !found {
		return nil, mcpArtifactReadOutput{}, fmt.Errorf("no artifact %s in the trees you can use. Take an id from %s", id, toolListArtifacts)
	}
	raw, err := os.ReadFile(art.Path)
	if err != nil {
		return nil, mcpArtifactReadOutput{}, mcpkit.ErrUnavailable
	}
	content, total, next := lineWindow(string(raw), in.FromLine)
	return nil, mcpArtifactReadOutput{
		mcpArtifact: t.summary(art, tree), Content: content, TotalLines: total, NextLine: next,
	}, nil
}

// findMCPArtifact looks for an artifact across trees, an exact id first.
func findMCPArtifact(trees []mcpTree, id string) (artifacts.Artifact, mcpTree, bool) {
	var fuzzy *artifacts.Artifact
	var fuzzyTree mcpTree
	for _, tree := range trees {
		art, err := artifacts.Find(tree.Root, id)
		if err != nil {
			continue
		}
		if strings.EqualFold(art.ID, id) {
			return art, tree, true
		}
		if fuzzy == nil {
			fuzzy, fuzzyTree = &art, tree
		}
	}
	if fuzzy == nil {
		return artifacts.Artifact{}, mcpTree{}, false
	}
	return *fuzzy, fuzzyTree, true
}
