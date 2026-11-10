package artifacts

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/FacileStudio/Mycelium/internal/encre"
)

const Dirname = "artifacts"

const timeFormat = time.RFC3339

// Artifact represents a recorded markdown document. Artifacts are
// permanent: nothing deletes one except an explicit remove.
type Artifact struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Path    string    `json:"path"`
	Machine string    `json:"machine"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

// Report is a backward-compatible alias for Artifact.
type Report = Artifact

// ErrNotMarkdown is returned for HTML content: artifacts are markdown only.
var ErrNotMarkdown = errors.New("artifacts are markdown only, HTML is not supported: write the document as markdown (.md)")

// Request holds parameters for adding a new artifact.
type Request struct {
	Source  string
	Title   string
	Machine string
}

// AddContent registers a new artifact directly from byte content. The body
// must be a valid Facile document: encre's refusal comes back as a
// *encre.ValidationError and nothing is written.
func AddContent(dataDir string, raw []byte, filename string, req Request, now time.Time) (Artifact, error) {
	if IsHTML(filename, raw) {
		return Artifact{}, ErrNotMarkdown
	}
	if err := validateBody(raw); err != nil {
		return Artifact{}, err
	}
	title := strings.TrimSpace(req.Title)
	if title == "" {
		title = titleFrom(raw, filename)
	}
	id := idFor(title, req.Machine, now)
	if id == "" {
		return Artifact{}, fmt.Errorf("cannot derive a name from %q, pass --title", filename)
	}
	targetDir := filepath.Join(Dir(dataDir), now.UTC().Format("2006/01"))
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return Artifact{}, fmt.Errorf("failed to create target directory: %w", err)
	}
	art := Artifact{
		ID:      id,
		Title:   title,
		Path:    filepath.Join(targetDir, id+".md"),
		Machine: req.Machine,
		Created: stored(now),
		Updated: stored(now),
	}
	stamped := stamp(raw, art)
	if err := os.WriteFile(art.Path, stamped, 0o644); err != nil {
		return Artifact{}, fmt.Errorf("failed to write %s: %w", art.Path, err)
	}
	return art, nil
}

// Add copies and registers a new artifact from a file path in the local store.
func Add(dataDir string, req Request, now time.Time) (Artifact, error) {
	raw, err := os.ReadFile(req.Source)
	if err != nil {
		return Artifact{}, fmt.Errorf("failed to read %s: %w", req.Source, err)
	}
	return AddContent(dataDir, raw, req.Source, req, now)
}

// validateBody asks encre whether the body under any frontmatter is a valid
// document. Without encre on this machine the body is stored as it came and
// the viewer renders it in degraded mode.
func validateBody(raw []byte) error {
	err := encre.ValidateMarkdown(context.Background(), []byte(stripFrontmatter(string(raw))))
	if errors.Is(err, encre.ErrUnavailable) {
		return nil
	}
	return err
}

// stored is t as the file records it: UTC, to the second. Returning anything
// finer would make a fresh artifact differ from the same one read back.
func stored(t time.Time) time.Time { return t.UTC().Truncate(time.Second) }

// URL returns the canonical web URL for an artifact given the server URL and artifact ID.
func URL(serverURL, id string) string {
	clean := strings.TrimSpace(serverURL)
	if clean == "" {
		return ""
	}
	return fmt.Sprintf("%s/artifacts/%s", strings.TrimRight(clean, "/"), id)
}

var (
	mdHeading    = regexp.MustCompile(`(?m)^#\s+(.+)$`)
	slugTrim     = regexp.MustCompile(`[^a-z0-9]+`)
	datePrefix   = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}-`)
	fullHashedID = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}-.+-[a-f0-9]{6}$`)
)

func idFor(title, machine string, now time.Time) string {
	slug := Slug(title)
	if slug == "" {
		return ""
	}
	if fullHashedID.MatchString(slug) {
		return slug
	}
	baseSlug := datePrefix.ReplaceAllString(slug, "")
	h := sha256.Sum256([]byte(machine + ":" + title + ":" + now.UTC().Format(time.RFC3339Nano)))
	hashStr := fmt.Sprintf("%x", h[:3])
	return now.UTC().Format("2006-01-02") + "-" + baseSlug + "-" + hashStr
}

// IsHTML reports whether a source is an HTML document, by its extension or
// by an opening doctype or html tag.
func IsHTML(source string, raw []byte) bool {
	ext := strings.ToLower(filepath.Ext(source))
	if ext == ".html" || ext == ".htm" {
		return true
	}
	head := bytes.ToLower(bytes.TrimSpace(raw))
	return bytes.HasPrefix(head, []byte("<!doctype")) || bytes.HasPrefix(head, []byte("<html"))
}

func titleFrom(raw []byte, source string) string {
	text := string(raw)
	if title := extractFrontmatterTitle(text); title != "" {
		return title
	}
	if m := mdHeading.FindStringSubmatch(text); len(m) > 1 {
		if t := strings.TrimSpace(m[1]); t != "" {
			return t
		}
	}
	base := strings.TrimSuffix(filepath.Base(source), filepath.Ext(source))
	return datePrefix.ReplaceAllString(base, "")
}

func extractFrontmatterTitle(content string) string {
	if !strings.HasPrefix(content, "---\n") && !strings.HasPrefix(content, "---\r\n") {
		return ""
	}
	end := strings.Index(content[4:], "\n---")
	if end < 0 {
		return ""
	}
	head := content[4 : 4+end]
	for _, line := range strings.Split(head, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "title:") {
			val := strings.TrimSpace(strings.TrimPrefix(trimmed, "title:"))
			return strings.Trim(val, `"'`)
		}
	}
	return ""
}

// Slug sanitizes a title into a URL and filename friendly slug.
func Slug(title string) string {
	return strings.Trim(slugTrim.ReplaceAllString(strings.ToLower(title), "-"), "-")
}
