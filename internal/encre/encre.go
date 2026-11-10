// Package encre judges artifact bodies with encre, the reference
// implementation of the Facile document format, through its Go module
// (github.com/FacileStudio/encre/go), which runs the encre binary as a
// subprocess. Artifacts are markdown in encre's lossless dialect, and encre
// is the one judge of whether a body is a valid document.
package encre

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	upstream "github.com/FacileStudio/encre/go"
)

// BinaryEnv names the environment variable that points at the encre binary,
// the same one the module reads. Unset, the binary is looked up on PATH as
// "encre".
const BinaryEnv = "ENCRE_BIN"

// ErrUnavailable is returned when no encre binary can be found.
var ErrUnavailable = errors.New("encre is not installed")

// Problem is one reason a document is invalid, as encre reports it. Line and
// Column are 1-based and name where the block holding the problem starts;
// they are zero when encre refused the input before reading it.
type Problem = upstream.Problem

// ValidationError is encre's verdict on a body it refuses.
type ValidationError struct {
	Problems []Problem
}

func (e *ValidationError) Error() string {
	messages := make([]string, 0, len(e.Problems))
	for _, p := range e.Problems {
		messages = append(messages, describe(p))
	}
	return "invalid document: " + strings.Join(messages, "; ")
}

// describe places the message in the source when encre knows where, which
// is what the API, the CLI and the MCP tool show.
func describe(p Problem) string {
	if p.Line > 0 {
		return fmt.Sprintf("line %d, column %d: %s", p.Line, p.Column, p.Message)
	}
	return p.Message
}

// Available reports whether an encre binary can be found.
func Available() bool {
	_, err := exec.LookPath(upstream.Binary())
	return err == nil
}

// ValidateMarkdown runs `encre validate --from markdown --json` on body. It
// answers nil for a valid document, a *ValidationError for one encre refuses,
// ErrUnavailable when encre is not installed, and another error when encre
// could not run.
func ValidateMarkdown(ctx context.Context, body []byte) error {
	if !Available() {
		return ErrUnavailable
	}
	verdict, err := upstream.Validate(ctx, body, upstream.Markdown)
	if err != nil {
		return refusal(err)
	}
	if verdict.Valid {
		return nil
	}
	return &ValidationError{Problems: verdict.Errors}
}

// refusal reads a failed run. An input encre cannot read (not UTF-8, too
// large) is a refusal of the document, without a position; anything else is
// encre failing.
func refusal(err error) error {
	var failed *upstream.Error
	if errors.As(err, &failed) && failed.Invalid() {
		return &ValidationError{Problems: []Problem{{Path: failed.Path, Message: failed.Message}}}
	}
	return fmt.Errorf("encre validate: %w", err)
}
