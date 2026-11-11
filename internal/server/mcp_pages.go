package server

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"github.com/FacileStudio/Mycelium/internal/memory"
	mcpkit "github.com/FacileStudio/tronc/mcp"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// waitingDelete tells the model that nothing happened yet, and what to do.
const waitingDelete = "Nothing was deleted: your client cannot ask the user to confirm. Give the user confirm_url and " +
	"ask them to confirm there. Do not call this tool again for the same page; the request stays open for 7 days."

// writeMemoryPage is mycelium_write_memory_page: memory.WritePage on the
// tree named, then the embedding the sync route would queue.
func (t *mcpTools) writeMemoryPage(_ context.Context, req *sdk.CallToolRequest, in mcpWritePageInput) (*sdk.CallToolResult, mcpWritePageOutput, error) {
	tree, err := t.s.mcpOneTree(req, in.SpaceID)
	if err != nil {
		return nil, mcpWritePageOutput{}, err
	}
	key, err := mcpPageKey(in.Path)
	if err != nil {
		return nil, mcpWritePageOutput{}, err
	}
	page := memory.Page{Key: key, Content: in.Content, Summary: in.Summary, Log: in.Log}
	res, err := memory.WritePage(tree.Root, page, time.Now())
	if err != nil {
		return nil, mcpWritePageOutput{}, t.s.pageFailure(err, key, tree)
	}
	for _, file := range []string{res.Page, "index.md", "log.md"} {
		t.s.enqueueEmbed(tree.Root, filepath.Join(tree.Root, "memory", filepath.FromSlash(file)))
	}
	return nil, mcpWritePageOutput{SpaceID: tree.SpaceID, Path: res.Page, Created: res.Created, Indexed: res.Indexed}, nil
}

// deleteMemoryPage is mycelium_delete_memory_page. It never deletes on the
// first call: the person confirms through MRTR when the client can ask
// them, or on the confirm page otherwise.
func (t *mcpTools) deleteMemoryPage(ctx context.Context, req *sdk.CallToolRequest, in mcpDeletePageInput) (*sdk.CallToolResult, mcpDeletePageOutput, error) {
	if mcpkit.Retrying(req) {
		return t.resumeDelete(ctx, req)
	}
	tree, err := t.s.mcpOneTree(req, in.SpaceID)
	if err != nil {
		return nil, mcpDeletePageOutput{}, err
	}
	key, err := mcpPageKey(in.Path)
	if err != nil {
		return nil, mcpDeletePageOutput{}, err
	}
	if _, _, err := memory.CheckDelete(tree.Root, key, in.Log); err != nil {
		return nil, mcpDeletePageOutput{}, t.s.pageFailure(err, key, tree)
	}
	target := encodeTarget(pageTarget{SpaceID: tree.SpaceID, Path: key, Log: strings.TrimSpace(in.Log)})
	question := fmt.Sprintf("Delete the memory page %s from %s for good? Reason given: %s", key, tree.Name, strings.TrimSpace(in.Log))
	result, outcome, err := t.confirm.Request(ctx, req, kindDeletePage, target, question)
	if err != nil {
		return nil, mcpDeletePageOutput{}, mcpkit.Failure(err, "Find page paths with "+toolListMemoryPages+".")
	}
	if result != nil {
		return result, mcpDeletePageOutput{}, nil
	}
	return nil, mcpDeletePageOutput{
		Status: outcome.Status, SpaceID: tree.SpaceID, Path: key, ConfirmURL: outcome.ConfirmURL, Message: waitingDelete,
	}, nil
}

// resumeDelete handles the retry that carries the person's answer.
func (t *mcpTools) resumeDelete(ctx context.Context, req *sdk.CallToolRequest) (*sdk.CallToolResult, mcpDeletePageOutput, error) {
	outcome, err := t.confirm.Resume(ctx, req, kindDeletePage)
	if err != nil {
		return nil, mcpDeletePageOutput{}, mcpkit.Failure(err, "Call "+toolDeleteMemoryPage+" again to ask the user afresh.")
	}
	if outcome.Status == mcpkit.StatusDeclined {
		return nil, mcpDeletePageOutput{Status: mcpkit.StatusDeclined, Message: "The user declined. Nothing happened, and the page is kept."}, nil
	}
	target, _ := targetOf(outcome.Action)
	return nil, mcpDeletePageOutput{
		Status: "deleted", SpaceID: target.SpaceID, Path: target.Path, Message: "The user confirmed and the page was deleted.",
	}, nil
}

// pageFailure turns the wiki's refusals into ones a model can act on. Its
// own messages pass through; a file error never does, since it carries
// server paths.
func (s *Server) pageFailure(err error, key string, tree mcpTree) error {
	if errors.Is(err, memory.ErrPageMissing) {
		return fmt.Errorf("there is no page %s in %s. List page paths with %s", key, tree.Name, toolListMemoryPages)
	}
	var pathErr *fs.PathError
	if !errors.As(err, &pathErr) {
		return err
	}
	s.Log.Error("mcp: writing a memory page failed", slog.Any("error", err))
	return mcpkit.ErrUnavailable
}
