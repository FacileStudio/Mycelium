package server

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/FacileStudio/Mycelium/internal/artifacts"
	apierrors "github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"
	mcpkit "github.com/FacileStudio/tronc/mcp"
)

const (
	defaultArtifactSearchLimit = 20
	maxArtifactSearchLimit     = 50
)

// ArtifactSearchHit is one artifact GET /api/artifacts/search found. Match is
// "title" when the title holds every word, else "content".
type ArtifactSearchHit struct {
	ArtifactSummary
	Match   string `json:"match"`
	Snippet string `json:"snippet"`
}

// ArtifactSearchResponse is one page of search hits.
type ArtifactSearchResponse struct {
	Results    []ArtifactSearchHit `json:"results"`
	Total      int                 `json:"total"`
	NextCursor string              `json:"next_cursor,omitempty"`
}

// artifactsSearch answers GET /api/artifacts/search: q over titles and
// bodies, in one tree (space_id), every readable tree (all=true), or the
// caller's default tree.
func (s *Server) artifactsSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	q := strings.TrimSpace(query.Get("q"))
	if q == "" || utf8.RuneCountInString(q) > artifacts.MaxQueryRunes {
		httpjson.WriteError(w, apierrors.Invalid("q must hold 1 to 200 characters"))
		return
	}
	var at mcpArtifactCursor
	if err := mcpkit.DecodeCursor(query.Get("cursor"), &at); err != nil || at.Offset < 0 {
		httpjson.WriteError(w, apierrors.Invalid("invalid cursor"))
		return
	}
	trees, ok := s.searchTrees(w, r)
	if !ok {
		return
	}
	hits, err := artifacts.Search(trees, q)
	if err != nil {
		s.Log.Error("artifacts search failed", slog.Any("error", err))
		httpjson.WriteError(w, apierrors.Internal("internal error", err))
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, searchPage(hits, at.Offset, searchLimit(query.Get("limit"))))
}

// searchTrees resolves the trees a search covers, through the same guards
// as the list routes.
func (s *Server) searchTrees(w http.ResponseWriter, r *http.Request) ([]artifacts.Tree, bool) {
	if r.URL.Query().Get("all") == "true" {
		return s.readableArtifactTrees(identityFrom(r)), true
	}
	root, ok := s.artifactRoot(w, r)
	if !ok {
		return nil, false
	}
	tree := artifacts.Tree{Root: root, SpaceID: r.URL.Query().Get("space_id")}
	if tree.SpaceID != "" {
		if space, _, found := s.spaceAccess(identityFrom(r), tree.SpaceID); found {
			tree.SpaceName = space.Name
		}
	}
	return []artifacts.Tree{tree}, true
}

func searchLimit(raw string) int {
	limit, err := strconv.Atoi(raw)
	if err != nil || limit <= 0 {
		return defaultArtifactSearchLimit
	}
	return min(limit, maxArtifactSearchLimit)
}

func searchPage(hits []artifacts.Hit, offset, limit int) ArtifactSearchResponse {
	out := ArtifactSearchResponse{Results: []ArtifactSearchHit{}, Total: len(hits)}
	rest := hits[min(offset, len(hits)):]
	if len(rest) > limit {
		rest = rest[:limit]
		out.NextCursor = mcpkit.EncodeCursor(mcpArtifactCursor{Offset: offset + limit})
	}
	for _, hit := range rest {
		summary := summaryOf(hit.Artifact)
		summary.SpaceID, summary.SpaceName = hit.SpaceID, hit.SpaceName
		out.Results = append(out.Results, ArtifactSearchHit{ArtifactSummary: summary, Match: hit.Match, Snippet: hit.Snippet})
	}
	return out
}
