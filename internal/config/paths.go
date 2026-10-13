package config

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
)

// The XDG layout, per CLI-STANDARD §6.1: configuration under ~/.config, the
// syncable store under ~/.local/share, machine-local state under
// ~/.local/state. The legacy home dotfiles are what every existing machine
// carries, so they are read as a fallback and migrated once, on first use.
//
// Migration is driven lazily from DataDir and ConfigPath rather than from a
// command, because the paths are also read by the MCP server, the daemon and
// the adapters, none of which goes through a command entry point.
const (
	appDirName      = "mycelium"
	configFileName  = "config.yml"
	legacyConfig    = ".mycelium.yml"
	legacyDataDir   = ".mycelium"
	dataDirOverride = "DATA_DIR"
	stateDirEnv     = "MYCELIUM_STATE_DIR"
)

// legacyStateNames are the files that lived inside the old flat data
// directory and belong in the state directory instead: a daemon log, an
// install stamp, scan offsets, consolidation cursors, flow trust pins and
// flow run artifacts. None of them syncs, and none of them is worth
// recovering after a reinstall.
var legacyStateNames = []string{
	"daemon.log",
	".last-install",
	".sessions-state.json",
	".consolidate-cursor.json",
	".consolidate-run.json",
	".flow-trust.json",
	"runs",
}

func homeDir() string {
	home, _ := os.UserHomeDir()
	return home
}

// XDGConfigHome returns the base directory for user configuration.
func XDGConfigHome() string {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return dir
	}
	return filepath.Join(homeDir(), ".config")
}

// XDGDataHome returns the base directory for user data.
func XDGDataHome() string {
	if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
		return dir
	}
	return filepath.Join(homeDir(), ".local", "share")
}

// XDGStateHome returns the base directory for state that must survive a
// restart but is not portable enough to be called data: logs, stamps, scan
// offsets, cursors and trust pins.
func XDGStateHome() string {
	if dir := os.Getenv("XDG_STATE_HOME"); dir != "" {
		return dir
	}
	return filepath.Join(homeDir(), ".local", "state")
}

// xdgConfigPath is the pure resolver, safe to call from inside a migration.
func xdgConfigPath() string {
	return filepath.Join(XDGConfigHome(), appDirName, configFileName)
}

// xdgDataDir is the pure resolver, safe to call from inside a migration.
func xdgDataDir() string {
	if dir := os.Getenv(dataDirOverride); dir != "" {
		return dir
	}
	return filepath.Join(XDGDataHome(), appDirName)
}

// StateDirFor returns the state directory that belongs to a given data
// directory. A store placed explicitly, by DATA_DIR in a container or a temp
// dir in a test, keeps its state beside it so the whole tree stays one
// self-contained volume; the default store uses the XDG state directory like
// every other tool.
func StateDirFor(dataDir string) string {
	if dir := os.Getenv(stateDirEnv); dir != "" {
		return dir
	}
	if os.Getenv(dataDirOverride) == "" && dataDir == defaultDataDir() {
		return filepath.Join(XDGStateHome(), appDirName)
	}
	return filepath.Join(dataDir, "state")
}

func defaultDataDir() string {
	return filepath.Join(XDGDataHome(), appDirName)
}

// ConfigPath is where the client configuration file lives.
func ConfigPath() string {
	ensureLayout()
	return xdgConfigPath()
}

// LegacyConfigPath is the pre-XDG client configuration file, still read when
// the new one is absent.
func LegacyConfigPath() string {
	return filepath.Join(homeDir(), legacyConfig)
}

// DataDir is the root of the local memory store, overridable via DATA_DIR.
func DataDir() string {
	ensureLayout()
	return xdgDataDir()
}

// LegacyDataDir is the pre-XDG data directory, still used when the new one
// does not exist yet.
func LegacyDataDir() string {
	return filepath.Join(homeDir(), legacyDataDir)
}

// StateDir is where machine-local state lives. MYCELIUM_STATE_DIR overrides
// it for a container that wants a single volume.
func StateDir() string {
	ensureLayout()
	return StateDirFor(xdgDataDir())
}

// MigrateLegacyLayout moves a pre-XDG install onto the XDG directories. It is
// idempotent, it never runs when DATA_DIR points somewhere explicit, and it
// reports why it failed so doctor can show it.
func MigrateLegacyLayout() error {
	ensureLayout()
	return migrateErr
}

var (
	migrating  atomic.Bool
	migrated   atomic.Bool
	migrateErr error
)

// ensureLayout runs the migration once. The migrating flag rather than a
// sync.Once is what makes it reentrant: DataDir is read from inside the
// migration itself, and a second call must resolve the path instead of
// blocking on the Once the first call still holds.
//
// A test never migrates. The paths are read from package-level helpers that
// tests reach without isolating HOME, and moving a developer's real config
// and wiki because a test ran is not an acceptable side effect. Tests that
// want the migration exercise migrate() directly.
func ensureLayout() {
	if migrated.Load() {
		return
	}
	if testing.Testing() {
		return
	}
	if !migrating.CompareAndSwap(false, true) {
		return
	}
	defer func() {
		migrated.Store(true)
		migrating.Store(false)
	}()
	migrateErr = migrate()
}

func migrate() error {
	if err := migrateConfig(); err != nil {
		return err
	}
	if os.Getenv(dataDirOverride) != "" {
		return nil
	}
	return migrateData()
}

func migrateConfig() error {
	legacy := LegacyConfigPath()
	current := xdgConfigPath()
	if exists(current) {
		return nil
	}
	if !exists(legacy) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(current), 0o700); err != nil {
		return err
	}
	if err := copyFile(legacy, current); err != nil {
		return err
	}
	if err := os.Chmod(current, 0o600); err != nil {
		return err
	}
	return os.Remove(legacy)
}

func migrateData() error {
	legacy := LegacyDataDir()
	current := xdgDataDir()
	if exists(current) {
		return liftState(current)
	}
	if !exists(legacy) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(current), 0o700); err != nil {
		return err
	}
	if err := os.Rename(legacy, current); err != nil {
		return err
	}
	return liftState(current)
}

// liftState moves the machine-local files out of the data directory and into
// the state directory. A missing file is expected: most machines have never
// run a flow or a daemon.
func liftState(dataDir string) error {
	stateDir := StateDirFor(dataDir)
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		return err
	}
	for _, name := range legacyStateNames {
		from := filepath.Join(dataDir, name)
		info, err := os.Lstat(from)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return err
		}
		to := filepath.Join(stateDir, name)
		if exists(to) {
			continue
		}
		if info.IsDir() {
			if err := moveTree(from, to); err != nil {
				return err
			}
			continue
		}
		if err := os.Rename(from, to); err != nil {
			return err
		}
	}
	return nil
}

func exists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}

// moveTree renames a directory, falling back to a copy-and-remove when the
// state directory sits on another filesystem.
func moveTree(from, to string) error {
	if err := os.Rename(from, to); err == nil {
		return nil
	}
	return copyTree(from, to)
}

func copyTree(from, to string) error {
	if err := os.MkdirAll(to, 0o700); err != nil {
		return err
	}
	entries, err := os.ReadDir(from)
	if err != nil {
		return err
	}
	for _, e := range entries {
		src := filepath.Join(from, e.Name())
		dst := filepath.Join(to, e.Name())
		if e.IsDir() {
			if err := copyTree(src, dst); err != nil {
				return err
			}
			continue
		}
		if err := copyFile(src, dst); err != nil {
			return err
		}
	}
	return os.RemoveAll(from)
}

func copyFile(from, to string) error {
	src, err := os.Open(from)
	if err != nil {
		return err
	}
	defer src.Close()
	info, err := src.Stat()
	if err != nil {
		return err
	}
	dst, err := os.OpenFile(to, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode().Perm())
	if err != nil {
		return err
	}
	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		return err
	}
	return dst.Close()
}
