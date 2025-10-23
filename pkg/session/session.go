package session

import (
	"fmt"
	"time"

	"golang.org/x/oauth2"
)

type Session struct {
	Token         *oauth2.Token `json:"token"`
	RefreshExpiry time.Time     `json:"refresh_expiry"`
	CreatedAt     time.Time     `json:"created_at"`
}

// FromOAuth2Token wraps an oauth2.Token into a Session.
func FromOAuth2Token(tok *oauth2.Token) *Session {
	if tok == nil {
		return nil
	}

	refreshExpiry := tok.Expiry.Add(30 * 24 * time.Hour)
	if rtExp, ok := tok.Extra("refresh_expires_in").(float64); ok {
		refreshExpiry = time.Now().Add(time.Duration(rtExp) * time.Second)
	}

	return &Session{
		Token:         tok,
		RefreshExpiry: refreshExpiry,
		CreatedAt:     time.Now(),
	}
}

// IsExpired reports whether the access token has expired.
func (s *Session) IsExpired() bool {
	if s == nil || s.Token == nil {
		return true
	}
	return !s.Token.Valid()
}

// HasValidRefresh reports whether the refresh token is still valid.
func (s *Session) HasValidRefresh() bool {
	if s == nil || s.Token == nil {
		return false
	}
	return s.Token.RefreshToken != "" && time.Now().Before(s.RefreshExpiry)
}

// IsValid reports whether the session has a valid (non-expired) access token
// or at least a valid refresh token that can be used to obtain a new one.
func (s *Session) IsValid() bool {
	if s == nil {
		return false
	}
	return !s.IsExpired() || s.HasValidRefresh()
}

// AccessToken returns the raw access token string.
func (s *Session) AccessToken() string {
	if s == nil || s.Token == nil {
		return ""
	}
	return s.Token.AccessToken
}

// RefreshToken returns the raw refresh token string.
func (s *Session) RefreshToken() string {
	if s == nil || s.Token == nil {
		return ""
	}
	return s.Token.RefreshToken
}

// AuthHeader returns the properly formatted Authorization header value,
// e.g., "Bearer <token>", or "" if no token is available.
func (s *Session) AuthHeader() string {
	if s == nil || s.Token == nil || s.Token.AccessToken == "" {
		return ""
	}
	tokenType := s.Token.TokenType
	if tokenType == "" {
		tokenType = "Bearer"
	}
	return fmt.Sprintf("%s %s", tokenType, s.Token.AccessToken)
}
