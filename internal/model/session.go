package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-deploy/dto/v2/body"
	"kthcloud-cli/internal/api"
	"kthcloud-cli/pkg/util"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type KeycloakSession struct {
	Token                 string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	ExpiresIn             int    `json:"expires_in"`
}

type AuthSession struct {
	SessionStart          *time.Time `json:"session_start,omitempty"`
	Token                 string     `json:"access_token"`
	RefreshToken          string     `json:"refresh_token"`
	RefreshTokenExpiresIn int        `json:"refresh_token_expires_in"`
	ExpiresIn             int        `json:"expires_in"`
}

type ApiKey struct {
	Key    string    `json:"key,omitempty"`
	Expiry time.Time `json:"expiry,omitempty"`
}

type Session struct {
	User        *body.UserRead `json:"user"`
	ApiKey      *ApiKey        `json:"api_key,omitempty"`
	AuthSession *AuthSession   `json:"auth_session"`
	Client      *api.Client
}

func NewAuthSession(token string, refreshToken string, expiresIn int, refreshTokenExpiresIn int) *AuthSession {
	now := time.Now()
	return &AuthSession{
		SessionStart:          &now,
		Token:                 token,
		RefreshToken:          refreshToken,
		ExpiresIn:             expiresIn,
		RefreshTokenExpiresIn: refreshTokenExpiresIn,
	}
}

func (k *KeycloakSession) ToAuthSession() *AuthSession {
	now := time.Now()
	return &AuthSession{
		SessionStart:          &now,
		Token:                 k.Token,
		RefreshToken:          k.RefreshToken,
		ExpiresIn:             k.ExpiresIn,
		RefreshTokenExpiresIn: k.RefreshTokenExpiresIn,
	}
}

func (as *AuthSession) SetTimeIfNotSet() {
	if as.SessionStart == nil {
		now := time.Now()
		as.SessionStart = &now
	}
}

func (as *AuthSession) IsExpired() bool {
	if as.SessionStart == nil {
		return true
	}
	expirationTime := as.SessionStart.Add(time.Duration(as.ExpiresIn) * time.Second)
	return time.Now().After(expirationTime)
}
func NewSession(auth *AuthSession) *Session {
	return &Session{
		AuthSession: auth,
		Client:      api.NewClient(viper.GetString("api-url"), auth.Token),
		User:        nil,
		ApiKey:      nil,
	}
}

func (s *Session) SetupClient() error {
	if s.ApiKey != nil && s.ApiKey.Key != "" && !time.Now().After(s.ApiKey.Expiry) {
		// api key is present and has not expired, lets use it
		s.Client = api.NewAPIClient(viper.GetString("api-url"), s.ApiKey.Key)
	} else if s.AuthSession.Token != "" && !s.AuthSession.IsExpired() {
		s.Client = api.NewClient(viper.GetString("api-url"), s.AuthSession.Token)
	} else {
		return errors.New("no authentication available")
	}
	return nil
}

func (s *Session) FetchUser() error {
	if s.Client == nil {
		s.SetupClient()
	}
	resp, err := s.Client.Req("/v2/users", "GET", nil)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return errors.New("non ok responsecode when fetching user")
	}

	users, err := util.ProcessResponseArr[body.UserRead](resp.String())
	if err != nil {
		return err
	}

	if len(users) != 1 {
		return errors.New("recieved more than one user")
	}

	s.User = &users[0]
	return nil
}

func (s *Session) Save(filename string) error {
	fmt.Println(filename)
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(s); err != nil {
		return fmt.Errorf("failed to encode session to JSON: %w", err)
	}

	return nil
}

func Load(filename string) (*Session, error) {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filename)
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var s Session
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&s); err != nil {
		return nil, err
	}

	return &s, nil
}
