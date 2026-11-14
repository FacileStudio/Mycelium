package artifacts

import (
	"strings"
	"unicode"
)

var inlineMarks = strings.NewReplacer("**", "", "__", "", "==", "", "`", "", "|", " ")

// plainText drops the Markdown syntax a reader should not see in a snippet:
// fence lines, table separator rows, heading, quote and list markers, task
// boxes and inline emphasis, code and highlight marks.
func plainText(markdown string) string {
	lines := strings.Split(markdown, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") || strings.Trim(trimmed, "|-: ") == "" {
			continue
		}
		out = append(out, inlineMarks.Replace(stripLineMarker(trimmed)))
	}
	return strings.Join(out, "\n")
}

func stripLineMarker(line string) string {
	line = strings.TrimLeft(line, "#> ")
	for _, marker := range []string{"- [ ] ", "- [x] ", "- ", "* ", "+ "} {
		if rest, ok := strings.CutPrefix(line, marker); ok {
			return rest
		}
	}
	digits := strings.TrimLeftFunc(line, unicode.IsDigit)
	if rest, ok := strings.CutPrefix(digits, ". "); ok && len(digits) < len(line) {
		return rest
	}
	return line
}
