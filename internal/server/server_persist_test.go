package server

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// A deploy is a restart, and a restart must keep a device login a human started
// before it. create on one store, rebuild from the same file, and the request
// is still there to approve and poll.
func TestDevicePendingSurvivesRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "devices.json")
	now := time.Now().UTC()

	first := newDeviceStore(path, nil)
	req, err := first.create("laptop", "1.2.3.4", now)
	if err != nil {
		t.Fatal(err)
	}

	second := newDeviceStore(path, nil)
	if status, _, ok := second.poll(req.DeviceCode); !ok || status != devicePending {
		t.Fatalf("restart lost the pending request: ok=%v status=%s", ok, status)
	}
	if _, ok := second.info(req.UserCode); !ok {
		t.Fatal("restart lost the user-code lookup")
	}
}

// The minted token is the one long-lived secret in a device request and must
// never reach disk, and the raw device code is stored hashed. After approve the
// request is dropped from the persisted set, so a restart cannot replay a token
// a poll already consumed.
func TestDeviceTokenNeverPersisted(t *testing.T) {
	path := filepath.Join(t.TempDir(), "devices.json")
	now := time.Now().UTC()

	store := newDeviceStore(path, nil)
	req, err := store.create("laptop", "1.2.3.4", now)
	if err != nil {
		t.Fatal(err)
	}
	store.approve(req.UserCode, "super-secret-token")

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "super-secret-token") {
		t.Fatal("the minted token reached disk")
	}
	if strings.Contains(string(data), req.DeviceCode) {
		t.Fatal("the raw device code reached disk")
	}

	restarted := newDeviceStore(path, nil)
	if _, _, ok := restarted.poll(req.DeviceCode); ok {
		t.Fatal("an approved-and-tokenized request must not survive restart for replay")
	}
}

// The login code is kept hashed at rest by construction; a restart must keep a
// code someone is mid-paste so the exchange still succeeds.
func TestLoginCodeSurvivesRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "login-codes.json")
	now := time.Now().UTC()

	first := newLoginCodeStore(path, nil)
	if !first.create(hashToken("code"), "yann@facile.studio", scopeUser, now) {
		t.Fatal("first create refused")
	}

	second := newLoginCodeStore(path, nil)
	code, ok := second.consume(hashToken("code"), now.Add(time.Second))
	if !ok {
		t.Fatal("restart lost the login code")
	}
	if code.Email != "yann@facile.studio" || code.Scope != scopeUser {
		t.Fatalf("recovered the wrong identity: %+v", code)
	}
}

// The persisted stores are server state, exactly like tokens.json, and must not
// be reachable through the file sync that moves wiki content between machines.
func TestPersistedStoresAreFencedFromSync(t *testing.T) {
	for _, name := range []string{"devices.json", "login-codes.json"} {
		if !syncSkip(name) {
			t.Errorf("%s is not fenced by syncSkip", name)
		}
	}
}
