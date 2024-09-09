package auth

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

// KeycloakLoginResponse represents the response from Keycloak token endpoint with MFA
type KeycloakLoginResponse struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int    `json:"expires_in"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// KeycloakLogin performs login to Keycloak and retrieves the access token
func KeycloakLogin(authURL, clientID, clientSecret, username, password string) (string, error) {
	client := resty.New()

	// Set form data for Keycloak token request
	formData := map[string]string{
		"grant_type":    "password",
		"client_id":     clientID,
		"client_secret": clientSecret,
		"username":      username,
		"password":      password,
	}

	// Make the login request
	resp, err := client.R().
		SetFormData(formData).
		Post(authURL)

	if err != nil {
		return "", fmt.Errorf("failed to login to Keycloak: %w", err)
	}

	// Handle MFA challenge if needed
	if resp.StatusCode() == 401 {
		return "", fmt.Errorf("authentication failed: %s", resp.String())
	}

	// Parse the login response
	var loginResp KeycloakLoginResponse
	err = json.Unmarshal(resp.Body(), &loginResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse Keycloak response: %w", err)
	}

	if loginResp.Error != "" {
		return "", fmt.Errorf("error from Keycloak: %s - %s", loginResp.Error, loginResp.ErrorDescription)
	}

	log.Infof("Successfully logged in, access token: %s", loginResp.AccessToken)
	return loginResp.AccessToken, nil
}

// GetAccessToken exchanges authorization code for access token
func GetAccessToken(code, clientID, clientSecret, tokenURL, redirectURI string) (string, error) {
	client := resty.New()
	resp, err := client.R().
		SetFormData(map[string]string{
			"grant_type":    "authorization_code",
			"code":          code,
			"redirect_uri":  redirectURI,
			"client_id":     clientID,
			"client_secret": clientSecret,
		}).
		Post(tokenURL)

	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}

	var tokenResp map[string]interface{}
	err = json.Unmarshal(resp.Body(), &tokenResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	if token, ok := tokenResp["access_token"].(string); ok {
		return token, nil
	}
	return "", fmt.Errorf("failed to get access token from response")
}
