package server

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/FacileStudio/Mycelium/internal/artifacts"
	"github.com/FacileStudio/Mycelium/internal/memory"
	mcpkit "github.com/FacileStudio/tronc/mcp"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// kindDeleteArtifact is the second sensitive action: deleting an artifact.
const kindDeleteArtifact = "delete_artifact"

// waitingArtifactDelete tells the model that nothing happened yet.
const waitingArtifactDelete = "Nothing was deleted: your client cannot ask the user to confirm. Give the user confirm_url " +
	"and ask them to confirm there. Do not call this tool again for the same artifact; the request stays open for 7 days."

// artifactIn finds the artifact with exactly id in one tree.
func artifactIn(tree mcpTree, raw string) (artifacts.Artifact, error) {
	id, ok := safeName(strings.TrimSpace(raw))
	if !ok {
		return artifacts.Artifact{}, errMCPArtifactID
	}
	art, err := artifacts.Find(tree.Root, id)
	if err != nil || !strings.EqualFold(art.ID, id) {
		return artifacts.Artifact{}, fmt.Errorf("no artifact %s in %s. Take an id from %s", id, tree.Name, toolListArtifacts)
	}
	return art, nil
}

// updateArtifact is mycelium_update_artifact: PUT /api/artifacts/{id} with
// the same checks, by exact id.
func (t *mcpTools) updateArtifact(_ context.Context, req *sdk.CallToolRequest, in mcpUpdateArtifactInput) (*sdk.CallToolResult, mcpPublishOutput, error) {
	write := ArtifactWriteRequest{Title: strings.TrimSpace(in.Title), Content: in.Content}
	if refused := validateArtifactWrite(write); refused != nil {
		return nil, mcpPublishOutput{}, errors.New(refused.Message)
	}
	tree, err := t.s.mcpOneTree(req, in.SpaceID)
	if err != nil {
		return nil, mcpPublishOutput{}, err
	}
	existing, err := artifactIn(tree, in.ID)
	if err != nil {
		return nil, mcpPublishOutput{}, err
	}
	art, err := artifacts.Update(tree.Root, existing.ID, write.Title, []byte(write.Content), time.Now())
	if err != nil {
		return nil, mcpPublishOutput{}, t.s.publishFailure(err)
	}
	return nil, mcpPublishOutput{mcpArtifact: t.summary(art, tree), Replaced: true}, nil
}

// deleteArtifact is mycelium_delete_artifact. Like a page delete, it never
// runs on the first call.
func (t *mcpTools) deleteArtifact(ctx context.Context, req *sdk.CallToolRequest, in mcpDeleteArtifactInput) (*sdk.CallToolResult, mcpDeleteArtifactOutput, error) {
	if mcpkit.Retrying(req) {
		return t.resumeArtifactDelete(ctx, req)
	}
	tree, err := t.s.mcpOneTree(req, in.SpaceID)
	if err != nil {
		return nil, mcpDeleteArtifactOutput{}, err
	}
	art, err := artifactIn(tree, in.ID)
	if err != nil {
		return nil, mcpDeleteArtifactOutput{}, err
	}
	target := encodeTarget(pageTarget{SpaceID: tree.SpaceID, Path: art.ID})
	question := fmt.Sprintf("Delete the artifact %q from %s for good?", art.Title, tree.Name)
	result, outcome, err := t.confirm.Request(ctx, req, kindDeleteArtifact, target, question)
	if err != nil {
		return nil, mcpDeleteArtifactOutput{}, mcpkit.Failure(err, "Take an id from "+toolListArtifacts+".")
	}
	if result != nil {
		return result, mcpDeleteArtifactOutput{}, nil
	}
	return nil, mcpDeleteArtifactOutput{
		Status: outcome.Status, SpaceID: tree.SpaceID, ID: art.ID, ConfirmURL: outcome.ConfirmURL, Message: waitingArtifactDelete,
	}, nil
}

// resumeArtifactDelete handles the retry that carries the person's answer.
func (t *mcpTools) resumeArtifactDelete(ctx context.Context, req *sdk.CallToolRequest) (*sdk.CallToolResult, mcpDeleteArtifactOutput, error) {
	outcome, err := t.confirm.Resume(ctx, req, kindDeleteArtifact)
	if err != nil {
		return nil, mcpDeleteArtifactOutput{}, mcpkit.Failure(err, "Call "+toolDeleteArtifact+" again to ask the user afresh.")
	}
	if outcome.Status == mcpkit.StatusDeclined {
		return nil, mcpDeleteArtifactOutput{Status: mcpkit.StatusDeclined, Message: "The user declined. Nothing happened, and the artifact is kept."}, nil
	}
	target, _ := targetOf(outcome.Action)
	return nil, mcpDeleteArtifactOutput{
		Status: "deleted", SpaceID: target.SpaceID, ID: target.Path, Message: "The user confirmed and the artifact was deleted.",
	}, nil
}

// artifactOf resolves the artifact a pending delete names.
func (t *mcpTools) artifactOf(a mcpkit.Action) (artifacts.Artifact, mcpTree, error) {
	target, tree, err := t.s.actionTree(a)
	if err != nil {
		return artifacts.Artifact{}, mcpTree{}, err
	}
	art, err := artifactIn(tree, target.Path)
	if err != nil {
		return artifacts.Artifact{}, mcpTree{}, errPageGone
	}
	return art, tree, nil
}

// artifactDigest fingerprints the artifact a pending delete names.
func (t *mcpTools) artifactDigest(_ context.Context, a mcpkit.Action) (string, error) {
	art, tree, err := t.artifactOf(a)
	if err != nil {
		return "", err
	}
	raw, err := os.ReadFile(art.Path)
	if err != nil {
		return "", errPageGone
	}
	return mcpkit.Digest(tree.SpaceID, art.ID, memory.Checksum(raw)), nil
}

// removeArtifact deletes the artifact a confirmed action names.
func (t *mcpTools) removeArtifact(_ context.Context, a mcpkit.Action) error {
	art, tree, err := t.artifactOf(a)
	if err != nil {
		return err
	}
	return artifacts.Remove(tree.Root, art.ID)
}
