package server

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/FacileStudio/Mycelium/internal/artifacts"
	apierrors "github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"
)

// ArtifactSummary represents the metadata of an artifact for API listings.
type ArtifactSummary struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Machine string    `json:"machine"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	// SpaceID and SpaceName name the tree an artifact lives in, and are set
	// only by the all=true listing; an empty SpaceID there is the common tree.
	SpaceID   string `json:"space_id,omitempty"`
	SpaceName string `json:"space_name,omitempty"`
}

// ReportSummary is a backward-compatible alias for ArtifactSummary.
type ReportSummary = ArtifactSummary

// ArtifactDetail represents the full content and metadata of an artifact.
type ArtifactDetail struct {
	ArtifactSummary
	Content string `json:"content"`
}

// ReportDetail is a backward-compatible alias for ArtifactDetail.
type ReportDetail = ArtifactDetail

func summaryOf(art artifacts.Artifact) ArtifactSummary {
	return ArtifactSummary{
		ID:      art.ID,
		Title:   art.Title,
		Machine: art.Machine,
		Created: art.Created,
		Updated: art.Updated,
	}
}

// detailOf reads the artifact's file back, frontmatter included, so every
// route answers with exactly what is on disk.
func detailOf(art artifacts.Artifact) (ArtifactDetail, error) {
	raw, err := os.ReadFile(art.Path)
	if err != nil {
		return ArtifactDetail{}, err
	}
	return ArtifactDetail{ArtifactSummary: summaryOf(art), Content: string(raw)}, nil
}

func (s *Server) artifactsList(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("all") == "true" {
		s.artifactsListAll(w, r)
		return
	}
	root, ok := s.artifactRoot(w, r)
	if !ok {
		return
	}
	all, err := artifacts.List(root)
	if err != nil {
		s.Log.Error("artifacts list failed", slog.Any("error", err))
		httpjson.WriteError(w, apierrors.Internal("internal error", err))
		return
	}
	out := make([]ArtifactSummary, 0, len(all))
	for _, art := range all {
		out = append(out, summaryOf(art))
	}
	httpjson.WriteJSON(w, http.StatusOK, out)
}

func (s *Server) artifactGet(w http.ResponseWriter, r *http.Request) {
	id, ok := safeName(pathParam(r, "id"))
	if !ok {
		httpjson.WriteError(w, apierrors.Invalid("invalid artifact id"))
		return
	}
	root, rootOK := s.artifactRoot(w, r)
	if !rootOK {
		return
	}
	art, err := artifacts.Find(root, id)
	if err != nil {
		httpjson.WriteError(w, apierrors.NotFound("not found"))
		return
	}
	detail, err := detailOf(art)
	if err != nil {
		httpjson.WriteError(w, apierrors.NotFound("not found"))
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, detail)
}

func (s *Server) artifactDelete(w http.ResponseWriter, r *http.Request) {
	id, ok := safeName(pathParam(r, "id"))
	if !ok {
		httpjson.WriteError(w, apierrors.Invalid("invalid artifact id"))
		return
	}
	root, rootOK := s.artifactRoot(w, r)
	if !rootOK {
		return
	}
	if err := artifacts.Remove(root, id); err != nil {
		httpjson.WriteError(w, apierrors.NotFound("not found"))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
