package session

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/resources"
	"github.com/kthcloud/go-deploy/dto/v2/body"
	"github.com/sirupsen/logrus"
)

type Session struct {
	Token      JWTToken             `json:"token" yaml:"token"`
	ExpiryTime time.Time            `json:"expiry_time" yaml:"expiry_time"`
	Resources  *resources.Resources `json:"resources,omitempty" yaml:"resources,omitempty"`
	ID         *string              `json:"id,omitempty" yaml:"id,omitempty"`
	ApiKey     *body.ApiKeyCreated  `json:"api_key,omitempty" yaml:"api_key,omitempty"`
	Cookies    []*http.Cookie       `json:"cookies,omitempty" yaml:"cookies,omitempty"`
}

func New(token JWTToken) *Session {
	expiryTime := time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	var id *string
	data, err := token.Decode()
	if err != nil {
		logrus.Warn("Could not parse the JWT token")
	}
	sub, ok := data["sub"].(string)
	if !ok {
		logrus.Warn("JWT 'sub' claim is not a string")
	}
	id = &sub

	return &Session{
		Token:      token,
		ExpiryTime: expiryTime,
		Resources:  nil,
		ID:         id,
	}
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiryTime)
}

func (s *Session) TimeUntilExpiry() time.Duration {
	return time.Until(s.ExpiryTime)
}

func (s *Session) RefreshAuthToken(newToken JWTToken) {
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

	if authToken.ID == nil {
		var id *string
		data, err := authToken.Token.Decode()
		if err != nil {
			logrus.Warn("Could not parse the JWT token")
		}
		sub, ok := data["sub"].(string)
		if !ok {
			logrus.Warn("JWT 'sub' claim is not a string")
		}
		id = &sub
		authToken.ID = id
	}

	return &authToken, nil
}
