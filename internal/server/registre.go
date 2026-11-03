package server

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/FacileStudio/porte/directory"
	apierrors "github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"
)

const (
	importRetryFirst = 30 * time.Second
	importRetryMax   = 5 * time.Minute
)

// registreRuntime holds the one directory syncer a server owns. It is built
// lazily so a Server literal in a test, which never ran New, still mounts the
// nudge route. backoff overrides the first import retry delay in tests.
type registreRuntime struct {
	once    sync.Once
	syncer  *directory.Syncer
	backoff time.Duration
}

// directorySyncer is the syncer that keeps spaces equal to Registre's. Without
// OIDC its client has no issuer and it is never run, so a nudge only fills a
// channel nobody reads.
func (s *Server) directorySyncer() *directory.Syncer {
	s.registre.once.Do(func() {
		syncer := &directory.Syncer{Apply: s.applySnapshot, Logger: s.Log}
		if s.OIDC != nil {
			syncer.Client = directory.Client{
				Issuer: s.OIDC.Issuer, ClientID: s.OIDC.ClientID, ClientSecret: s.OIDC.ClientSecret,
			}
		}
		s.registre.syncer = syncer
	})
	return s.registre.syncer
}

// RunRegistreSync hands Registre the spaces this instance shared before
// Registre owned them, retrying with backoff until it accepts, and only then
// starts pulling snapshots, because a snapshot applied first would drop every
// member Registre has not heard of yet. It blocks until ctx ends and returns
// at once when OIDC is not configured.
func (s *Server) RunRegistreSync(ctx context.Context) {
	if s.OIDC == nil {
		return
	}
	syncer := s.directorySyncer()
	if !s.importSpaces(ctx, syncer.Client) {
		return
	}
	syncer.Run(ctx)
}

// importSpaces calls Import until it succeeds or ctx ends, doubling the wait
// up to importRetryMax so a Registre that is not deployed yet costs one log
// line every few minutes and nothing else.
func (s *Server) importSpaces(ctx context.Context, client directory.Client) bool {
	delay := s.registre.backoff
	if delay <= 0 {
		delay = importRetryFirst
	}
	for {
		_, err := client.Import(ctx, s.sharedSpaces())
		if err == nil {
			return true
		}
		s.Log.Warn("registre: space import failed, retrying", slog.Any("error", err), slog.Duration("in", delay))
		select {
		case <-ctx.Done():
			return false
		case <-time.After(delay):
		}
		delay = min(delay*2, importRetryMax)
	}
}

// sharedSpaces is what this instance offers Registre on import: every
// unlinked space with more than one member. A space with a single member is
// someone's private memory tree, and Mycelium has no other sharing signal, so
// it is never sent and never merged by name with anyone else's space.
func (s *Server) sharedSpaces() []directory.ImportSpace {
	s.mu.RLock()
	spaces := s.loadSpaces()
	s.mu.RUnlock()
	out := []directory.ImportSpace{}
	for _, space := range spaces {
		if space.RegistreID != "" || len(space.Members) < 2 {
			continue
		}
		members := make([]directory.ImportMember, 0, len(space.Members))
		for email, role := range space.Members {
			members = append(members, directory.ImportMember{Email: email, Role: role})
		}
		sort.Slice(members, func(i, j int) bool { return members[i].Email < members[j].Email })
		out = append(out, directory.ImportSpace{Name: space.Name, Members: members})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// registreSpacesURL is where a person manages spaces in Registre, or one
// space when id is set. It is empty when OIDC is not configured.
func (s *Server) registreSpacesURL(id string) string {
	if s.OIDC == nil {
		return ""
	}
	base := strings.TrimRight(s.OIDC.Issuer, "/") + "/spaces"
	if id == "" {
		return base
	}
	return base + "/" + url.PathEscape(id)
}

// refuseManaged answers 409 when Registre owns the change. A linked space's
// name and existence are Registre's; with OIDC configured, so is every
// space's membership, which keeps a personal space from quietly becoming a
// shared one Registre never hears of.
func (s *Server) refuseManaged(w http.ResponseWriter, space *Space, membership bool) bool {
	if space.RegistreID == "" && (!membership || s.OIDC == nil) {
		return false
	}
	message := "this space is managed in Registre"
	if link := s.registreSpacesURL(space.RegistreID); link != "" {
		message += ": " + link
	}
	httpjson.WriteError(w, apierrors.Conflict(message))
	return true
}
