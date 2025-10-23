package session

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"
)

type OAuth2RefreshImpl struct {
	Config *oauth2.Config
}

func NewOAuth2Refresher(conf *oauth2.Config) *OAuth2RefreshImpl {
	return &OAuth2RefreshImpl{Config: conf}
}

func (r *OAuth2RefreshImpl) Refresh(refreshToken string) (*Session, error) {
	if r.Config == nil {
		return nil, fmt.Errorf("oauth2 config is nil")
	}
	if refreshToken == "" {
		return nil, fmt.Errorf("refresh token is empty")
	}

	ctx := context.Background()
	src := r.Config.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken})
	newTok, err := src.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	s := FromOAuth2Token(newTok)

	// Preserve the original refresh token if the provider did not include a new one
	if s.Token.RefreshToken == "" {
		s.Token.RefreshToken = refreshToken
	}

	// Ensure a default refresh expiry if not set by provider
	if s.RefreshExpiry.IsZero() {
		s.RefreshExpiry = time.Now().Add(30 * 24 * time.Hour)
	}

	return s, nil
}
