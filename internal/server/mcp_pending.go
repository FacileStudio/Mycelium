package server

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"

	mcpkit "github.com/FacileStudio/tronc/mcp"
)

// pendingFile holds the actions waiting for a person. It is a dotfile under
// the data directory, so sync never carries it.
const pendingFile = ".mcp-actions.json"

// pendingKeep is how long a resolved or expired action stays on file.
const pendingKeep = 30 * 24 * time.Hour

// pendingRow is one stored action with what Action leaves out of its JSON.
type pendingRow struct {
	mcpkit.Action
	Digest     string     `json:"digest"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	Outcome    string     `json:"outcome,omitempty"`
}

// filePending is mcpkit.Pending on a JSON file. Mycelium has no database;
// one mutex makes every Claim atomic, which is what makes a confirmation
// single use.
type filePending struct {
	mu   sync.Mutex
	path string
}

// newFilePending stores actions in dataDir.
func newFilePending(dataDir string) *filePending {
	return &filePending{path: filepath.Join(dataDir, pendingFile)}
}

// load reads every row; a missing file is an empty store.
func (p *filePending) load() ([]pendingRow, error) {
	data, err := os.ReadFile(p.path)
	if errors.Is(err, os.ErrNotExist) {
		return []pendingRow{}, nil
	}
	if err != nil {
		return nil, err
	}
	rows := []pendingRow{}
	return rows, json.Unmarshal(data, &rows)
}

// save writes the rows back, dropping those long settled.
func (p *filePending) save(rows []pendingRow) error {
	cutoff := time.Now().Add(-pendingKeep)
	rows = slices.DeleteFunc(rows, func(r pendingRow) bool {
		return r.ExpiresAt.Before(cutoff) || (r.ResolvedAt != nil && r.ResolvedAt.Before(cutoff))
	})
	data, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		return err
	}
	return writeFileAtomic(p.path, data, 0o600)
}

// update runs change on the rows under the lock and saves them when it
// returns no error.
func (p *filePending) update(change func(rows []pendingRow) ([]pendingRow, error)) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	rows, err := p.load()
	if err != nil {
		return err
	}
	if rows, err = change(rows); err != nil {
		return err
	}
	return p.save(rows)
}

// open reports whether a row waits for userID at now.
func (r pendingRow) open(userID string, now time.Time) bool {
	return r.UserID == userID && r.ResolvedAt == nil && now.Before(r.ExpiresAt)
}

// Create stores a new action.
func (p *filePending) Create(_ context.Context, a mcpkit.Action) error {
	return p.update(func(rows []pendingRow) ([]pendingRow, error) {
		return append(rows, pendingRow{Action: a, Digest: a.Digest}), nil
	})
}

// Open lists userID's open actions, newest first.
func (p *filePending) Open(_ context.Context, userID string) ([]mcpkit.Action, error) {
	p.mu.Lock()
	rows, err := p.load()
	p.mu.Unlock()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	actions := []mcpkit.Action{}
	for _, r := range rows {
		if r.open(userID, now) {
			actions = append(actions, r.withDigest())
		}
	}
	slices.SortStableFunc(actions, func(a, b mcpkit.Action) int { return b.CreatedAt.Compare(a.CreatedAt) })
	return actions, nil
}

// withDigest is the row's action with its digest restored.
func (r pendingRow) withDigest() mcpkit.Action {
	a := r.Action
	a.Digest = r.Digest
	return a
}

// Claim marks one open action running and returns it, or ErrGone.
func (p *filePending) Claim(_ context.Context, userID, id string) (mcpkit.Action, error) {
	var claimed mcpkit.Action
	err := p.update(func(rows []pendingRow) ([]pendingRow, error) {
		now := time.Now()
		for i := range rows {
			if rows[i].ID == id && rows[i].open(userID, now) {
				rows[i].ResolvedAt, rows[i].Outcome = &now, mcpkit.OutcomeRunning
				claimed = rows[i].withDigest()
				return rows, nil
			}
		}
		return nil, mcpkit.ErrGone
	})
	return claimed, err
}

// Finish records the outcome of a claimed action.
func (p *filePending) Finish(_ context.Context, id, outcome string) error {
	return p.set(id, func(r *pendingRow) { r.Outcome = outcome })
}

// Release reopens a claimed action whose run failed.
func (p *filePending) Release(_ context.Context, id string) error {
	return p.set(id, func(r *pendingRow) { r.ResolvedAt, r.Outcome = nil, "" })
}

// Resolve closes one of userID's open actions without running it.
func (p *filePending) Resolve(_ context.Context, userID, id, outcome string) error {
	return p.update(func(rows []pendingRow) ([]pendingRow, error) {
		now := time.Now()
		for i := range rows {
			if rows[i].ID == id && rows[i].open(userID, now) {
				rows[i].ResolvedAt, rows[i].Outcome = &now, outcome
				return rows, nil
			}
		}
		return nil, mcpkit.ErrGone
	})
}

// Erase deletes every action of the people named.
func (p *filePending) Erase(_ context.Context, userIDs ...string) error {
	return p.update(func(rows []pendingRow) ([]pendingRow, error) {
		return slices.DeleteFunc(rows, func(r pendingRow) bool { return slices.Contains(userIDs, r.UserID) }), nil
	})
}

// set changes the row with id, if there is one.
func (p *filePending) set(id string, change func(*pendingRow)) error {
	return p.update(func(rows []pendingRow) ([]pendingRow, error) {
		for i := range rows {
			if rows[i].ID == id {
				change(&rows[i])
			}
		}
		return rows, nil
	})
}
