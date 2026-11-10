package server

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"testing"
)

func listedTrees(t *testing.T, f delegatedFixture, token string) map[string]string {
	t.Helper()
	var items []ArtifactSummary
	rec := expect(t, f, apiCall{Method: "GET", Path: "/api/artifacts?all=true", Token: token}, http.StatusOK)
	if err := json.Unmarshal(rec.Body.Bytes(), &items); err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	for _, item := range items {
		out[item.Title] = item.SpaceID
	}
	return out
}

func postArtifact(t *testing.T, f delegatedFixture, token, query, title string) int {
	t.Helper()
	body := `{"title":"` + title + `","content":"# x\n","machine":"bulle"}`
	return spReq(t, f.h, apiCall{Method: "POST", Path: "/api/artifacts" + query, Token: token, Body: body}).Code
}

func TestDelegatedListsEveryReadableArtifact(t *testing.T) {
	f := newDelegatedFixture(t, "mycelium")
	other := f.issuer.tokenFor(t, otherEmail, "mycelium", "mycelium.read mycelium.write")
	codes := []int{
		postArtifact(t, f, f.rw, "", "Common"),
		postArtifact(t, f, f.rw, "?space_id="+f.own, "Own"),
		postArtifact(t, f, other, "?space_id="+f.foreign, "Foreign"),
		postArtifact(t, f, other, "", "Nope"),
	}
	if want := []int{201, 201, 201, 403}; !slices.Equal(codes, want) {
		t.Fatalf("creates answered %v, want %v", codes, want)
	}
	admin := listedTrees(t, f, f.ro)
	if len(admin) != 2 || admin["Common"] != "" || admin["Own"] != f.own {
		t.Fatalf("the admin sees the common tree and their spaces only: %v", admin)
	}
	if got := listedTrees(t, f, other); len(got) != 1 || got["Foreign"] != f.foreign {
		t.Fatalf("a non-admin never sees the common tree: %v", got)
	}
}

func TestSpacesSayWhetherTheCommonTreeHoldsTheirArtifacts(t *testing.T) {
	f := newDelegatedFixture(t, "mycelium")
	other := f.issuer.tokenFor(t, otherEmail, "mycelium", "mycelium.read")
	for token, want := range map[string]string{f.ro: `"common_artifacts":true`, other: `"common_artifacts":false`} {
		rec := expect(t, f, apiCall{Method: "GET", Path: "/api/spaces", Token: token}, http.StatusOK)
		if !strings.Contains(rec.Body.String(), want) {
			t.Fatalf("want %s: %s", want, rec.Body.String())
		}
	}
}
