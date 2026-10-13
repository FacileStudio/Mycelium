package server

import "testing"

// TestServerSkipsMachineLocalStateToo keeps the two copies of the state rule
// in step. The container runs with DATA_DIR=/data, so its state directory sits
// inside the tree it serves; if the server could see it and the client could
// not, the reconcile would read every scan offset as a deletion.
func TestServerSkipsMachineLocalStateToo(t *testing.T) {
	for _, path := range []string{
		"state/.sessions-state.json",
		"state/.flow-trust.json",
		"state/runs/flow/x.json",
		"state/daemon.log",
	} {
		if !syncSkip(path) {
			t.Errorf("%q would sync", path)
		}
	}
	if syncSkip("memory/state.md") {
		t.Error("a page with the word in its name must still sync")
	}
}

// TestServerSkipsInstalledPackagesToo keeps the two copies of this rule in
// step. A path the client excludes and the server does not is a file the
// reconcile sees on one side only, which it reads as a deletion.
func TestServerSkipsInstalledPackagesToo(t *testing.T) {
	for _, path := range []string{
		"node_modules/left-pad/index.js",
		"extensions/models/node_modules/zod/package.json",
	} {
		if !syncSkip(path) {
			t.Errorf("%q would sync", path)
		}
	}
	if syncSkip("memory/tools/my-node_modules.md") {
		t.Error("a page with the word in its name must still sync")
	}
}
