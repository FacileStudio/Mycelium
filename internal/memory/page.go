package memory

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

// MaxPageBytes caps one page written whole, frontmatter included.
const MaxPageBytes = 256 << 10

// pageTypes are the frontmatter types a page written whole may declare. A
// standard is missing on purpose: normative pages are written by a human.
var pageTypes = []string{"bug", "tool", "project", "convention", "synthesis", "actor"}

// Log operations a whole-page write or a delete records.
const (
	opCreatePage = "create-page"
	opLint       = "lint"
)

// ErrPageMissing is the refusal for a page that is not in the wiki.
var ErrPageMissing = errors.New("there is no such page")

// Page is one page written whole: its key, its full markdown with the
// frontmatter, the index hook for a page the index does not name yet, and
// the log line.
type Page struct {
	Key     string
	Content string
	Summary string
	Log     string
}

// PageResult reports what one WritePage did.
type PageResult struct {
	Page    string
	Created bool
	Indexed bool
}

// Protected refuses the pages no whole-page write or delete may touch:
// index.md and log.md, which the tools maintain themselves, and every page
// under standards/, which a human ratifies.
func Protected(key string) error {
	switch {
	case key == indexFile:
		return fmt.Errorf("index.md is maintained by the wiki tools: write or delete a page and its index line follows")
	case key == logFile:
		return fmt.Errorf("log.md is append-only history: its entries are never written or deleted through this tool")
	case strings.HasPrefix(key, "standards/"):
		return fmt.Errorf("%s is a normative page: standards are changed by a human and ratified, never through this tool", key)
	}
	return nil
}

// WritePage creates or replaces one page and does the bookkeeping a finding
// gets: the `updated:` stamp, `created:` kept from the old page or set
// today, an index pointer when the index does not name the page, and the log
// line. The page must carry YAML frontmatter with a title and a descriptive
// type, and be English, like everything else the wiki holds.
func WritePage(dataDir string, p Page, now time.Time) (PageResult, error) {
	path, key, err := pagePath(dataDir, p.Key)
	if err != nil {
		return PageResult{}, err
	}
	if err := Protected(key); err != nil {
		return PageResult{}, err
	}
	if err := checkPage(p); err != nil {
		return PageResult{}, err
	}
	day := now.Format(dayLayout)
	release, err := lockWiki(dataDir)
	if err != nil {
		return PageResult{}, err
	}
	defer release()
	old, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return PageResult{}, err
	}
	created := err != nil
	if !created && isStandard(string(old)) {
		return PageResult{}, Protected("standards/" + key)
	}
	content, err := stampPage(p.Content, string(old), day)
	if err != nil {
		return PageResult{}, err
	}
	p.Content = content
	return commitPage(dataDir, pageWrite{path: path, key: key, day: day, created: created}, p)
}

// pageWrite is where a checked page lands and what kind of write it is.
type pageWrite struct {
	path, key, day string
	created        bool
}

// commitPage writes the page, the index and the log together.
func commitPage(dataDir string, w pageWrite, p Page) (PageResult, error) {
	path, key, content, day, created := w.path, w.key, p.Content, w.day, w.created
	front, _ := frontmatter(content)
	index := filepath.Join(memoryRoot(dataDir), indexFile)
	history := filepath.Join(memoryRoot(dataDir), logFile)
	hook := strings.TrimSpace(p.Summary)
	if hook == "" {
		hook = strings.TrimSpace(p.Log)
	}
	nextIndex, pointed := indexPointer(readOrEmpty(index), key, front.title, hook)
	op := opLint
	if created {
		op = opCreatePage
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return PageResult{}, err
	}
	err := writeAll([]edit{
		{path: path, data: strings.TrimRight(content, "\n") + "\n"},
		{path: index, data: nextIndex},
		{path: history, data: appendLogAs(readOrEmpty(history), day, op, p.Log)},
	})
	if err != nil {
		return PageResult{}, err
	}
	return PageResult{Page: key, Created: created, Indexed: pointed}, nil
}

// CheckDelete runs every refusal DeletePage would, without touching the
// wiki, so a person is never asked to confirm a delete that would fail. It
// returns the page's key and its current bytes.
func CheckDelete(dataDir, page, log string) (string, []byte, error) {
	path, key, err := pagePath(dataDir, page)
	if err != nil {
		return "", nil, err
	}
	if err := Protected(key); err != nil {
		return "", nil, err
	}
	if strings.TrimSpace(log) == "" {
		return "", nil, fmt.Errorf("log is empty: pass one line for log.md saying why the page goes")
	}
	if err := refuseFrench(log); err != nil {
		return "", nil, err
	}
	old, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil, ErrPageMissing
	}
	if err != nil {
		return "", nil, err
	}
	if isStandard(string(old)) {
		return "", nil, Protected("standards/" + key)
	}
	return key, old, nil
}

// DeletePage removes one page, drops its index line and logs why.
func DeletePage(dataDir, page, log string, now time.Time) (string, error) {
	release, err := lockWiki(dataDir)
	if err != nil {
		return "", err
	}
	defer release()
	key, _, err := CheckDelete(dataDir, page, log)
	if err != nil {
		return "", err
	}
	index := filepath.Join(memoryRoot(dataDir), indexFile)
	history := filepath.Join(memoryRoot(dataDir), logFile)
	day := now.Format(dayLayout)
	edits := []edit{{path: history, data: appendLogAs(readOrEmpty(history), day, opLint, "deleted "+key+": "+log)}}
	if current := readOrEmpty(index); current != "" {
		edits = append(edits, edit{path: index, data: dropPointer(current, key)})
	}
	if err := writeAll(edits); err != nil {
		return "", err
	}
	return key, os.Remove(filepath.Join(memoryRoot(dataDir), filepath.FromSlash(key)))
}

// checkPage refuses a page that breaks the wiki's rules before anything is
// locked: empty or too large, no frontmatter, no title, a type that is not a
// descriptive one, no log line, or French prose.
func checkPage(p Page) error {
	switch {
	case strings.TrimSpace(p.Content) == "":
		return fmt.Errorf("content is empty: pass the whole page, frontmatter included")
	case len(p.Content) > MaxPageBytes:
		return fmt.Errorf("content exceeds %d KiB: split the page", MaxPageBytes>>10)
	case strings.TrimSpace(p.Log) == "":
		return fmt.Errorf("log is empty: pass one line for log.md saying what changed")
	}
	if !strings.HasPrefix(p.Content, "---\n") || !strings.Contains(p.Content[4:], "\n---") {
		return fmt.Errorf("the page needs YAML frontmatter between --- lines, with at least a title and a type field")
	}
	front, _ := frontmatter(p.Content)
	if strings.TrimSpace(front.title) == "" {
		return fmt.Errorf("the frontmatter needs a title field")
	}
	if !slices.Contains(pageTypes, front.kind) {
		return fmt.Errorf("the frontmatter type: must be one of %s; standards are never written through this tool",
			strings.Join(pageTypes, ", "))
	}
	return refuseFrench(p.Content + "\n" + p.Log + "\n" + p.Summary)
}

// isStandard reports whether a page declares itself normative.
func isStandard(content string) bool {
	front, _ := frontmatter(content)
	return front.kind == normativeType
}
