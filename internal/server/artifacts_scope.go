package server

import (
	"log/slog"
	"net/http"
	"sort"

	"github.com/FacileStudio/Mycelium/internal/artifacts"
	apierrors "github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"
)

// artifactCommonAllowed widens commonAllowed for the artifact routes only: a
// delegated suite token of an admin reaches the admin's own artifacts in the
// common tree, which is that admin's personal tree, so Bulle shows the same
// list as the CLI. Rules, memory, skills and sync stay closed to it.
func (s *Server) artifactCommonAllowed(id Identity) bool {
	if s.commonAllowed(id) {
		return true
	}
	if !id.Delegated || id.Email == "" {
		return false
	}
	s.mu.RLock()
	user, ok := s.loadUsers()[id.Email]
	s.mu.RUnlock()
	return ok && user.Admin
}

// artifactRoot is scopeRoot for the artifact routes, with the common tree
// gated by artifactCommonAllowed.
func (s *Server) artifactRoot(w http.ResponseWriter, r *http.Request) (string, bool) {
	if r.URL.Query().Get("space_id") == "" && s.artifactCommonAllowed(identityFrom(r)) {
		return s.DataDir, true
	}
	return s.scopeRoot(w, r)
}

// memberSpaces answers the spaces the caller is a member of, by name.
func (s *Server) memberSpaces(id Identity) []*Space {
	s.mu.RLock()
	spaces := s.loadSpaces()
	s.mu.RUnlock()
	out := []*Space{}
	for _, space := range spaces {
		if _, ok := space.Members[id.Email]; ok && id.Email != "" {
			out = append(out, space)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// readableArtifactTrees is every artifact tree the caller can read: the
// common tree when artifactCommonAllowed lets them, then each space they are a
// member of, by name.
func (s *Server) readableArtifactTrees(id Identity) []artifacts.Tree {
	trees := []artifacts.Tree{}
	if s.artifactCommonAllowed(id) {
		trees = append(trees, artifacts.Tree{Root: s.DataDir})
	}
	for _, space := range s.memberSpaces(id) {
		trees = append(trees, artifacts.Tree{Root: s.spaceDir(space.ID), SpaceID: space.ID, SpaceName: space.Name})
	}
	return trees
}

// artifactsListAll lists every artifact the caller can read, newest first.
// Each summary carries the space it lives in.
func (s *Server) artifactsListAll(w http.ResponseWriter, r *http.Request) {
	out := []ArtifactSummary{}
	for _, tree := range s.readableArtifactTrees(identityFrom(r)) {
		all, err := artifacts.List(tree.Root)
		if err != nil {
			s.Log.Error("artifacts list failed", slog.Any("error", err))
			httpjson.WriteError(w, apierrors.Internal("internal error", err))
			return
		}
		for _, art := range all {
			summary := summaryOf(art)
			summary.SpaceID, summary.SpaceName = tree.SpaceID, tree.SpaceName
			out = append(out, summary)
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Created.After(out[j].Created) })
	httpjson.WriteJSON(w, http.StatusOK, out)
}
