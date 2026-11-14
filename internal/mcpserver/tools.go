package mcpserver

import (
	"strings"

	"github.com/FacileStudio/Mycelium/internal/flow"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// searchMemoryTool declares the wiki search.
//
// ReadOnlyHint puts it in the client's read-only group, which a user approves
// once in bulk instead of confirming on every call. OpenWorldHint is false: it
// reads this machine's own wiki and nothing else.
func searchMemoryTool() *mcp.Tool {
	return &mcp.Tool{
		Name:  "search_memory",
		Title: "Search memory",
		Description: "Search the shared agent wiki, best match first. It ranks the findings of this " +
			"machine's synced markdown folder by keyword (BM25, prefix match, accents folded). Read this " +
			"before non-trivial work rather than rediscovering what is already written down.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    true,
			DestructiveHint: hint(false),
			IdempotentHint:  true,
			OpenWorldHint:   hint(false),
		},
	}
}

// listFlowsTool declares the flow inventory.
//
// It reads two files per flow and changes neither, so it carries the same
// read-only annotations as the search. Trust comes back as a field rather than
// a rendered column: a model deciding whether to call run_flow should read a
// value, not parse a table built for a terminal.
func listFlowsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:  "list_flows",
		Title: "List flows",
		Description: "List the recorded procedures on this machine with their step count and trust " +
			"state. A flow runs only when a human on this machine has pinned its exact content, so " +
			"trust is what decides whether run_flow executes a flow or refuses it.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    true,
			DestructiveHint: hint(false),
			IdempotentHint:  true,
			OpenWorldHint:   hint(false),
		},
	}
}

// runFlowTool declares the one tool here that executes anything.
//
// It takes a flow name and nothing else. Command injection is the dominant way
// MCP servers are abused, and a fixed set of procedures a human read and pinned
// is this server's whole defence against it; an args or command passthrough
// would hand that back. DestructiveHint and OpenWorldHint are both true because
// a flow's steps are arbitrary shell commands that may reach anything.
//
// Those annotations still only shape a prompt. The spec requires clients to
// treat them as untrusted, so the pin checked in trust.go is the real gate.
func runFlowTool() *mcp.Tool {
	return &mcp.Tool{
		Name:  "run_flow",
		Title: "Run a flow",
		Description: "Run a recorded procedure by name and return each step's exit code, output and " +
			"the path of the run record. Takes a flow name and nothing else: no command or argument " +
			"can be passed, and a flow no human has pinned on this machine is refused with the " +
			"command that fixes it.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: hint(true),
			IdempotentHint:  false,
			OpenWorldHint:   hint(true),
		},
	}
}

// publishArtifactTool declares the tool that records rendered artifacts.
func publishArtifactTool() *mcp.Tool {
	return &mcp.Tool{
		Name:  "publish_artifact",
		Title: "Publish an artifact",
		Description: "Record a markdown document in the synced tree, open it in the browser, " +
			"and return its canonical web URL and file path. Accepts either inline 'content' or a file 'path'. " +
			"Artifacts are markdown only: HTML is refused.\n\n" +
			"Use it when the answer is structural rather than linear: a comparison across many items, " +
			"a timeline, a graph, or an extensive overview. Answer in the conversation first and " +
			"record an artifact as an attachment to that answer, never in place of it. A durable finding belongs " +
			"in the wiki as memory; an artifact is a rendered presentation, and it is permanent.\n\n" +
			"The page must carry everything it needs inline. The document's title becomes the artifact's name, " +
			"so recording the same title replaces it rather than piling up copies.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: hint(true),
			IdempotentHint:  true,
			OpenWorldHint:   hint(false),
		},
	}
}

// flowToTool converts a recorded flow into an MCP tool that runs it. The
// wrapper is thin: the name is "run_flow_<name>" (snake_case) and the
// annotations match the generic run_flow tool because a flow's steps are
// arbitrary shell commands. New only calls it for a flow a human pinned, so
// the name and description it copies are text that human approved.
func flowToTool(f *flow.Flow) *mcp.Tool {
	return &mcp.Tool{
		Name:        "run_flow_" + strings.ReplaceAll(f.Name, "-", "_"),
		Title:       "Run flow: " + f.Name,
		Description: f.Description,
		Annotations: runFlowTool().Annotations,
	}
}
