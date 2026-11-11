package server

// mcpWritePageInput is what mycelium_write_memory_page takes.
type mcpWritePageInput struct {
	Path    string `json:"path" jsonschema:"page path such as projects/mycelium.md; a missing page is created"`
	Content string `json:"content" jsonschema:"the whole page in markdown (English in your personal tree, any language in a space), starting with YAML frontmatter that holds at least title: and type: (bug, tool, project, convention, synthesis or actor); updated: and created: are stamped for you"`
	Summary string `json:"summary,omitempty" jsonschema:"one line describing the page, for its index.md entry when the index does not name it yet; empty uses log"`
	Log     string `json:"log" jsonschema:"one line for log.md saying what changed"`
	SpaceID string `json:"space_id,omitempty" jsonschema:"a space_id from mycelium_read_memory_index; empty for your personal tree"`
}

// mcpWritePageOutput is what mycelium_write_memory_page answers.
type mcpWritePageOutput struct {
	SpaceID string `json:"space_id"`
	Path    string `json:"path"`
	Created bool   `json:"created" jsonschema:"true when the page did not exist before"`
	Indexed bool   `json:"indexed" jsonschema:"true when index.md gained a line for the page"`
}

// mcpDeletePageInput is what mycelium_delete_memory_page takes.
type mcpDeletePageInput struct {
	Path    string `json:"path" jsonschema:"the page to delete, such as people/alice.md"`
	Log     string `json:"log" jsonschema:"one line for log.md saying why the page goes"`
	SpaceID string `json:"space_id,omitempty" jsonschema:"a space_id from mycelium_read_memory_index; empty for your personal tree"`
}

// mcpDeletePageOutput is where a delete stands.
type mcpDeletePageOutput struct {
	Status     string `json:"status" jsonschema:"pending_confirmation, declined or deleted"`
	SpaceID    string `json:"space_id,omitempty"`
	Path       string `json:"path,omitempty"`
	ConfirmURL string `json:"confirm_url,omitempty" jsonschema:"the page where the user confirms, when status is pending_confirmation"`
	Message    string `json:"message"`
}

// mcpListPagesInput is what mycelium_list_memory_pages takes.
type mcpListPagesInput struct {
	SpaceID string `json:"space_id,omitempty" jsonschema:"a space_id from mycelium_read_memory_index; empty for your personal tree"`
	Cursor  string `json:"cursor,omitempty" jsonschema:"next_cursor of the previous page"`
}

// mcpPageEntry is one page of a wiki as the list shows it.
type mcpPageEntry struct {
	Path    string `json:"path"`
	Title   string `json:"title,omitempty"`
	Size    int64  `json:"size"`
	Updated string `json:"updated" jsonschema:"when the file last changed, RFC 3339"`
}

// mcpListPagesOutput is one page of a wiki's pages, by path.
type mcpListPagesOutput struct {
	SpaceID    string         `json:"space_id"`
	SpaceName  string         `json:"space_name"`
	Pages      []mcpPageEntry `json:"pages"`
	Total      int            `json:"total"`
	NextCursor string         `json:"next_cursor,omitempty"`
}

// mcpCreateSpaceInput is what mycelium_create_space takes.
type mcpCreateSpaceInput struct {
	Name        string `json:"name" jsonschema:"the space's name; a space of yours with this name is reused"`
	Description string `json:"description,omitempty" jsonschema:"what the space is for"`
}

// mcpCreateSpaceOutput is the space as created or found.
type mcpCreateSpaceOutput struct {
	SpaceID string `json:"space_id"`
	Name    string `json:"name"`
	Role    string `json:"role"`
	Created bool   `json:"created" jsonschema:"false when a space of yours with this name already existed"`
}

// mcpUpdateArtifactInput is what mycelium_update_artifact takes.
type mcpUpdateArtifactInput struct {
	ID      string `json:"id" jsonschema:"the artifact's exact id, from mycelium_list_artifacts"`
	Title   string `json:"title,omitempty" jsonschema:"a new title; empty keeps the current one"`
	Content string `json:"content" jsonschema:"the whole new Facile Markdown, at most 1 MiB; it replaces the old content"`
	SpaceID string `json:"space_id,omitempty" jsonschema:"the artifact's space_id; empty for your personal tree"`
}

// mcpDeleteArtifactInput is what mycelium_delete_artifact takes.
type mcpDeleteArtifactInput struct {
	ID      string `json:"id" jsonschema:"the artifact's exact id, from mycelium_list_artifacts"`
	SpaceID string `json:"space_id,omitempty" jsonschema:"the artifact's space_id; empty for your personal tree"`
}

// mcpDeleteArtifactOutput is where an artifact delete stands.
type mcpDeleteArtifactOutput struct {
	Status     string `json:"status" jsonschema:"pending_confirmation, declined or deleted"`
	SpaceID    string `json:"space_id,omitempty"`
	ID         string `json:"id,omitempty"`
	ConfirmURL string `json:"confirm_url,omitempty" jsonschema:"the page where the user confirms, when status is pending_confirmation"`
	Message    string `json:"message"`
}
