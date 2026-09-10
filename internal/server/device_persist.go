package server

import (
	"crypto/rand"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

// generateUserCode makes the code a person reads and approves with. It is the
// value persisted in devices.json, so it lives beside the persistence.
func generateUserCode() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	out := make([]byte, 0, 9)
	for i, c := range b {
		if i == 4 {
			out = append(out, '-')
		}
		out = append(out, userCodeAlphabet[int(c)%len(userCodeAlphabet)])
	}
	return string(out), nil
}

// Persistence. create, approve, deny and the terminal poll each leave the
// in-memory map changed, so each rewrites the pending set. Only pending
// requests are ever written: an approved one holds the minted token, which must
// not reach disk, and it is consumed within seconds anyway, so a restart in
// that window loses a login the terminal simply retries.

// deviceKey is the store key for a device code. It is the code's hash, never
// the code: the at-rest rule that governs tokens.json applies to these too, and
// a leaked devices.json must not hand out a working bearer.
func deviceKey(code string) string {
	return hashToken(code)
}

func (d *deviceStore) load() {
	if d.path == "" {
		return
	}
	data, err := os.ReadFile(d.path)
	if err != nil {
		return
	}
	var raw map[string]deviceRequest
	if err := json.Unmarshal(data, &raw); err != nil {
		if d.log != nil {
			d.log.Error("devices: failed to parse", slog.String("path", d.path), slog.Any("error", err))
		}
		return
	}
	now := time.Now()
	for key, req := range raw {
		if req.Status != devicePending || now.After(req.Expires) {
			continue
		}
		p := req
		d.byDevice[key] = &p
		d.byUser[normalizeUserCode(p.UserCode)] = &p
	}
}

func (d *deviceStore) persist() {
	if d.path == "" {
		return
	}
	out := make(map[string]deviceRequest)
	for key, req := range d.byDevice {
		if req.Status == devicePending {
			out[key] = *req
		}
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		if d.log != nil {
			d.log.Error("devices: marshal failed", slog.Any("error", err))
		}
		return
	}
	if err := atomicWriteFile(d.path, data); err != nil {
		if d.log != nil {
			d.log.Error("devices: write failed", slog.String("path", d.path), slog.Any("error", err))
		}
	}
}

// atomicWriteFile stages the data in a temp file beside the target and renames
// it into place, so a crash mid-write leaves the old file intact instead of a
// torn one that would fail to parse on the next load and come up empty. A torn
// devices.json or login-codes.json would silently drop the very pending logins
// this persistence exists to protect.
func atomicWriteFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".*")
	if err != nil {
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}
	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	if err := os.Rename(tmp.Name(), path); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	return nil
}
