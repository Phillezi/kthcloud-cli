package session

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/token"
)

type Session struct {
	Token      token.JWTToken
	ExpiryTime time.Time
}

func New(token token.JWTToken) *Session {
	expiryTime := time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	return &Session{
		Token:      token,
		ExpiryTime: expiryTime,
	}
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiryTime)
}

func (s *Session) TimeUntilExpiry() time.Duration {
	return time.Until(s.ExpiryTime)
}

func (s *Session) RefreshAuthToken(newToken token.JWTToken) {
	s.Token = newToken
	s.ExpiryTime = time.Now().Add(time.Duration(newToken.ExpiresIn) * time.Second)
}

func (s *Session) Save(filepath string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write token to file: %w", err)
	}

	return nil
}

func Load(filepath string) (*Session, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read token file: %w", err)
	}

	var authToken Session
	err = json.Unmarshal(data, &authToken)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	return &authToken, nil
}
