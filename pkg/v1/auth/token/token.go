package token

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type JWTToken struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	IdToken          string `json:"id_token"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type AuthToken struct {
	Token      JWTToken
	ExpiryTime time.Time
}

func NewAuthToken(token JWTToken) *AuthToken {
	expiryTime := time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	return &AuthToken{
		Token:      token,
		ExpiryTime: expiryTime,
	}
}

func (a *AuthToken) IsExpired() bool {
	return time.Now().After(a.ExpiryTime)
}

func (a *AuthToken) TimeUntilExpiry() time.Duration {
	return time.Until(a.ExpiryTime)
}

func (a *AuthToken) RefreshAuthToken(newToken JWTToken) {
	a.Token = newToken
	a.ExpiryTime = time.Now().Add(time.Duration(newToken.ExpiresIn) * time.Second)
}

func (a *AuthToken) Save(filepath string) error {
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write token to file: %w", err)
	}

	return nil
}

func LoadAuthToken(filepath string) (*AuthToken, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read token file: %w", err)
	}

	var authToken AuthToken
	err = json.Unmarshal(data, &authToken)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	return &authToken, nil
}
