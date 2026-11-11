package server

import (
	"errors"
	"strings"

	"github.com/FacileStudio/porte/mcpauth"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// personalTree is the name the tools give the common tree, which is the
// instance owner's own wiki and artifacts, not a shared bucket.
const personalTree = "Personal"

var (
	errMCPNoCaller = errors.New("this request carries no valid access token for a Mycelium account. Reconnect Mycelium")
	errMCPNoTree   = errors.New("you have no personal tree on this server, only spaces. Pass a space_id: " +
		"call mycelium_read_memory_index to list the spaces you belong to")
	errMCPSpace = errors.New("no space with this space_id is yours. " +
		"Call mycelium_read_memory_index to list the spaces you belong to")
)

// mcpTree is one tree a person can read and write: the common tree (empty
// SpaceID) or a space they are a member of.
type mcpTree struct {
	SpaceID string
	Name    string
	Role    string
	Root    string
}

// mcpCaller resolves the person a tool call acts for to the Identity a
// browser session of theirs carries: admin when their account is admin,
// user otherwise, never delegated. Every tree guard downstream then answers
// exactly as it does for that person in the web app.
func (s *Server) mcpCaller(req *sdk.CallToolRequest) (Identity, error) {
	if req == nil || req.Extra == nil {
		return Identity{}, errMCPNoCaller
	}
	info := req.Extra.TokenInfo
	if _, ok := mcpauth.Identity(info); !ok {
		return Identity{}, errMCPNoCaller
	}
	subject, _ := info.Extra[mcpauth.ExtraSubject].(string)
	actor, _ := info.Extra[mcpauth.ExtraActor].(string)
	s.mu.RLock()
	users := s.loadUsers()
	s.mu.RUnlock()
	key, ok := findBySubject(users, s.OIDC.Issuer, subject)
	if subject == "" || !ok {
		return Identity{}, errMCPNoCaller
	}
	scope := scopeUser
	if users[key].Admin {
		scope = scopeAdmin
	}
	return Identity{Email: key, Scope: scope, Actor: actor}, nil
}

// mcpTrees lists every tree the person can use, the common tree first when
// they may reach it, then their spaces by name.
func (s *Server) mcpTrees(id Identity) []mcpTree {
	trees := []mcpTree{}
	if s.commonAllowed(id) {
		trees = append(trees, mcpTree{Name: personalTree, Role: roleOwner, Root: s.DataDir})
	}
	for _, space := range s.memberSpaces(id) {
		trees = append(trees, mcpTree{
			SpaceID: space.ID, Name: space.Name, Role: space.Members[id.Email], Root: s.spaceDir(space.ID),
		})
	}
	return trees
}

// mcpTreeFor resolves the tree a tool call names: the common tree for an
// empty space_id, else a space after the same membership guard every API
// route uses.
func (s *Server) mcpTreeFor(id Identity, spaceID string) (mcpTree, error) {
	spaceID = strings.TrimSpace(spaceID)
	if spaceID == "" {
		if !s.commonAllowed(id) {
			return mcpTree{}, errMCPNoTree
		}
		return mcpTree{Name: personalTree, Role: roleOwner, Root: s.DataDir}, nil
	}
	if strings.ContainsAny(spaceID, "/\\.") {
		return mcpTree{}, errMCPSpace
	}
	space, role, ok := s.spaceAccess(id, spaceID)
	if !ok {
		return mcpTree{}, errMCPSpace
	}
	return mcpTree{SpaceID: space.ID, Name: space.Name, Role: role, Root: s.spaceDir(space.ID)}, nil
}

// mcpScope resolves the trees a read covers: the one named, or every tree
// the person can use when space_id is empty.
func (s *Server) mcpScope(req *sdk.CallToolRequest, spaceID string) (Identity, []mcpTree, error) {
	id, err := s.mcpCaller(req)
	if err != nil {
		return id, nil, err
	}
	if strings.TrimSpace(spaceID) == "" {
		return id, s.mcpTrees(id), nil
	}
	tree, err := s.mcpTreeFor(id, spaceID)
	if err != nil {
		return id, nil, err
	}
	return id, []mcpTree{tree}, nil
}

// mcpOneTree resolves the one tree a page read or a write uses.
func (s *Server) mcpOneTree(req *sdk.CallToolRequest, spaceID string) (mcpTree, error) {
	id, err := s.mcpCaller(req)
	if err != nil {
		return mcpTree{}, err
	}
	return s.mcpTreeFor(id, spaceID)
}
