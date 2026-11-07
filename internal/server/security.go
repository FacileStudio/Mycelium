package server

import (
	"net/url"
	"strings"

	"github.com/FacileStudio/tronc/middleware"
)

const iconifyAPI = "https://api.iconify.design"

func (s *Server) securityConfig() middleware.SecurityConfig {
	connect := []string{iconifyAPI}
	if origin := originOf(s.JournalBrowserURL); origin != "" {
		connect = append(connect, origin)
	}
	sources := map[string][]string{"connect-src": connect}
	if len(s.ScriptSources) > 0 {
		sources["script-src"] = s.ScriptSources
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
