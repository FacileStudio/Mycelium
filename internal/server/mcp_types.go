package server

// mcpSearchInput is what mycelium_search_memory takes.
type mcpSearchInput struct {
	Query   string `json:"query" jsonschema:"the words to look for across the wiki"`
	SpaceID string `json:"space_id,omitempty" jsonschema:"one tree to search; empty searches every tree you can use"`
	Limit   int    `json:"limit,omitempty" jsonschema:"how many results at most, 1 to 50; 0 means 10"`
}

// mcpHit is one ranked section of a page.
type mcpHit struct {
	SpaceID   string  `json:"space_id" jsonschema:"the tree the page lives in; empty is your personal tree"`
	SpaceName string  `json:"space_name"`
	Path      string  `json:"path" jsonschema:"page path relative to the wiki root, for mycelium_read_memory_page"`
	Heading   string  `json:"heading,omitempty"`
	Line      int     `json:"line" jsonschema:"line the matching section starts on"`
	Date      string  `json:"date,omitempty" jsonschema:"the finding's date when it has one"`
	Excerpt   string  `json:"excerpt"`
	Score     float64 `json:"score" jsonschema:"BM25 score, comparable only within this answer"`
}

// mcpSearchOutput is what mycelium_search_memory answers.
type mcpSearchOutput struct {
	Results []mcpHit `json:"results"`
}

// mcpPageInput is what mycelium_read_memory_page takes.
type mcpPageInput struct {
	Path     string `json:"path" jsonschema:"page path such as projects/mycelium.md, as search and the index give it"`
	SpaceID  string `json:"space_id,omitempty" jsonschema:"a space_id from mycelium_read_memory_index; empty for your personal tree"`
	FromLine int    `json:"from_line,omitempty" jsonschema:"first line to return, from next_line of a cut answer; 0 means 1"`
}

// mcpPageOutput is one page, or the window of it that fits.
type mcpPageOutput struct {
	SpaceID    string `json:"space_id"`
	SpaceName  string `json:"space_name"`
	Path       string `json:"path"`
	Content    string `json:"content" jsonschema:"the page from from_line on, frontmatter included"`
	TotalLines int    `json:"total_lines"`
	NextLine   int    `json:"next_line,omitempty" jsonschema:"set when the content was cut: pass it back as from_line"`
}

// mcpIndexInput is what mycelium_read_memory_index takes.
type mcpIndexInput struct {
	SpaceID string `json:"space_id,omitempty" jsonschema:"a space_id from spaces; empty for your personal tree"`
}

// mcpSpace is one tree the person can use.
type mcpSpace struct {
	SpaceID string `json:"space_id" jsonschema:"empty for your personal tree"`
	Name    string `json:"name"`
	Role    string `json:"role" jsonschema:"owner, admin or member"`
}

// mcpIndexOutput is the trees and the index of one of them.
type mcpIndexOutput struct {
	Spaces    []mcpSpace `json:"spaces" jsonschema:"every tree you can use, your personal tree first when you have one"`
	SpaceID   string     `json:"space_id"`
	SpaceName string     `json:"space_name,omitempty"`
	Index     string     `json:"index" jsonschema:"the tree's index.md: one line per page, by section"`
	Truncated bool       `json:"truncated,omitempty"`
	Note      string     `json:"note,omitempty"`
}

// mcpFindingInput is what mycelium_add_finding takes.
type mcpFindingInput struct {
	Page    string `json:"page" jsonschema:"an existing page such as projects/mycelium.md; findings are never filed on a new page"`
	Title   string `json:"title" jsonschema:"a short heading for the finding"`
	Source  string `json:"source" jsonschema:"where it comes from: a URL, a file path or direct observation"`
	Body    string `json:"body" jsonschema:"two to six lines of markdown (English in your personal tree): what was learned and how it changes what to do"`
	Log     string `json:"log" jsonschema:"one line for log.md saying what changed"`
	SpaceID string `json:"space_id,omitempty" jsonschema:"a space_id from mycelium_read_memory_index; empty for your personal tree"`
}

// mcpFindingOutput is what mycelium_add_finding answers.
type mcpFindingOutput struct {
	SpaceID string `json:"space_id"`
	Page    string `json:"page"`
	Indexed bool   `json:"indexed" jsonschema:"true when the page was added to index.md by this finding"`
}

// mcpArtifactListInput is what mycelium_list_artifacts takes.
type mcpArtifactListInput struct {
	SpaceID string `json:"space_id,omitempty" jsonschema:"one tree to list; empty lists every tree you can use"`
	Query   string `json:"query,omitempty" jsonschema:"keep artifacts whose title holds these words"`
	Cursor  string `json:"cursor,omitempty" jsonschema:"next_cursor of the previous page"`
}

// mcpArtifact is one artifact as the tools list it.
type mcpArtifact struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	SpaceID   string `json:"space_id" jsonschema:"empty is your personal tree"`
	SpaceName string `json:"space_name"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
	URL       string `json:"url,omitempty" jsonschema:"the web page, for artifacts of your personal tree"`
	Machine   string `json:"machine,omitempty" jsonschema:"the machine or app that recorded it, when known"`
}

// mcpArtifactListOutput is one page of artifacts.
type mcpArtifactListOutput struct {
	Artifacts  []mcpArtifact `json:"artifacts" jsonschema:"newest first"`
	Total      int           `json:"total"`
	NextCursor string        `json:"next_cursor,omitempty"`
}

// mcpArtifactSearchInput is what mycelium_search_artifacts takes.
type mcpArtifactSearchInput struct {
	Query   string `json:"query" jsonschema:"the words to look for in artifact titles and bodies, 1 to 200 characters"`
	SpaceID string `json:"space_id,omitempty" jsonschema:"one tree to search; empty searches every tree you can use"`
	Cursor  string `json:"cursor,omitempty" jsonschema:"next_cursor of the previous page, with the same query"`
}

// mcpArtifactHit is one artifact a search found.
type mcpArtifactHit struct {
	mcpArtifact
	Match   string `json:"match" jsonschema:"title when the title holds every word, else content"`
	Snippet string `json:"snippet" jsonschema:"about 160 characters of the body around the first word found"`
}

// mcpArtifactSearchOutput is one page of search hits.
type mcpArtifactSearchOutput struct {
	Results    []mcpArtifactHit `json:"results" jsonschema:"title matches first, then content matches, newest first"`
	Total      int              `json:"total"`
	NextCursor string           `json:"next_cursor,omitempty"`
}

// mcpArtifactReadInput is what mycelium_read_artifact takes.
type mcpArtifactReadInput struct {
	ID       string `json:"id" jsonschema:"an artifact id from mycelium_list_artifacts"`
	SpaceID  string `json:"space_id,omitempty" jsonschema:"the artifact's space_id; empty looks in every tree you can use"`
	FromLine int    `json:"from_line,omitempty" jsonschema:"first line to return, from next_line of a cut answer; 0 means 1"`
}

// mcpArtifactReadOutput is one artifact with its content.
type mcpArtifactReadOutput struct {
	mcpArtifact
	Content    string `json:"content" jsonschema:"the stored markdown from from_line on, frontmatter included"`
	TotalLines int    `json:"total_lines"`
	NextLine   int    `json:"next_line,omitempty" jsonschema:"set when the content was cut: pass it back as from_line"`
}

// mcpPublishInput is what mycelium_publish_artifact takes.
type mcpPublishInput struct {
	Title   string `json:"title" jsonschema:"the artifact's name; an artifact of this tree with the same title is replaced"`
	Content string `json:"content" jsonschema:"Facile Markdown, at most 1 MiB; HTML is refused"`
	SpaceID string `json:"space_id,omitempty" jsonschema:"a space_id from mycelium_read_memory_index; empty for your personal tree"`
}

// mcpPublishOutput is the artifact as published.
type mcpPublishOutput struct {
	mcpArtifact
	Replaced bool `json:"replaced" jsonschema:"true when an artifact with this title was replaced in place"`
}
