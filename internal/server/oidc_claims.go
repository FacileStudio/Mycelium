package server

import (
	"net/http"
	"strings"

	gooidc "github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

const uiLocalesParam = "ui_locales"

func uiLocalesOption(r *http.Request) []oauth2.AuthCodeOption {
	switch locale := r.URL.Query().Get(uiLocalesParam); locale {
	case "fr", "en", "es", "de":
		return []oauth2.AuthCodeOption{oauth2.SetAuthURLParam(uiLocalesParam, locale)}
	}
	return nil
}

func verifiedClaim(value any) bool {
	switch typed := value.(type) {
	case bool:
		return typed
	case string:
		return strings.EqualFold(typed, "true") || typed == "1"
	}
	return false
}

func stringClaim(value any) string {
	text, _ := value.(string)
	return text
}

func identityFromToken(idToken *gooidc.IDToken, issuer string) (oidcIdentity, bool) {
	var claims struct {
		Email             string `json:"email"`
		EmailVerified     any    `json:"email_verified"`
		Name              string `json:"name"`
		PreferredUsername string `json:"preferred_username"`
		Picture           string `json:"picture"`
		Locale            any    `json:"locale"`
	}
	if err := idToken.Claims(&claims); err != nil || claims.Email == "" || idToken.Subject == "" {
		return oidcIdentity{}, false
	}
	name := claims.Name
	if name == "" {
		name = claims.PreferredUsername
	}
	return oidcIdentity{
		Issuer: issuer, Subject: idToken.Subject, Email: claims.Email,
		EmailVerified: verifiedClaim(claims.EmailVerified), Name: name,
		Picture: claims.Picture, Locale: stringClaim(claims.Locale),
	}, true
}
