package server

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	mcpkit "github.com/FacileStudio/tronc/mcp"
)

// mcpTextBudget caps the text one read returns, well under tronc's result
// budget once the JSON around it is counted.
const mcpTextBudget = 60_000

// errMCPPagePath is the refusal for a page path outside the wiki.
var errMCPPagePath = errors.New("the path must be a page inside the wiki, such as projects/mycelium.md")

// mcpPageKey cleans a page path the way the wiki names pages: relative to
// the wiki root, slash separated, ending in .md, with no dot segment.
func mcpPageKey(raw string) (string, error) {
	key := strings.TrimPrefix(path.Clean("/"+strings.TrimSpace(raw)), "/")
	if key == "" || key == "." {
		return "", errMCPPagePath
	}
	for _, segment := range strings.Split(key, "/") {
		if strings.HasPrefix(segment, ".") {
			return "", errMCPPagePath
		}
	}
	if !strings.HasSuffix(key, ".md") {
		key += ".md"
	}
	return key, nil
}

// readMCPPage reads one page of a tree's wiki.
func readMCPPage(tree mcpTree, raw string) (string, string, error) {
	key, err := mcpPageKey(raw)
	if err != nil {
		return "", "", err
	}
	data, err := os.ReadFile(filepath.Join(tree.Root, "memory", filepath.FromSlash(key)))
	if errors.Is(err, os.ErrNotExist) {
		return "", "", fmt.Errorf("there is no page %s in %s. Find page paths with %s or %s",
			key, tree.Name, toolSearchMemory, toolReadMemoryIndex)
	}
	if err != nil {
		return "", "", mcpkit.ErrUnavailable
	}
	return key, string(data), nil
}

// lineWindow returns text from line from (1-based) on, cut on a line
// boundary to the read budget, with the total line count and the line to
// continue from when it was cut.
func lineWindow(text string, from int) (string, int, int) {
	lines := strings.SplitAfter(text, "\n")
	if n := len(lines); n > 0 && lines[n-1] == "" {
		lines = lines[:n-1]
	}
	from = max(from, 1)
	var out strings.Builder
	for i := from - 1; i < len(lines); i++ {
		if out.Len()+len(lines[i]) <= mcpTextBudget {
			out.WriteString(lines[i])
			continue
		}
		if out.Len() == 0 {
			clipped, _ := mcpkit.Clip(lines[i], mcpTextBudget)
			out.WriteString(clipped)
			i++
		}
		return out.String(), len(lines), nextLine(i, len(lines))
	}
	return out.String(), len(lines), 0
}

// nextLine is the 1-based line to continue from at 0-based index i, or 0
// when nothing is left.
func nextLine(i, total int) int {
	if i >= total {
		return 0
	}
	return i + 1
}
