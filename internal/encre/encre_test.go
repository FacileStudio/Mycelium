package encre_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/FacileStudio/Mycelium/internal/encre"
	"github.com/FacileStudio/Mycelium/internal/encre/encretest"
)

func TestValidateMarkdownPlacesEncreRefusals(t *testing.T) {
	encretest.Install(t)
	ctx := context.Background()
	if err := encre.ValidateMarkdown(ctx, []byte("# Fine\n")); err != nil {
		t.Fatalf("valid body: %v", err)
	}
	var invalid *encre.ValidationError
	err := encre.ValidateMarkdown(ctx, []byte("# INVALID\n"))
	if !errors.As(err, &invalid) {
		t.Fatalf("a verdict on stdout must become a ValidationError, got %v", err)
	}
	want := "invalid document: line 3, column 1: text is empty or holds a forbidden character"
	if err.Error() != want {
		t.Fatalf("got %q, want %q", err.Error(), want)
	}
	if p := invalid.Problems[0]; p.Path != "doc.content[1]" || p.Line != 3 || p.Column != 1 {
		t.Fatalf("the position is lost: %+v", p)
	}
	err = encre.ValidateMarkdown(ctx, []byte("OVERSIZE"))
	want = "invalid document: input is larger than 8388608 bytes; send less or raise --max-bytes"
	if !errors.As(err, &invalid) || err.Error() != want {
		t.Fatalf("an input error on stderr must become a ValidationError without a position, got %v", err)
	}
	err = encre.ValidateMarkdown(ctx, []byte("CRASH"))
	if errors.As(err, &invalid) || err == nil || !strings.Contains(err.Error(), "exit 3") {
		t.Fatalf("an internal failure is not a refusal, got %v", err)
	}
}

func TestValidateMarkdownWithoutEncre(t *testing.T) {
	t.Setenv(encre.BinaryEnv, "/nonexistent/encre")
	if encre.Available() {
		t.Fatal("a missing binary must not be available")
	}
	if err := encre.ValidateMarkdown(context.Background(), []byte("x")); !errors.Is(err, encre.ErrUnavailable) {
		t.Fatalf("got %v, want ErrUnavailable", err)
	}
}

func TestRealEncrePlacesRefusals(t *testing.T) {
	if !encre.Available() {
		t.Skip("encre is not installed")
	}
	ctx := context.Background()
	body := "# Hi\n\n```mermaid\nflowchart LR\n  A-->B\n```\n\nA [badge:success kept] ==mark== [status:pending] line.\n"
	if err := encre.ValidateMarkdown(ctx, []byte(body)); err != nil {
		t.Fatalf("the schema 2 constructs are valid: %v", err)
	}
	if err := encre.ValidateMarkdown(ctx, []byte(encretest.NestedCallouts(8))); err != nil {
		t.Fatalf("eight nested custom blocks are allowed: %v", err)
	}
	var invalid *encre.ValidationError
	err := encre.ValidateMarkdown(ctx, []byte("# Hi\n\na\x01b\n"))
	want := "invalid document: line 3, column 1: text is empty or holds a forbidden character"
	if !errors.As(err, &invalid) || err.Error() != want {
		t.Fatalf("a control character must be refused at its line, got %v", err)
	}
	err = encre.ValidateMarkdown(ctx, []byte(encretest.NestedCallouts(9)))
	if !errors.As(err, &invalid) || !strings.HasPrefix(err.Error(), "invalid document: line 11, column 1: custom blocks nest") {
		t.Fatalf("nine nested custom blocks must be refused at the ninth fence, got %v", err)
	}
	if p := invalid.Problems[0]; p.Path != "doc" || p.Line != 11 || p.Column != 1 {
		t.Fatalf("the position is lost: %+v", p)
	}
}
