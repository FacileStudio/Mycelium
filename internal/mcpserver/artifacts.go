package mcpserver

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/FacileStudio/Mycelium/internal/artifacts"
	"github.com/FacileStudio/Mycelium/internal/browser"
	"github.com/FacileStudio/Mycelium/internal/config"
	hsync "github.com/FacileStudio/Mycelium/internal/sync"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type publishArtifactInput struct {
	Path    string `json:"path,omitempty" jsonschema:"path to the markdown (.md) file to record; artifacts are markdown only (Facile Markdown: CommonMark with GFM tables and task lists plus callout, stat, chart and columns fences, mermaid/d2/dot diagram fences, [badge:tone text], [status:state] and ==highlight==), at most 1 MiB; HTML is refused and an invalid document is refused with encre's reason and its line and column"`
	Content string `json:"content,omitempty" jsonschema:"inline markdown content to record as an artifact (alternative to path), at most 1 MiB; HTML is refused and an invalid document is refused with encre's reason and its line and column"`
	Title   string `json:"title,omitempty" jsonschema:"overrides the document's own title as the artifact's name"`
}

type publishArtifactOutput struct {
	ID       string   `json:"id" jsonschema:"the artifact's name, used by mycelium artifact open and rm"`
	Path     string   `json:"path" jsonschema:"where the page was written on this machine"`
	URL      string   `json:"url,omitempty" jsonschema:"the canonical web URL of the artifact on the Mycelium server"`
	Opened   bool     `json:"opened" jsonschema:"true when a browser was opened here, false when this machine has no browser to open"`
	Unusable []string `json:"unresolved_refs,omitempty" jsonschema:"relative markdown links or images that cannot load from disk"`
}

func recordArtifact(in publishArtifactInput) (artifacts.Artifact, []byte, error) {
	path := strings.TrimSpace(in.Path)
	content := strings.TrimSpace(in.Content)
	if path == "" && content == "" {
		return artifacts.Artifact{}, nil, errors.New("publish_artifact needs either path or inline content")
	}
	req := artifacts.Request{
		Title:   in.Title,
		Machine: config.MachineName(),
	}
	if content != "" {
		raw := []byte(in.Content)
		art, err := artifacts.AddContent(config.DataDir(), raw, "inline.md", req, time.Now())
		return art, raw, err
	}
	req.Source = path
	raw, _ := os.ReadFile(path)
	art, err := artifacts.Add(config.DataDir(), req, time.Now())
	return art, raw, err
}

func syncArtifact() string {
	cfg, err := config.LoadMyceliumConfig()
	if err != nil || cfg == nil || cfg.ServerURL() == "" {
		return ""
	}
	client := hsync.NewClient(cfg.ServerURL(), cfg.AuthToken())
	client.Space = cfg.SpaceID()
	//nolint:errcheck // best-effort sync: the artifact is already recorded, the URL must survive a dead server
	client.Sync(config.DataDir())
	return cfg.ServerURL()
}

func publishArtifact(_ context.Context, _ *mcp.CallToolRequest, in publishArtifactInput) (*mcp.CallToolResult, publishArtifactOutput, error) {
	art, raw, err := recordArtifact(in)
	if err != nil {
		return nil, publishArtifactOutput{}, err
	}
	serverURL := syncArtifact()
	artURL := artifacts.URL(serverURL, art.ID)
	out := publishArtifactOutput{ID: art.ID, Path: art.Path, URL: artURL}
	if len(raw) > 0 {
		out.Unusable = artifacts.ExternalRefs(raw)
	}
	if browser.Available() {
		target := art.Path
		if artURL != "" {
			target = artURL
		}
		out.Opened = browser.Open(target) == nil
	}
	return nil, out, nil
}
