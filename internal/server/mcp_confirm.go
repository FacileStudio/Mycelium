package server

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"path/filepath"
	"strconv"
	"time"

	"github.com/FacileStudio/Mycelium/internal/memory"
	apierrors "github.com/FacileStudio/tronc/errors"
	mcpkit "github.com/FacileStudio/tronc/mcp"
)

// kindDeletePage is the one sensitive action: deleting a wiki page.
const kindDeletePage = "delete_page"

// webActor marks, in the audit, a confirmation made on the confirm page.
const webActor = "mycelium-web"

// errPageGone is the refusal for a pending delete whose page or tree is no
// longer the person's.
var errPageGone = apierrors.NotFound("this page is gone, or no longer yours")

// errBrowserOnly refuses a confirm page call that is not a browser session:
// a CLI or API token is what agents hold, and an agent never confirms its
// own request.
var errBrowserOnly = apierrors.Forbidden("only your own browser session can confirm actions: open the link in the web app")

// errNoSubject refuses an account that has no Registre subject on record.
var errNoSubject = apierrors.Forbidden("sign out and in again so this account is linked to your Facile identity")

// pageTarget is what a pending delete names: the tree, the page and the
// log line the delete will write.
type pageTarget struct {
	SpaceID string `json:"space_id"`
	Path    string `json:"path"`
	Log     string `json:"log"`
}

// mcpStateKey derives the confirmation key from the OIDC client secret,
// namespaced so it collides with no other use of the secret.
func mcpStateKey(secret string) []byte {
	sum := sha256.Sum256([]byte("mycelium-mcp-state:" + secret))
	return sum[:]
}

// mcpKinds are the sensitive actions. The digest is the page's bytes, so an
// edit after the question closes the action instead of deleting content
// nobody saw.
func (t *mcpTools) kinds() map[string]mcpkit.Kind {
	return map[string]mcpkit.Kind{
		kindDeletePage: {Tool: toolDeleteMemoryPage, Fingerprint: t.pageDigest, Run: t.removePage},
	}
}

// targetOf reads the tree, page and log a pending delete names.
func targetOf(a mcpkit.Action) (pageTarget, error) {
	var target pageTarget
	if err := json.Unmarshal([]byte(a.Target), &target); err != nil || target.Path == "" {
		return pageTarget{}, errPageGone
	}
	return target, nil
}

// encodeTarget is the Target string of a pending delete.
func encodeTarget(target pageTarget) string {
	raw, _ := json.Marshal(target)
	return string(raw)
}

// actionTree resolves the tree of a pending action for the person it
// belongs to, as their browser session would see it now.
func (s *Server) actionTree(a mcpkit.Action) (pageTarget, mcpTree, error) {
	target, err := targetOf(a)
	if err != nil {
		return pageTarget{}, mcpTree{}, err
	}
	id, ok := s.accountOf(a.UserID)
	if !ok {
		return pageTarget{}, mcpTree{}, errPageGone
	}
	tree, err := s.mcpTreeFor(id, target.SpaceID)
	if err != nil {
		return pageTarget{}, mcpTree{}, errPageGone
	}
	return target, tree, nil
}

// pageDigest fingerprints the page a pending delete names.
func (t *mcpTools) pageDigest(_ context.Context, a mcpkit.Action) (string, error) {
	target, tree, err := t.s.actionTree(a)
	if err != nil {
		return "", err
	}
	key, data, err := memory.CheckDelete(tree.Root, target.Path, target.Log)
	if err != nil {
		return "", errPageGone
	}
	return mcpkit.Digest(tree.SpaceID, key, memory.Checksum(data)), nil
}

// removePage deletes the page a confirmed action names.
func (t *mcpTools) removePage(_ context.Context, a mcpkit.Action) error {
	target, tree, err := t.s.actionTree(a)
	if err != nil {
		return err
	}
	key, err := memory.DeletePage(tree.Root, target.Path, target.Log, time.Now())
	if err != nil {
		return t.s.pageFailure(err, target.Path, tree)
	}
	for _, file := range []string{key, "index.md", "log.md"} {
		t.s.enqueueEmbed(tree.Root, filepath.Join(tree.Root, "memory", filepath.FromSlash(file)))
	}
	return nil
}

// accountOf finds the account whose MCP user id is userID, and the
// Identity its browser session carries.
func (s *Server) accountOf(userID string) (Identity, bool) {
	s.mu.RLock()
	users := s.loadUsers()
	s.mu.RUnlock()
	for email, user := range users {
		if user.Subject == "" || strconv.FormatInt(subjectNumber(user.Issuer, user.Subject), 10) != userID {
			continue
		}
		scope := scopeUser
		if user.Admin {
			scope = scopeAdmin
		}
		return Identity{Email: email, Scope: scope}, true
	}
	return Identity{}, false
}
