package auth

import (
	"fmt"
	"net/url"
	"strings"

	"math/rand/v2"
)

func generateRandomState() string {
	const charset = "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx"
	randomize := func(c rune) rune {
		r := rand.IntN(16)
		if c == 'x' {
			return rune(fmt.Sprintf("%x", r)[0])
		}
		return rune(fmt.Sprintf("%x", (r&0x3)|0x8)[0])
	}
	return strings.Map(randomize, charset)
}

func (c *Client) constructRedirectURI() (string, error) {
	redirectURI := fmt.Sprintf("%s%s", c.redirectHost, c.redirectPath)
	_, err := url.Parse(redirectURI)
	if err != nil {
		return "", err
	}

	return redirectURI, nil
}

func (c *Client) constructKeycloakURL() (string, error) {
	state := generateRandomState()
	nonce := generateRandomState()

	keycloakURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/auth?client_id=%s&redirect_uri=%s&response_type=code&response_mode=query&scope=openid&nonce=%s&state=%s",
		c.keycloakBaseURL, c.keycloakRealm, c.keycloakClientID, url.QueryEscape(c.redirectURI), nonce, state)
	_, err := url.Parse(keycloakURL)
	if err != nil {
		return "", err
	}

	return keycloakURL, nil
}

func (c *Client) constructOAuthTokenURL() (string, error) {
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.keycloakBaseURL, c.keycloakRealm)
	_, err := url.Parse(tokenURL)
	if err != nil {
		return "", err
	}

	return tokenURL, nil
}
