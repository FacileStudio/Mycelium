// Package encretest installs a stand-in encre binary for tests, so the write
// paths can be exercised without the real tool on the machine.
package encretest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/FacileStudio/Mycelium/internal/encre"
)

// LogEnv names the file the fake appends every input it receives to.
const LogEnv = "FAKE_ENCRE_LOG"

const separator = "\n---\n"

const script = `#!/bin/sh
input=$(cat)
if [ -n "${FAKE_ENCRE_LOG:-}" ]; then printf '%s\n---\n' "$input" >> "$FAKE_ENCRE_LOG"; fi
case "$input" in
  *INVALID*) printf '%s\n' '{"errors":[{"message":"text is empty or holds a forbidden character","path":"doc.content[1]","line":3,"column":1}],"valid":false}'; exit 1 ;;
  *NESTED*) printf '%s\n' '{"errors":[{"message":"custom blocks nest more than 8 levels deep","path":"doc","line":11,"column":1}],"valid":false}'; exit 1 ;;
  *OVERSIZE*) printf '%s\n' '{"error":{"code":"invalid","exit":1,"message":"input is larger than 8388608 bytes; send less or raise --max-bytes","path":null}}' >&2; exit 1 ;;
  *CRASH*) printf '%s\n' '{"error":{"code":"internal","exit":3,"message":"panic","path":null}}' >&2; exit 3 ;;
esac
printf '%s\n' '{"errors":[],"valid":true}'
`

// Install writes the fake and points ENCRE_BIN at it for the rest of the
// test. The fake answers the verdicts encre 0.9.0 gives: any input holding
// INVALID is refused at line 3, column 1 as a forbidden character, any input
// holding NESTED at line 11, column 1 as custom blocks nested too deep, any
// input holding OVERSIZE with an input error, and CRASH exits 3. It answers
// the path of the log file every input is appended to.
func Install(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	bin := filepath.Join(dir, "encre")
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	log := filepath.Join(dir, "inputs.log")
	t.Setenv(encre.BinaryEnv, bin)
	t.Setenv(LogEnv, log)
	return log
}

// Inputs answers everything the fake has validated so far, one entry per call.
func Inputs(t *testing.T, log string) []string {
	t.Helper()
	raw, err := os.ReadFile(log)
	if err != nil {
		return nil
	}
	return strings.Split(strings.TrimSuffix(string(raw), separator), separator)
}

// NestedCallouts writes a document whose callout fences nest depth levels
// deep, each fence one backtick longer than the one inside it.
func NestedCallouts(depth int) string {
	inner := "x\n"
	for i := 0; i < depth; i++ {
		ticks := strings.Repeat("`", 3+i)
		inner = ticks + "callout tone=\"info\"\n" + inner + ticks + "\n"
	}
	return "# Deep\n\n" + inner
}
