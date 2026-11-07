package server

import (
	"net/http"
	"net/netip"
	"net/url"
	"strings"

	"github.com/FacileStudio/tronc/httpjson"
	"github.com/FacileStudio/tronc/middleware"
)

type EdgeConfig struct {
	TrustedProxies []netip.Prefix
	CDNProxies     []netip.Prefix
	CDNHeader      string
	ScriptSources  []string
}

const iconifyAPI = "https://api.iconify.design"

func (s *Server) securityConfig() middleware.SecurityConfig {
	connect := []string{iconifyAPI}
	if origin := originOf(s.JournalBrowserURL); origin != "" {
		connect = append(connect, origin)
	}
	sources := map[string][]string{"connect-src": connect}
	if len(s.Edge.ScriptSources) > 0 {
		sources["script-src"] = s.Edge.ScriptSources
	}
	return middleware.SecurityConfig{CSPSources: sources}
}

func originOf(raw string) string {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ""
	}
	return parsed.Scheme + "://" + parsed.Host
}

func (s *Server) versionHandler(w http.ResponseWriter, _ *http.Request) {
	version := s.Version
	if version == "" {
		version = "dev"
	}
	httpjson.WriteJSON(w, http.StatusOK, map[string]string{"version": version})
}
