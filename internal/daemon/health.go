package daemon

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	timerUnit       = "mycelium-sync.timer"
	staleIntervals  = 10
	systemdTimeForm = "Mon 2006-01-02 15:04:05 MST"
	reinstallHint   = " — run 'mycelium daemon install'"
)

// Health reports whether the background service is installed and actually
// scheduling runs. An installed unit file proves nothing on its own: a
// systemd timer can stay active with no next run.
func Health(now time.Time) (string, bool) {
	if !Installed() {
		return "not installed" + reinstallHint, false
	}
	switch runtime.GOOS {
	case "linux":
		out, err := exec.Command("systemctl", "--user", "show", timerUnit,
			"-p", "NextElapseUSecRealtime", "-p", "NextElapseUSecMonotonic", "-p", "LastTriggerUSec").Output()
		if err != nil {
			return fmt.Sprintf("installed, timer state unknown (%v)", err), true
		}
		return timerHealth(string(out), now)
	case "darwin":
		if err := exec.Command("launchctl", "list", Label).Run(); err != nil {
			return "installed but not loaded in launchd" + reinstallHint, false
		}
		return "installed and loaded", true
	default:
		return "installed", true
	}
}

// timerHealth reads `systemctl show` output for the sync timer: it must have
// a next run, and its last run must be no older than staleIntervals runs.
func timerHealth(show string, now time.Time) (string, bool) {
	props := map[string]string{}
	for _, line := range strings.Split(show, "\n") {
		if key, value, ok := strings.Cut(strings.TrimSpace(line), "="); ok {
			props[key] = strings.TrimSpace(value)
		}
	}
	if !scheduled(props["NextElapseUSecRealtime"]) && !scheduled(props["NextElapseUSecMonotonic"]) {
		return "installed but the timer has no next run" + reinstallHint, false
	}
	last := props["LastTriggerUSec"]
	if !scheduled(last) {
		return "running, first run pending", true
	}
	at, err := time.ParseInLocation(systemdTimeForm, last, time.Local)
	if err != nil {
		return "running, last run " + last, true
	}
	age := now.Sub(at).Truncate(time.Second)
	if age > staleIntervals*IntervalSeconds*time.Second {
		return fmt.Sprintf("installed but last ran %s ago", age) + reinstallHint, false
	}
	return fmt.Sprintf("running, last ran %s ago", age), true
}

func scheduled(value string) bool {
	return value != "" && value != "0" && value != "n/a" && value != "infinity"
}
