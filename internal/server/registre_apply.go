package server

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/FacileStudio/porte/directory"
)

// snapshotApply is one pass of a Registre snapshot over the two stores. It
// records whether anything changed so a repeated snapshot writes nothing.
type snapshotApply struct {
	spaces      map[string]*Space
	users       map[string]User
	linked      map[string]*Space
	byName      map[string]*Space
	created     []string
	seen        map[string]bool
	spacesDirty bool
	usersDirty  bool
	issuer      string
	now         string
}

// applySnapshot makes the spaces and users stores match Registre's snapshot.
// Spaces link by Registre id, then by name to an unlinked shared space, and
// are created otherwise. Members of a linked space become exactly the
// snapshot's: every member email here came from a Registre sign-in, so there
// is no local-only member to spare. People not seen before are provisioned as
// non-admin accounts so their first sign-in lands on them. A linked space
// Registre no longer lists is logged and left alone, and personal spaces are
// never touched. Holding the write lock for the whole pass makes it atomic to
// every handler.
func (s *Server) applySnapshot(_ context.Context, snapshot directory.Snapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	spaces, err := readStore[*Space](s.spacesPath())
	if err != nil {
		return fmt.Errorf("registre: read spaces: %w", err)
	}
	users, err := readStore[User](s.usersPath())
	if err != nil {
		return fmt.Errorf("registre: read users: %w", err)
	}
	pass := newSnapshotApply(spaces, users, s.issuer())
	for _, remote := range snapshot.Spaces {
		pass.apply(remote)
	}
	for registreID, space := range pass.linked {
		if !pass.seen[registreID] {
			s.Log.Warn("registre: linked space missing from the snapshot, left alone",
				slog.String("space", space.ID), slog.String("registre_space_id", registreID))
		}
	}
	return s.saveApplied(pass)
}

// newSnapshotApply indexes the local spaces a snapshot can land on.
func newSnapshotApply(spaces map[string]*Space, users map[string]User, issuer string) *snapshotApply {
	pass := &snapshotApply{
		spaces: spaces, users: users, issuer: issuer, linked: map[string]*Space{}, byName: map[string]*Space{},
		seen: map[string]bool{}, now: time.Now().UTC().Format(time.RFC3339),
	}
	ids := make([]string, 0, len(spaces))
	for id := range spaces {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		space := spaces[id]
		if space.RegistreID != "" {
			pass.linked[space.RegistreID] = space
			continue
		}
		key := strings.ToLower(space.Name)
		if _, taken := pass.byName[key]; !taken && len(space.Members) > 1 {
			pass.byName[key] = space
		}
	}
	return pass
}

// apply lands one Registre space on its local space and sets its name and
// members from the snapshot.
func (p *snapshotApply) apply(remote directory.Space) {
	p.seen[remote.ID] = true
	space := p.local(remote)
	members := map[string]string{}
	for _, member := range remote.Members {
		if email := p.resolve(member); email != "" {
			members[email] = localRole(member.Role)
		}
	}
	if space.Name == remote.Name && maps.Equal(space.Members, members) {
		return
	}
	space.Name, space.Members, space.UpdatedAt = remote.Name, members, p.now
	p.spacesDirty = true
}

// local finds or creates the local space a Registre space maps to.
func (p *snapshotApply) local(remote directory.Space) *Space {
	if space, ok := p.linked[remote.ID]; ok {
		return space
	}
	key := strings.ToLower(remote.Name)
	space, ok := p.byName[key]
	if ok {
		delete(p.byName, key)
	} else {
		space = &Space{ID: newSpaceID(), Members: map[string]string{}, CreatedAt: p.now, UpdatedAt: p.now}
		p.spaces[space.ID] = space
		p.created = append(p.created, space.ID)
	}
	space.RegistreID = remote.ID
	p.linked[remote.ID] = space
	p.spacesDirty = true
	return space
}

// resolve maps a Registre member onto a stored account by email, case
// insensitively, and provisions one when there is none. A provisioned account
// is never an admin, as with ensureDelegatedUser. The member's subject is
// stored on an account that has none yet, so a deletion notice finds it.
func (p *snapshotApply) resolve(member directory.Member) string {
	email := strings.TrimSpace(member.Email)
	if email == "" {
		return ""
	}
	key, ok := findByEmail(p.users, email)
	if !ok {
		key = email
		p.users[key] = User{Email: email, Name: member.Name, AvatarURL: pictureURL(member.Picture), CreatedAt: p.now}
		p.usersDirty = true
	}
	if user := p.users[key]; user.Subject == "" && member.Subject != "" && p.issuer != "" {
		user.Issuer, user.Subject = p.issuer, member.Subject
		p.users[key] = user
		p.usersDirty = true
	}
	return key
}

// localRole maps a Registre role onto Mycelium's, which share a vocabulary.
// Anything unknown is an ordinary member rather than a grant.
func localRole(role string) string {
	if validRole(role) {
		return role
	}
	return roleMember
}

// saveApplied writes what the pass changed, users first so no space ever
// names a member the users store has not recorded, then gives every created
// space its content directories.
func (s *Server) saveApplied(pass *snapshotApply) error {
	if pass.usersDirty {
		if err := s.saveUsers(pass.users); err != nil {
			return fmt.Errorf("registre: save users: %w", err)
		}
	}
	if pass.spacesDirty {
		if err := s.saveSpaces(pass.spaces); err != nil {
			return fmt.Errorf("registre: save spaces: %w", err)
		}
	}
	for _, id := range pass.created {
		if err := s.makeSpaceDirs(id); err != nil {
			return fmt.Errorf("registre: create space %s: %w", id, err)
		}
	}
	return nil
}

// makeSpaceDirs gives a space the content directories every client expects.
func (s *Server) makeSpaceDirs(id string) error {
	for _, dir := range []string{"memory", "rules", "skills", "sessions"} {
		if err := os.MkdirAll(filepath.Join(s.spaceDir(id), dir), 0o755); err != nil {
			return err
		}
	}
	return nil
}
