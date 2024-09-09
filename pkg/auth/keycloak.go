package auth

import (
	"encoding/json"
	"fmt"
	"kthcloud-cli/internal/model"
	"kthcloud-cli/pkg/util"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type KeycloakLoginResponse struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int    `json:"expires_in"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func KeycloakLogin(authURL, clientID, clientSecret, username, password string) (string, error) {
	client := resty.New()

	formData := map[string]string{
		"grant_type":    "password",
		"client_id":     clientID,
		"client_secret": clientSecret,
		"username":      username,
		"password":      password,
	}

	resp, err := client.R().
		SetFormData(formData).
		Post(authURL)

	if err != nil {
		return "", fmt.Errorf("failed to login to Keycloak: %w", err)
	}

	if resp.StatusCode() == 401 {
		return "", fmt.Errorf("authentication failed: %s", resp.String())
	}

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

func GetAuthSession(code, clientID, clientSecret, tokenURL, redirectURI string) (*model.AuthSession, error) {
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
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	KeycloakSession, err := util.ProcessResponse[model.KeycloakSession](string(resp.Body()))
	if err != nil {
		return nil, fmt.Errorf("failed to parse auth session response: %w", err)
	}

	return KeycloakSession.ToAuthSession(), nil
}
