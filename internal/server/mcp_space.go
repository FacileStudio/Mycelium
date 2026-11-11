package server

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	mcpkit "github.com/FacileStudio/tronc/mcp"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// maxSpaceName caps a space name made through MCP.
const maxSpaceName = 100

// createSpace is mycelium_create_space: the web app's space creation, made
// idempotent by name, so an agent bootstrapping its space twice gets one.
func (t *mcpTools) createSpace(_ context.Context, req *sdk.CallToolRequest, in mcpCreateSpaceInput) (*sdk.CallToolResult, mcpCreateSpaceOutput, error) {
	id, err := t.s.mcpCaller(req)
	if err != nil {
		return nil, mcpCreateSpaceOutput{}, err
	}
	name := strings.TrimSpace(in.Name)
	if name == "" || len(name) > maxSpaceName {
		return nil, mcpCreateSpaceOutput{}, fmt.Errorf("name must be 1 to %d characters", maxSpaceName)
	}
	space, created, err := t.s.ownSpace(id.Email, name, strings.TrimSpace(in.Description))
	if err != nil {
		t.s.Log.Error("mcp: creating a space failed", slog.Any("error", err))
		return nil, mcpCreateSpaceOutput{}, mcpkit.ErrUnavailable
	}
	return nil, mcpCreateSpaceOutput{SpaceID: space.ID, Name: space.Name, Role: roleOwner, Created: created}, nil
}

// ownSpace returns the space named name that email owns and Registre does
// not manage, creating it when there is none.
func (s *Server) ownSpace(email, name, description string) (*Space, bool, error) {
	s.mu.Lock()
	spaces := s.loadSpaces()
	for _, space := range spaces {
		if space.RegistreID == "" && space.Members[email] == roleOwner && strings.EqualFold(space.Name, name) {
			s.mu.Unlock()
			return space, false, nil
		}
	}
	now := time.Now().UTC().Format(time.RFC3339)
	space := &Space{
		ID: newSpaceID(), Name: name, Description: description,
		Members: map[string]string{email: roleOwner}, CreatedAt: now, UpdatedAt: now,
	}
	spaces[space.ID] = space
	err := s.saveSpaces(spaces)
	s.mu.Unlock()
	if err != nil {
		return nil, false, err
	}
	return space, true, s.makeSpaceDirs(space.ID)
}
