package artifacts

import (
	"os"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	// MaxQueryRunes caps a search query.
	MaxQueryRunes = 200
	snippetRunes  = 160
	snippetLead   = 60
)

// Where an artifact matched: every word in its title, or only across its
// title and body together.
const (
	MatchTitle   = "title"
	MatchContent = "content"
)

// Tree is one artifact tree to search: its data root and the space it is.
// An empty SpaceID is the common tree.
type Tree struct {
	Root      string
	SpaceID   string
	SpaceName string
}

// Hit is one matching artifact with the tree it lives in, where it matched
// and an excerpt of its body.
type Hit struct {
	Artifact
	SpaceID   string
	SpaceName string
	Match     string
	Snippet   string
}

// Search finds the artifacts of trees holding every word of query in their
// title or their body, frontmatter excluded. Matching is a case-insensitive
// substring test. Title matches come first, then content matches, each
// newest updated first.
func Search(trees []Tree, query string) ([]Hit, error) {
	words := strings.Fields(strings.ToLower(query))
	if len(words) == 0 {
		return nil, nil
	}
	hits := []Hit{}
	for _, tree := range trees {
		list, err := List(tree.Root)
		if err != nil {
			return nil, err
		}
		for _, art := range list {
			if hit, ok := match(art, tree, words); ok {
				hits = append(hits, hit)
			}
		}
	}
	sort.SliceStable(hits, func(i, j int) bool { return before(hits[i], hits[j]) })
	return hits, nil
}

func match(art Artifact, tree Tree, words []string) (Hit, bool) {
	raw, err := os.ReadFile(art.Path)
	if err != nil {
		return Hit{}, false
	}
	body := []rune(strings.Join(strings.Fields(plainText(stripFrontmatter(string(raw)))), " "))
	lowerBody := lowerRunes(body)
	title := strings.ToLower(art.Title)
	hit := Hit{Artifact: art, SpaceID: tree.SpaceID, SpaceName: tree.SpaceName, Match: MatchTitle}
	if !holdsAll(title, words) {
		if !holdsEach(title, lowerBody, words) {
			return Hit{}, false
		}
		hit.Match = MatchContent
	}
	hit.Snippet = snippet(body, firstOffset(lowerBody, words))
	return hit, true
}

func holdsAll(text string, words []string) bool {
	for _, word := range words {
		if !strings.Contains(text, word) {
			return false
		}
	}
	return true
}

func holdsEach(title, body string, words []string) bool {
	for _, word := range words {
		if !strings.Contains(title, word) && !strings.Contains(body, word) {
			return false
		}
	}
	return true
}

func lowerRunes(text []rune) string {
	out := make([]rune, len(text))
	for i, r := range text {
		out[i] = unicode.ToLower(r)
	}
	return string(out)
}

func firstOffset(lowerBody string, words []string) int {
	first := -1
	for _, word := range words {
		at := strings.Index(lowerBody, word)
		if at < 0 {
			continue
		}
		runes := utf8.RuneCountInString(lowerBody[:at])
		if first < 0 || runes < first {
			first = runes
		}
	}
	return max(first, 0)
}

func snippet(body []rune, at int) string {
	start := max(at-snippetLead, 0)
	end := min(start+snippetRunes, len(body))
	out := strings.TrimSpace(string(body[start:end]))
	if start > 0 {
		out = "…" + out
	}
	if end < len(body) {
		out += "…"
	}
	return out
}

func before(a, b Hit) bool {
	if a.Match != b.Match {
		return a.Match == MatchTitle
	}
	if !a.Updated.Equal(b.Updated) {
		return a.Updated.After(b.Updated)
	}
	return a.ID < b.ID
}
