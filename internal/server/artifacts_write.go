package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/FacileStudio/Mycelium/internal/artifacts"
	apierrors "github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"
)

const (
	maxArtifactContent = 1 << 20
	maxArtifactTitle   = 200
	maxArtifactBody    = 8 << 20
)

var htmlMarkers = []string{"<!doctype", "<html", "<head>", "<head ", "<body", "<script", "<iframe"}

// ArtifactWriteRequest is the body of POST and PUT /api/artifacts. Machine is
// read on create only.
type ArtifactWriteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Machine string `json:"machine,omitempty"`
}

// looksLikeHTML reports whether content is an HTML document or carries
// active markup, which the write routes refuse: they accept markdown only.
func looksLikeHTML(content string) bool {
	lower := strings.ToLower(content)
	for _, marker := range htmlMarkers {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

// decodeArtifactWrite reads and validates a write body, answering the error
// itself when the body is unusable.
func decodeArtifactWrite(w http.ResponseWriter, r *http.Request) (ArtifactWriteRequest, bool) {
	var req ArtifactWriteRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxArtifactBody)).Decode(&req); err != nil {
		var tooBig *http.MaxBytesError
		if errors.As(err, &tooBig) {
			httpjson.WriteError(w, apierrors.TooLarge("content exceeds 1 MB"))
			return req, false
		}
		httpjson.WriteError(w, apierrors.Invalid("invalid body"))
		return req, false
	}
	req.Title = strings.TrimSpace(req.Title)
	if err := validateArtifactWrite(req); err != nil {
		httpjson.WriteError(w, err)
		return req, false
	}
	return req, true
}

func validateArtifactWrite(req ArtifactWriteRequest) *apierrors.Error {
	switch {
	case len(req.Content) > maxArtifactContent:
		return apierrors.TooLarge("content exceeds 1 MB")
	case utf8.RuneCountInString(req.Title) > maxArtifactTitle:
		return apierrors.Invalid("title exceeds 200 characters")
	case strings.TrimSpace(req.Content) == "":
		return apierrors.Invalid("content is required")
	case looksLikeHTML(req.Content):
		return apierrors.Invalid("content must be markdown, not HTML")
	}
	return nil
}

func (s *Server) artifactCreate(w http.ResponseWriter, r *http.Request) {
	root, ok := s.artifactRoot(w, r)
	if !ok {
		return
	}
	req, ok := decodeArtifactWrite(w, r)
	if !ok {
		return
	}
	art, err := artifacts.AddContent(root, []byte(req.Content), "api.md", artifacts.Request{
		Title:   req.Title,
		Machine: strings.TrimSpace(req.Machine),
	}, time.Now())
	if err != nil {
		httpjson.WriteError(w, apierrors.Internal("internal error", err))
		return
	}
	s.answerArtifact(w, art, http.StatusCreated)
}

func (s *Server) artifactUpdate(w http.ResponseWriter, r *http.Request) {
	id, ok := safeName(pathParam(r, "id"))
	if !ok {
		httpjson.WriteError(w, apierrors.Invalid("invalid artifact id"))
		return
	}
	root, ok := s.artifactRoot(w, r)
	if !ok {
		return
	}
	req, ok := decodeArtifactWrite(w, r)
	if !ok {
		return
	}
	existing, err := artifacts.Find(root, id)
	if err != nil || !strings.EqualFold(existing.ID, id) {
		httpjson.WriteError(w, apierrors.NotFound("not found"))
		return
	}
	if artifactFormat(existing.Path) != "markdown" {
		httpjson.WriteError(w, apierrors.Invalid("only markdown artifacts can be updated"))
		return
	}
	art, err := artifacts.Update(root, id, req.Title, []byte(req.Content), time.Now())
	if err != nil {
		httpjson.WriteError(w, apierrors.NotFound("not found"))
		return
	}
	s.answerArtifact(w, art, http.StatusOK)
}

func (s *Server) answerArtifact(w http.ResponseWriter, art artifacts.Artifact, status int) {
	detail, err := detailOf(art)
	if err != nil {
		httpjson.WriteError(w, apierrors.Internal("internal error", err))
		return
	}
	httpjson.WriteJSON(w, status, detail)
}
