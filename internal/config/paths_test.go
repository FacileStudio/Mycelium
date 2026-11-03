package config

import (
	"os"
	"path/filepath"
	"testing"
)

// The migration moves a user's wiki and their token, and it runs itself on the
// first command anybody types. A mistake here is invisible until a machine
// comes up empty, so the four properties that make it safe are pinned: the
// config moves with its contents, the state files leave the synced tree, a
// second run changes nothing, and DATA_DIR keeps the migration off entirely.
func TestMigrateLegacyLayout(t *testing.T) {
	home := t.TempDir()
	setLegacyLayout(t, home)
	writeLegacyInstall(t, home)

	if err := migrate(); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	moved, err := os.ReadFile(ConfigPath())
	if err != nil {
		t.Fatalf("config did not move: %v", err)
	}
	if string(moved) != "url: https://legacy.example/\ntoken: keep-me\n" {
		t.Fatalf("config contents changed: %q", moved)
	}
	if _, err := os.Stat(LegacyConfigPath()); !os.IsNotExist(err) {
		t.Error("the legacy config file was left behind")
	}
	info, err := os.Stat(ConfigPath())
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm()&0o077 != 0 {
		t.Fatalf("config is mode %v, want no group or other bits", info.Mode().Perm())
	}
	if _, err := os.Stat(filepath.Join(DataDir(), "memory", "index.md")); err != nil {
		t.Fatalf("the memory tree did not move: %v", err)
	}
	for _, stateFile := range []string{"daemon.log", ".sessions-state.json", ".flow-trust.json"} {
		if _, err := os.Stat(filepath.Join(StateDir(), stateFile)); err != nil {
			t.Errorf("%s did not reach the state directory: %v", stateFile, err)
		}
		if _, err := os.Stat(filepath.Join(DataDir(), stateFile)); !os.IsNotExist(err) {
			t.Errorf("%s was left in the synced tree, where it would travel", stateFile)
		}
	}
	if _, err := os.Stat(filepath.Join(DataDir(), "events", "ruche", "2026-10.jsonl")); err != nil {
		t.Errorf("synced content did not move with the tree: %v", err)
	}
}

// Idempotence is what makes this safe to run from a lazy call site. The second
// pass must not re-copy, re-rename or report an error now that nothing legacy
// is left.
func TestMigrateLegacyLayoutIsIdempotent(t *testing.T) {
	home := t.TempDir()
	setLegacyLayout(t, home)
	writeLegacyInstall(t, home)

	if err := migrate(); err != nil {
		t.Fatalf("first migrate: %v", err)
	}
	before, err := os.ReadFile(ConfigPath())
	if err != nil {
		t.Fatal(err)
	}
	if err := migrate(); err != nil {
		t.Fatalf("second migrate: %v", err)
	}
	after, err := os.ReadFile(ConfigPath())
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatalf("the second run rewrote the config:\nbefore %q\nafter  %q", before, after)
	}
}

// DATA_DIR names a store that belongs to something else — a container volume, a
// test temp dir. Moving a legacy tree into it would mix two installs, so the
// data migration stands down and only the config path resolves as usual.
func TestMigrateLegacyLayoutLeavesAnExplicitDataDirAlone(t *testing.T) {
	home := t.TempDir()
	setLegacyLayout(t, home)
	writeLegacyInstall(t, home)
	explicit := t.TempDir()
	t.Setenv(dataDirOverride, explicit)

	if err := migrate(); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if DataDir() != explicit {
		t.Fatalf("DataDir() = %q, want the explicit %q", DataDir(), explicit)
	}
	if _, err := os.Stat(filepath.Join(explicit, "memory", "index.md")); !os.IsNotExist(err) {
		t.Error("the legacy tree was moved into the explicit DATA_DIR")
	}
	if _, err := os.Stat(LegacyDataDir()); err != nil {
		t.Errorf("the legacy tree should be untouched, got %v", err)
	}
}

func setLegacyLayout(t *testing.T, home string) {
	t.Helper()
	IsolateHome(t, home)
	t.Setenv(dataDirOverride, "")
	t.Setenv(stateDirEnv, "")
}

func writeLegacyInstall(t *testing.T, home string) {
	t.Helper()
	writeFixture(t, filepath.Join(home, legacyConfig), "url: https://legacy.example/\ntoken: keep-me\n")
	for path, body := range map[string]string{
		filepath.Join(legacyDataDir, "memory", "index.md"):               "# Index",
		filepath.Join(legacyDataDir, "events", "ruche", "2026-10.jsonl"): `{"t":"x"}`,
		filepath.Join(legacyDataDir, "daemon.log"):                       "started\n",
		filepath.Join(legacyDataDir, ".sessions-state.json"):             `{"files":{}}`,
		filepath.Join(legacyDataDir, ".flow-trust.json"):                 `{}`,
	} {
		writeFixture(t, filepath.Join(home, path), body)
	}
}

func writeFixture(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}
