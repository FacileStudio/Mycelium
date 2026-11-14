package consolidate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const dayLayout = "2006-01-02"

// Writer applies consolidation decisions to the wiki under MemoryPath,
// following the prose conventions the pages themselves use: a `### <short
// title>` block opening with **Date**: and **Source**: lines for CREATE, and
// ~~struck text~~ plus a [SUPERSEDED by: ...] marker with the correction
// beneath it for SUPERSEDE. Nothing is ever deleted.
type Writer struct {
	MemoryPath string
}

// Create appends the candidate as a finding block. A decision carrying a
// match appends to that page; one without starts a new page under the
// top-level dir its text classifies into. Returns the path written.
func (w *Writer) Create(dec Decision, c Candidate, now time.Time) (string, error) {
	block := findingBlock(shortTitle(c.Text), c.Text, sourceLine(c), now)
	if dec.Match == nil {
		return w.createPage(c, block, now)
	}
	return appendBlock(filepath.Join(w.MemoryPath, filepath.FromSlash(dec.Match.Path)), block)
}

// appendBlock adds a finding to the end of an existing page, leaving one blank
// line between it and whatever was already there.
func appendBlock(path, block string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return path, writeFileAtomic(path, []byte(strings.TrimRight(string(data), "\n")+"\n\n"+block))
}

// createPage starts a page for a candidate that matched nothing. A slug can
// collide with a page that already exists — two unrelated findings whose first
// line kebabs the same way — and overwriting it would delete a claim, which is
// the one thing this stage must never do. A collision appends instead.
func (w *Writer) createPage(c Candidate, block string, now time.Time) (string, error) {
	title := shortTitle(c.Text)
	dir := classifyDir(c.Text)
	pagePath := filepath.Join(w.MemoryPath, dir, slug(title)+".md")
	if _, err := os.Stat(pagePath); err == nil {
		return appendBlock(pagePath, block)
	} else if !os.IsNotExist(err) {
		return "", err
	}
	content := fmt.Sprintf(`---
title: %s
type: %s
sources:
  - consolidation
related: []
confidence: low
created: %s
updated: %s
---

# %s

%s`, title, dir, now.Format(dayLayout), now.Format(dayLayout), title, block)
	if err := os.MkdirAll(filepath.Dir(pagePath), 0o755); err != nil {
		return "", err
	}
	return pagePath, writeFileAtomic(pagePath, []byte(content))
}

func findingBlock(title, text, source string, now time.Time) string {
	return fmt.Sprintf("### %s\n%s\n%s\n", title, provenanceLines(source, now), text)
}

func provenanceLines(source string, now time.Time) string {
	return fmt.Sprintf("**Date**: %s\n**Source**: %s", now.Format(dayLayout), source)
}

func sourceLine(c Candidate) string {
	return "consolidation, " + strings.Join(c.EpisodeRefs, "; ")
}

// classifyDir picks the top-level wiki dir for a new page from the text's own
// vocabulary: failures go to bugs/, tooling vocabulary to tools/, people words
// to people/, everything else to projects/.
func classifyDir(text string) string {
	lowered := strings.ToLower(text)
	failureWords := []string{"error", "fail", "bug", "crash", "panic", "refused", "the fix was", "gotcha"}
	toolWords := []string{"cli", "command", "flag", ".go", ".ts", ".sh", "tool", "script", "binary"}
	personWords := []string{"collaborator", "stakeholder", "contractor", "team member"}
	if containsAny(lowered, failureWords) {
		return "bugs"
	}
	if containsAny(lowered, toolWords) {
		return "tools"
	}
	if containsAny(lowered, personWords) {
		return "people"
	}
	return "projects"
}

func containsAny(s string, needles []string) bool {
	for _, n := range needles {
		if strings.Contains(s, n) {
			return true
		}
	}
	return false
}

// shortTitle condenses a finding's first line into the heading style the wiki
// uses: a short declarative sentence, cut on a word boundary when needed.
func shortTitle(text string) string {
	line := strings.TrimSpace(strings.SplitN(strings.TrimSpace(text), "\n", 2)[0])
	line = strings.TrimPrefix(line, "- ")
	line = strings.Trim(line, "`*# ")
	runes := []rune(line)
	if len(runes) <= 70 {
		return line
	}
	cut := 69
	for cut > 40 && runes[cut] != ' ' {
		cut--
	}
	return strings.TrimRight(string(runes[:cut]), ",;:")
}

// slug turns a title into a filename-safe lowercase kebab string.
func slug(title string) string {
	var b strings.Builder
	lastDash := false
	for _, r := range strings.ToLower(title) {
		switch {
		case r >= 'a' && r <= 'z' || r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		default:
			if !lastDash && b.Len() > 0 {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}
	out := strings.Trim(b.String(), "-")
	if len(out) > 48 {
		out = out[:48]
	}
	return strings.Trim(out, "-")
}

func writeFileAtomic(path string, data []byte) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
