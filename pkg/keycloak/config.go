package keycloak

import (
	"fmt"

	"golang.org/x/oauth2"
)

func Config(clientID, baseURL, redirectURL, realm string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:    clientID,
		Scopes:      []string{"openid", "profile", "email"},
		RedirectURL: redirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/realms/%s/protocol/openid-connect/auth", baseURL, realm),
			TokenURL: fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", baseURL, realm),
		},
	}
}
