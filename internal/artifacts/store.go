package artifacts

import (
	"cmp"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Dir returns the path to the artifacts directory.
func Dir(dataDir string) string { return filepath.Join(dataDir, Dirname) }

func allDirs(dataDir string) []string {
	dirs := []string{filepath.Join(dataDir, "artifacts")}
	repDir := filepath.Join(dataDir, "reports")
	if _, err := os.Stat(repDir); err == nil {
		dirs = append(dirs, repDir)
	}
	return dirs
}

// applyField sets one metadata field read from a file. Unknown keys, such as
// the expires line older artifacts still carry, are ignored: artifacts are
// permanent and that line no longer means anything.
func applyField(art *Artifact, key, value string) {
	switch key {
	case "title":
		art.Title = value
	case "machine":
		art.Machine = value
	case "created":
		art.Created, _ = time.Parse(timeFormat, value)
	case "updated":
		art.Updated, _ = time.Parse(timeFormat, value)
	}
}

func parseHTMLMeta(raw []byte, art *Artifact) {
	for _, m := range metaPair.FindAllSubmatch(raw, -1) {
		applyField(art, strings.ToLower(string(m[1])), html.UnescapeString(string(m[2])))
	}
}

func parseFrontmatter(content string, art *Artifact) {
	if !strings.HasPrefix(content, "---\n") && !strings.HasPrefix(content, "---\r\n") {
		return
	}
	end := strings.Index(content[4:], "\n---")
	if end < 0 {
		return
	}
	head := content[4 : 4+end]
	for _, line := range strings.Split(head, "\n") {
		line = strings.TrimSpace(line)
		colon := strings.Index(line, ":")
		if colon < 0 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(line[:colon]))
		applyField(art, key, strings.Trim(strings.TrimSpace(line[colon+1:]), `"'`))
	}
}

func parse(path string, raw []byte) Artifact {
	ext := strings.ToLower(filepath.Ext(path))
	id := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	art := Artifact{ID: id, Path: path}

	if ext == ".html" || ext == ".htm" {
		parseHTMLMeta(raw, &art)
	} else {
		parseFrontmatter(string(raw), &art)
	}

	if art.Title == "" {
		art.Title = titleFrom(raw, path)
	}
	if art.Created.IsZero() {
		if info, err := os.Stat(path); err == nil {
			art.Created = info.ModTime()
		} else {
			art.Created = time.Now()
		}
	}
	if art.Updated.IsZero() {
		art.Updated = art.Created
	}
	return art
}

// List returns all stored artifacts, sorted newest first.
func List(dataDir string) ([]Artifact, error) {
	seen := map[string]bool{}
	var out []Artifact

	for _, root := range allDirs(dataDir) {
		walkArtifacts(root, func(path, id string) bool {
			if seen[id] {
				return false
			}
			raw, err := os.ReadFile(path)
			if err != nil {
				return false
			}
			seen[id] = true
			out = append(out, parse(path, raw))
			return false
		})
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Created.After(out[j].Created) })
	return out, nil
}

var hashSuffix = regexp.MustCompile(`-[a-f0-9]{6}$`)

func matchesID(base, target, cleanSlug string) bool {
	if strings.EqualFold(base, target) || strings.EqualFold(base, cleanSlug) {
		return true
	}
	if strings.HasSuffix(base, "-"+cleanSlug) || strings.HasPrefix(base, cleanSlug+"-") {
		return true
	}
	strippedDate := datePrefix.ReplaceAllString(base, "")
	strippedAll := hashSuffix.ReplaceAllString(strippedDate, "")
	if strings.EqualFold(strippedAll, cleanSlug) || strings.EqualFold(strippedDate, cleanSlug) {
		return true
	}
	return strings.Contains(base, "-"+cleanSlug+"-") || strings.Contains(base, "-"+cleanSlug)
}

// walkArtifacts calls visit with the path and id of every artifact file under
// root, and stops as soon as visit returns true.
func walkArtifacts(root string, visit func(path, base string) bool) {
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d == nil || d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(d.Name()))
		if ext != ".md" && ext != ".html" && ext != ".htm" {
			return nil
		}
		if visit(path, strings.TrimSuffix(d.Name(), filepath.Ext(d.Name()))) {
			return filepath.SkipAll
		}
		return nil
	})
}

// findInRoot walks one root for the artifact named target. An exact id wins
// over a fuzzy slug match, so a lookup by full id never lands on a sibling
// whose name merely contains it.
func findInRoot(root, target, cleanSlug string) (string, []byte) {
	var exact, fuzzy string
	walkArtifacts(root, func(path, base string) bool {
		if strings.EqualFold(base, target) {
			exact = path
			return true
		}
		if fuzzy == "" && matchesID(base, target, cleanSlug) {
			fuzzy = path
		}
		return false
	})
	matched := cmp.Or(exact, fuzzy)
	if matched == "" {
		return "", nil
	}
	raw, err := os.ReadFile(matched)
	if err != nil {
		return "", nil
	}
	return matched, raw
}

// Find retrieves a single artifact by its ID or title slug.
func Find(dataDir, id string) (Artifact, error) {
	target := strings.TrimSpace(id)
	cleanSlug := Slug(target)

	for _, root := range allDirs(dataDir) {
		path, raw := findInRoot(root, target, cleanSlug)
		if path != "" {
			return parse(path, raw), nil
		}
	}

	return Artifact{}, fmt.Errorf("no artifact named %q", id)
}

// Remove deletes an artifact by its ID.
func Remove(dataDir, id string) error {
	art, err := Find(dataDir, id)
	if err != nil {
		return err
	}
	if err := os.Remove(art.Path); err != nil {
		return fmt.Errorf("failed to delete %s: %w", art.Path, err)
	}
	return nil
}

// Update rewrites an artifact in place. The id, path, machine and created time
// are kept; the body is replaced, the title changes when one is given, and
// updated is set to now. Only an exact id is accepted.
func Update(dataDir, id, title string, raw []byte, now time.Time) (Artifact, error) {
	art, err := Find(dataDir, id)
	if err != nil || !strings.EqualFold(art.ID, strings.TrimSpace(id)) {
		return Artifact{}, fmt.Errorf("no artifact named %q", id)
	}
	if t := strings.TrimSpace(title); t != "" {
		art.Title = t
	}
	art.Updated = stored(now)
	ext := strings.ToLower(filepath.Ext(art.Path))
	stamped := stamp(raw, art, ext == ".html" || ext == ".htm")
	if err := os.WriteFile(art.Path, stamped, 0o644); err != nil {
		return Artifact{}, fmt.Errorf("failed to write %s: %w", art.Path, err)
	}
	return art, nil
}
