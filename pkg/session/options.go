package session

import (
	"context"

	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type Option func(m *DefaultManager)

func WithContext(ctx context.Context) Option {
	return func(m *DefaultManager) {
		m.ctx = ctx
	}
}

func WithOAuth2Config(conf *oauth2.Config) Option {
	return func(s *DefaultManager) {
		s.oauth2Config = conf
	}
}

func WithService(service string) Option {
	return func(m *DefaultManager) {
		m.service = service
	}
}

func WithFallbackStoreDir(fallbackStoreDir string) Option {
	return func(m *DefaultManager) {
		m.fallbackdir = fallbackStoreDir
	}
}

func WithSessionStore(sessionStore SessionStore) Option {
	return func(m *DefaultManager) {
		m.store = sessionStore
	}
}

func WithTokenRefresher(refresher Refresher) Option {
	return func(m *DefaultManager) {
		m.refresher = refresher
	}
}

func WithLogger(l *zap.Logger) Option {
	return func(m *DefaultManager) {
		m.l = l
	}
}
