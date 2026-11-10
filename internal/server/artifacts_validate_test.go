package server

import (
	"net/http"
	"strings"
	"testing"

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
