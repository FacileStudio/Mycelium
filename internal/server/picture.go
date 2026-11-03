package server

import (
	"net/url"
	"strings"
)

// pictureURL keeps an OIDC picture claim only when it is an https URL, so a
// data: or http: value from the provider never reaches an <img> in the
// dashboard. Registre's claim is stable per subject, so storing it once is
// enough for every later change to show on the next page load.
func pictureURL(claim string) string {
	if strings.HasPrefix(claim, "https://") {
		return claim
	}
	return ""
}

// subjectPictureURL derives the picture URL Registre publishes for a subject,
// for callers that carry a subject but no picture claim. The issuer is the
// operator's configuration rather than a claim, so it is trusted as given.
func subjectPictureURL(issuer, subject string) string {
	if issuer == "" || subject == "" {
		return ""
	}
	return strings.TrimRight(issuer, "/") + "/avatars/" + url.PathEscape(subject)
}
