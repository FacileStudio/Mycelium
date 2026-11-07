package server

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/FacileStudio/tronc/httpjson"
)

// writeStatusError answers with the suite's error envelope at a status
// tronc's code-to-status map does not cover, so 405 and 503 stay themselves
// instead of collapsing into the generic 500 WriteError would produce.
func writeStatusError(w http.ResponseWriter, status int, code, message string) {
	httpjson.WriteJSON(w, status, map[string]any{
		"error": map[string]string{"code": code, "message": message},
	})
}

// pathParam reads a chi URL parameter and unescapes it.
//
// chi matches on the raw path whenever the request carries any percent-encoding
// and hands the parameter back still encoded, where net/http's ServeMux — which
// this router replaced — decoded it. Without this, `%2F` would sail past the
// traversal guards and an encodeURIComponent'd member email would never match a
// stored one.
func pathParam(r *http.Request, key string) string {
	raw := chi.URLParam(r, key)
	decoded, err := url.PathUnescape(raw)
	if err != nil {
		return raw
	}
	return decoded
}

func safeName(name string) (string, bool) {
	if name == "" || strings.ContainsAny(name, `/\`) || strings.Contains(name, "..") {
		return "", false
	}
	return name, true
}

func writeFile(dir, name string, body io.Reader) error {
	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, name), data, 0o644)
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func listMdNames(dir string) []string {
	return listNamesWithExt(dir, ".md")
}

func listNamesWithExt(dir, ext string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return []string{}
	}
	names := []string{}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ext) {
			names = append(names, strings.TrimSuffix(e.Name(), ext))
		}
	}
	sort.Strings(names)
	return names
}

// readStore decodes one JSON store file. A missing file is an empty store; an
// unreadable or corrupt one is an error, so a caller that writes the store
// back never replaces data it failed to read with an empty map.
func readStore[T any](path string) (map[string]T, error) {
	store := make(map[string]T)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	if err := json.Unmarshal(data, &store); err != nil {
		return make(map[string]T), err
	}
	if store == nil {
		store = make(map[string]T)
	}
	return store, nil
}
