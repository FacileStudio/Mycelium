package server

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

const testIssuer = "https://sso.example"

func mustLink(t *testing.T, srv *Server, email, name string) User {
	t.Helper()
	user, err := srv.linkIdentity(context.Background(), oidcIdentity{
		Issuer: testIssuer, Subject: "sub-" + email, Email: email, EmailVerified: true, Name: name,
	})
	if err != nil {
		t.Fatalf("link %s: %v", email, err)
	}
	return user
}

func TestLinkRefusesAnUnverifiedEmail(t *testing.T) {
	srv := New(t.TempDir(), "pw")
	mustLink(t, srv, "owner@facile.studio", "Owner")
	_, err := srv.linkIdentity(context.Background(), oidcIdentity{Issuer: testIssuer, Subject: "squat", Email: "owner@facile.studio"})
	if err != errUnverifiedEmail {
		t.Fatalf("an unverified email linked: %v", err)
	}
}

func TestLinkRefusesASecondSubjectOfTheSameIssuer(t *testing.T) {
	srv := New(t.TempDir(), "pw")
	mustLink(t, srv, "owner@facile.studio", "Owner")
	_, err := srv.linkIdentity(context.Background(), oidcIdentity{
		Issuer: testIssuer, Subject: "recycled", Email: "Owner@facile.studio", EmailVerified: true,
	})
	if err != errSubjectConflict {
		t.Fatalf("a recycled address took over a linked account: %v", err)
	}
}

func TestLinkAdoptsAnAccountWithoutSubjectByVerifiedEmail(t *testing.T) {
	srv := New(t.TempDir(), "pw")
	users := map[string]User{"Old@facile.studio": {Email: "Old@facile.studio", Name: "Old", Admin: true}}
	if err := srv.saveUsers(users); err != nil {
		t.Fatal(err)
	}
	user, err := srv.linkIdentity(context.Background(), oidcIdentity{
		Issuer: testIssuer, Subject: "s1", Email: "old@facile.studio", EmailVerified: true, Locale: "de-AT",
	})
	if err != nil || !user.Admin || user.Subject != "s1" || user.Locale != "de" {
		t.Fatalf("adopt: %+v %v", user, err)
	}
	if again := mustLinkSubject(t, srv, "s1"); again.Email != "Old@facile.studio" {
		t.Fatalf("second login did not resolve by subject: %+v", again)
	}
}

func mustLinkSubject(t *testing.T, srv *Server, subject string) User {
	t.Helper()
	user, err := srv.linkIdentity(context.Background(), oidcIdentity{Issuer: testIssuer, Subject: subject, Email: "other@x.test"})
	if err != nil {
		t.Fatal(err)
	}
	return user
}

func TestErasePersonIsScopedAndIdempotent(t *testing.T) {
	srv := New(t.TempDir(), "pw")
	mustLink(t, srv, "admin@facile.studio", "Admin")
	gone := mustLink(t, srv, "gone@facile.studio", "Gone")
	keep := mustLink(t, srv, "keep@facile.studio", "Keep")
	srv.mu.Lock()
	srv.tokens["t1"] = TokenInfo{Hash: "t1", Name: "session", Scope: scopeUser, UserEmail: gone.Email}
	srv.tokens["t2"] = TokenInfo{Hash: "t2", Name: "session", Scope: scopeUser, UserEmail: keep.Email}
	srv.saveSpaces(map[string]*Space{
		"solo":   {ID: "solo", Members: map[string]string{gone.Email: roleOwner}},
		"shared": {ID: "shared", Members: map[string]string{gone.Email: roleOwner, keep.Email: roleMember}},
	})
	srv.mu.Unlock()
	solo := srv.spaceDir("solo")
	if err := os.MkdirAll(solo+"/memory", 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(solo+"/memory/a.md", []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	for range 2 {
		if err := srv.erasePerson(context.Background(), testIssuer, "sub-"+gone.Email); err != nil {
			t.Fatal(err)
		}
	}

	users := srv.loadUsers()
	if _, ok := users[gone.Email]; ok || len(users) != 2 {
		t.Fatalf("users after erasure: %v", users)
	}
	if _, ok := srv.tokens["t1"]; ok {
		t.Error("the erased person's token survived")
	}
	if _, ok := srv.tokens["t2"]; !ok {
		t.Error("another person's token was removed")
	}
	spaces := srv.loadSpaces()
	if _, ok := spaces["solo"]; ok || pathExists(solo) {
		t.Error("a space only the erased person used survived")
	}
	if spaces["shared"].Members[keep.Email] != roleOwner || len(spaces["shared"].Members) != 1 {
		t.Errorf("shared space after erasure: %v", spaces["shared"].Members)
	}
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func TestConcurrentLinksNeverCorruptTheUsersFile(t *testing.T) {
	srv := New(t.TempDir(), "")
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			email := fmt.Sprintf("u%d@example.com", i)
			_, err := srv.linkIdentity(context.Background(), oidcIdentity{Issuer: "i", Subject: email, Email: email, EmailVerified: true})
			if err != nil {
				t.Error(err)
			}
		}(i)
	}
	wg.Wait()
	if got := len(srv.loadUsers()); got != 20 {
		t.Fatalf("users = %d, want 20", got)
	}
	matches, _ := filepath.Glob(filepath.Join(srv.DataDir, "*.tmp"))
	if len(matches) != 0 {
		t.Fatalf("leftover temp files: %v", matches)
	}
}
