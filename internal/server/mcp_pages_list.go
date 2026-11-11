package server

import (
	"context"
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	mcpkit "github.com/FacileStudio/tronc/mcp"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// pagesPerList is how many pages one mycelium_list_memory_pages answer holds.
const pagesPerList = 200

// listMemoryPages is mycelium_list_memory_pages.
func (t *mcpTools) listMemoryPages(_ context.Context, req *sdk.CallToolRequest, in mcpListPagesInput) (*sdk.CallToolResult, mcpListPagesOutput, error) {
	tree, err := t.s.mcpOneTree(req, in.SpaceID)
	if err != nil {
		return nil, mcpListPagesOutput{}, err
	}
	offset := 0
	if err := mcpkit.DecodeCursor(in.Cursor, &offset); err != nil || offset < 0 {
		return nil, mcpListPagesOutput{}, mcpkit.ErrCursor
	}
	pages, err := wikiPages(filepath.Join(tree.Root, "memory"))
	if err != nil {
		t.s.Log.Error("mcp: listing memory pages failed", slog.Any("error", err))
		return nil, mcpListPagesOutput{}, mcpkit.ErrUnavailable
	}
	out := mcpListPagesOutput{SpaceID: tree.SpaceID, SpaceName: tree.Name, Pages: []mcpPageEntry{}, Total: len(pages)}
	if offset < len(pages) {
		end := min(offset+pagesPerList, len(pages))
		out.Pages = pages[offset:end]
		if end < len(pages) {
			out.NextCursor = mcpkit.EncodeCursor(end)
		}
	}
	return nil, out, nil
}

// wikiPages lists the markdown pages under dir by path, leaving out hidden
// files, index.md and log.md.
func wikiPages(dir string) ([]mcpPageEntry, error) {
	pages := []mcpPageEntry{}
	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if errors.Is(err, fs.ErrNotExist) && path == dir {
			return filepath.SkipDir
		}
		if err != nil || path == dir {
			return err
		}
		if strings.HasPrefix(entry.Name(), ".") {
			return skipHidden(entry)
		}
		rel, _ := filepath.Rel(dir, path)
		if key := filepath.ToSlash(rel); isListedPage(key, entry) {
			pages = append(pages, pageEntry(path, key, entry))
		}
		return nil
	})
	sort.Slice(pages, func(i, j int) bool { return pages[i].Path < pages[j].Path })
	return pages, err
}

// skipHidden leaves a hidden file out, and a hidden directory's whole tree.
func skipHidden(entry fs.DirEntry) error {
	if entry.IsDir() {
		return filepath.SkipDir
	}
	return nil
}

// isListedPage reports whether a wiki file is a page the list shows.
func isListedPage(key string, entry fs.DirEntry) bool {
	return !entry.IsDir() && strings.HasSuffix(key, ".md") && key != "index.md" && key != "log.md"
}

// pageEntry describes one page file.
func pageEntry(path, key string, entry fs.DirEntry) mcpPageEntry {
	out := mcpPageEntry{Path: key}
	if info, err := entry.Info(); err == nil {
		out.Size, out.Updated = info.Size(), info.ModTime().UTC().Format(time.RFC3339)
	}
	if data, err := os.ReadFile(path); err == nil {
		out.Title = frontmatterTitle(string(data))
	}
	return out
}

// frontmatterTitle reads a page's title: line, or "" without one.
func frontmatterTitle(content string) string {
	head, ok := strings.CutPrefix(content, "---\n")
	if !ok {
		return ""
	}
	head, _, ok = strings.Cut(head, "\n---")
	if !ok {
		return ""
	}
	for _, line := range strings.Split(head, "\n") {
		if value, found := strings.CutPrefix(line, "title:"); found {
			return strings.Trim(strings.TrimSpace(value), `"'`)
		}
	}
	return ""
}
