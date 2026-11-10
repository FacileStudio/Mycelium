package server

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/FacileStudio/Mycelium/internal/encre"
	"github.com/FacileStudio/Mycelium/internal/encre/encretest"
)

func TestArtifactWritesRefuseWhatEncreRefuses(t *testing.T) {
	encretest.Install(t)
	f := newArtifactFixture(t)
	first := f.create(t, `{"title":"Fine","content":"# Fine\n\nA valid body."}`)

	rec := spReq(t, f.h, apiCall{Method: "POST", Path: f.path(""), Token: f.owner, Body: `{"title":"Bad","content":"# Bad\n\nINVALID"}`})
	if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "forbidden character") {
		t.Fatalf("POST of an invalid document: got %d %s, want 400 with encre's reason", rec.Code, rec.Body.String())
	}
	rec = spReq(t, f.h, apiCall{Method: "PUT", Path: f.path("/" + first.ID), Token: f.owner, Body: `{"content":"INVALID"}`})
	if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "invalid document") {
		t.Fatalf("PUT of an invalid document: got %d %s, want 400 with encre's reason", rec.Code, rec.Body.String())
	}
	rec = spReq(t, f.h, apiCall{Method: "GET", Path: f.path("/" + first.ID), Token: f.owner})
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "A valid body.") {
		t.Fatalf("a refused update must leave the artifact as it was: %d %s", rec.Code, rec.Body.String())
	}
	rec = spReq(t, f.h, apiCall{Method: "PUT", Path: f.path("/" + first.ID), Token: f.owner, Body: `{"content":"Still fine."}`})
	if rec.Code != http.StatusOK {
		t.Fatalf("PUT of a valid document: %d %s", rec.Code, rec.Body.String())
	}
}

func TestArtifactWritesPlaceEncreRefusals(t *testing.T) {
	encretest.Install(t)
	f := newArtifactFixture(t)
	rec := spReq(t, f.h, apiCall{Method: "POST", Path: f.path(""), Token: f.owner, Body: `{"title":"Deep","content":"# Deep\n\nNESTED"}`})
	want := "invalid document: line 11, column 1: custom blocks nest more than 8 levels deep"
	if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), want) {
		t.Fatalf("got %d %s, want 400 carrying %q", rec.Code, rec.Body.String(), want)
	}
}

var placedRefusal = regexp.MustCompile(`invalid document: line \d+, column \d+: custom blocks nest`)

func TestRealEncreRefusesDeepNestingWithAPosition(t *testing.T) {
	if !encre.Available() {
		t.Skip("encre is not installed")
	}
	f := newArtifactFixture(t)
	body, _ := json.Marshal(map[string]string{"title": "Deep", "content": encretest.NestedCallouts(9)})
	rec := spReq(t, f.h, apiCall{Method: "POST", Path: f.path(""), Token: f.owner, Body: string(body)})
	if rec.Code != http.StatusBadRequest || !placedRefusal.MatchString(rec.Body.String()) {
		t.Fatalf("nine nested custom blocks: got %d %s, want 400 placed at a line and column", rec.Code, rec.Body.String())
	}
	body, _ = json.Marshal(map[string]string{"title": "Fine", "content": encretest.NestedCallouts(8)})
	if rec = spReq(t, f.h, apiCall{Method: "POST", Path: f.path(""), Token: f.owner, Body: string(body)}); rec.Code != http.StatusCreated {
		t.Fatalf("eight nested custom blocks: got %d %s, want 201", rec.Code, rec.Body.String())
	}
}
