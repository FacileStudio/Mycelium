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

// storedSubject is the issuer's subject a .users.json record belongs to, read
// the two ways findBySubject matches it: a stored subject, else the picture
// URL Registre publishes for that subject. It answers "" for neither.
func storedSubject(issuer string, user User) string {
	if user.Subject != "" {
		if user.Issuer != issuer {
			return ""
		}
		return user.Subject
	}
	prefix := subjectPictureURL(issuer, "x")
	escaped, ok := strings.CutPrefix(user.AvatarURL, strings.TrimSuffix(prefix, "x"))
	if !ok || escaped == "" {
		return ""
	}
	subject, err := url.PathUnescape(escaped)
	if err != nil || subjectPictureURL(issuer, subject) != user.AvatarURL {
		return ""
	}
	return subject
}
