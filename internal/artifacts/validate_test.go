package artifacts

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/FacileStudio/Mycelium/internal/encre"
	"github.com/FacileStudio/Mycelium/internal/encre/encretest"
)

func TestWritesAreValidatedUnderTheFrontmatter(t *testing.T) {
	log := encretest.Install(t)
	dir := t.TempDir()
	now := time.Date(2026, 11, 10, 10, 0, 0, 0, time.UTC)
	var invalid *encre.ValidationError

	_, err := AddContent(dir, []byte("# Bad\n\nINVALID text\n"), "stdin.md", Request{Title: "Bad"}, now)
	if !errors.As(err, &invalid) || !strings.Contains(err.Error(), "forbidden character") {
		t.Fatalf("an invalid body must be refused with encre's reason, got %v", err)
	}
	if all, _ := List(dir); len(all) != 0 {
		t.Fatalf("a refused body was written: %+v", all)
	}

	body := "---\ntitle: \"Mine\"\nmachine: \"lucy\"\n---\n\n# Plan\n\nStep one.\n"
	art, err := AddContent(dir, []byte(body), "plan.md", Request{Machine: "ruche"}, now)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	stored, _ := os.ReadFile(art.Path)
	if !strings.HasPrefix(string(stored), "---\ntitle: \"Mine\"\nmachine: \"ruche\"\n") || !strings.HasSuffix(string(stored), "# Plan\n\nStep one.\n") {
		t.Fatalf("the frontmatter was not rewritten around the untouched body:\n%s", stored)
	}
	inputs := encretest.Inputs(t, log)
	if len(inputs) != 2 || strings.TrimSpace(inputs[1]) != "# Plan\n\nStep one." {
		t.Fatalf("encre must validate the body alone, got %q", inputs)
	}

	_, err = Update(dir, art.ID, "", []byte("INVALID again"), now.Add(time.Minute))
	if !errors.As(err, &invalid) {
		t.Fatalf("an invalid update must be refused, got %v", err)
	}
	if again, _ := os.ReadFile(art.Path); string(again) != string(stored) {
		t.Fatal("a refused update changed the file")
	}
	if _, err := Update(dir, art.ID, "", []byte("---\ntitle: \"x\"\n---\n\nStep two.\n"), now.Add(time.Minute)); err != nil {
		t.Fatalf("update: %v", err)
	}
	if inputs := encretest.Inputs(t, log); strings.TrimSpace(inputs[len(inputs)-1]) != "Step two." {
		t.Fatalf("the update was validated with its frontmatter: %q", inputs[len(inputs)-1])
	}
}

func TestOversizedBodiesAreRefusedBeforeEncreRuns(t *testing.T) {
	log := encretest.Install(t)
	dir := t.TempDir()
	now := time.Date(2026, 11, 10, 10, 0, 0, 0, time.UTC)
	huge := append([]byte("# Big\n\n"), bytes.Repeat([]byte("a"), MaxBodyBytes)...)

	if _, err := AddContent(dir, huge, "big.md", Request{Title: "Big"}, now); !errors.Is(err, ErrTooLarge) {
		t.Fatalf("a body over the cap must be refused, got %v", err)
	}
	if inputs := encretest.Inputs(t, log); len(inputs) != 0 {
		t.Fatalf("encre ran on a body over the cap: %q", inputs)
	}
	if all, _ := List(dir); len(all) != 0 {
		t.Fatalf("a refused body was written: %+v", all)
	}
	art, err := AddContent(dir, huge[:MaxBodyBytes], "big.md", Request{Title: "Big"}, now)
	if err != nil {
		t.Fatalf("a body at the cap is validated and written: %v", err)
	}
	if _, err := Update(dir, art.ID, "", huge, now.Add(time.Minute)); !errors.Is(err, ErrTooLarge) {
		t.Fatalf("an oversized update must be refused, got %v", err)
	}
	t.Setenv(encre.BinaryEnv, "/nonexistent/encre")
	if _, err := AddContent(dir, huge, "big.md", Request{Title: "Bigger"}, now); !errors.Is(err, ErrTooLarge) {
		t.Fatalf("the cap must hold without encre, got %v", err)
	}
}

func TestWritesGoThroughWithoutEncre(t *testing.T) {
	t.Setenv(encre.BinaryEnv, "/nonexistent/encre")
	dir := t.TempDir()
	art, err := AddContent(dir, []byte("# Unchecked\n"), "stdin.md", Request{}, time.Now())
	if err != nil {
		t.Fatalf("without encre the body is stored as it came: %v", err)
	}
	if _, err := Find(dir, art.ID); err != nil {
		t.Fatal(err)
	}
}
