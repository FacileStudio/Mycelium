package artifacts

import (
	"fmt"
	"regexp"
	"strings"
)

var mdLink = regexp.MustCompile(`!?\[.*?\]\(([^)]+)\)`)

// stamp writes the metadata frontmatter above the body. The body is kept
// byte for byte, so fenced blocks of any kind survive untouched.
func stamp(raw []byte, art Artifact) []byte {
	return []byte(frontmatterBlock(art) + stripFrontmatter(string(raw)))
}

func stripFrontmatter(content string) string {
	if !strings.HasPrefix(content, "---\n") && !strings.HasPrefix(content, "---\r\n") {
		return content
	}
	end := strings.Index(content[4:], "\n---")
	if end < 0 {
		return content
	}
	rest := content[4+end+1:]
	cut := strings.Index(rest, "\n")
	if cut < 0 {
		return ""
	}
	return strings.TrimLeft(rest[cut+1:], "\r\n")
}

func frontmatterBlock(art Artifact) string {
	var b strings.Builder
	b.WriteString("---\n")
	fmt.Fprintf(&b, "title: %q\n", art.Title)
	if art.Machine != "" {
		fmt.Fprintf(&b, "machine: %q\n", art.Machine)
	}
	fmt.Fprintf(&b, "created: %s\n", art.Created.UTC().Format(timeFormat))
	fmt.Fprintf(&b, "updated: %s\n", art.Updated.UTC().Format(timeFormat))
	b.WriteString("type: artifact\n")
	b.WriteString("---\n\n")
	return b.String()
}

var schemeRef = regexp.MustCompile(`(?i)^[a-z][a-z0-9+.-]*:`)

// ExternalRefs scans raw content for relative links and assets that will not resolve from disk.
func ExternalRefs(raw []byte) []string {
	seen := map[string]bool{}
	var out []string
	for _, m := range mdLink.FindAllSubmatch(raw, -1) {
		ref := strings.TrimSpace(string(m[1]))
		if ref == "" || seen[ref] || strings.HasPrefix(ref, "#") || schemeRef.MatchString(ref) {
			continue
		}
		seen[ref] = true
		out = append(out, ref)
	}
	return out
}
