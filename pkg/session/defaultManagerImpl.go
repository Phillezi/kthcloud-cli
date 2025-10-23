package session

import (
	"context"
	"errors"
	"net/http"

	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type DefaultManager struct {
	ctx context.Context

	store     SessionStore
	refresher Refresher

	fallbackdir  string
	service      string
	oauth2Config *oauth2.Config

	l *zap.Logger
}

// NOTE: make sure you set the oauth2config
func NewManager(opts ...Option) *DefaultManager {
	m := &DefaultManager{
		ctx: context.Background(),

		l: zap.NewNop(),
	}

	for _, opt := range opts {
		opt(m)
	}

	if m.refresher == nil {
		if m.oauth2Config != nil {
			m.refresher = NewOAuth2Refresher(m.oauth2Config)
		} else {
			m.refresher = NopRefresherImpl{}
		}
	}

	if m.store == nil {
		store, err := NewSessionStore(m.service, m.fallbackdir)
		if err != nil {
			panic(err) // TODO: make this better
		}
		m.store = store
	}

	if _, isInsecure := m.store.(*FileStoreImpl); isInsecure {
		m.l.Warn("Using insecure (unencrypted) keystore implementation, your session keys will be stored in PLAINTEXT")
	}

	return m
}

func (m *DefaultManager) SaveSession(key string, session *Session) error {
	m.l.Debug("saving session", zap.String("key", key))
	return m.store.Set(key, session)
}

func (m *DefaultManager) GetSession(key string) (*Session, error) {
	m.l.Debug("getting session", zap.String("key", key))
	session, err := m.store.Get(key)
	if err != nil {
		m.l.Debug("failed to get session by key", zap.String("key", key), zap.Error(err))
		return nil, err
	}

	if session.IsExpired() && session.HasValidRefresh() && m.refresher != nil {
		m.l.Debug("session is expired but has valid refresh and a non nil refresher, will try to refresh using refreshtoken")
		newSession, err := m.refresher.Refresh(session.RefreshToken())
		if err != nil {
			m.l.Debug("failed refreshing", zap.Error(err))
			return nil, err
		}
		m.l.Debug("refresh was succesful", zap.Error(err))
		_ = m.store.Set(key, newSession)
		return newSession, nil
	}

	return session, nil
}

func (m *DefaultManager) DeleteSession(key string) error {
	m.l.Debug("deleting session", zap.String("key", key))
	return m.store.Delete(key)
}

func (m *DefaultManager) Clear() error {
	m.l.Debug("clearing all sessions")
	return m.store.Clear()
}

// AuthMiddleware injects a valid OAuth2 access token into the given request.
// It retrieves (and refreshes if needed) the session using the DefaultManager.
func (m *DefaultManager) AuthMiddleware(_ context.Context, req *http.Request) error {
	if req == nil {
		return ErrMiddlewareOnNilReq
	}

	m.l.Debug("auth middleware invoked", zap.String("url", req.URL.String()))

	session, err := m.GetSession("default")
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			m.l.Info("no existing session found, user needs to log in")
			return errors.Join(err, ErrLoginRequired)
		}
		m.l.Error("failed to retrieve session", zap.Error(err))
		return err
	}

	if session == nil || !session.IsValid() {
		m.l.Info("session is invalid or expired, user needs to reauthenticate")
		return ErrLoginRequired
	}

	authHeader := session.AuthHeader()
	if authHeader == "" {
		m.l.Warn("session missing valid token")
		return errors.New("empty or invalid access token")
	}

	req.Header.Set("Authorization", authHeader)
	m.l.Debug("attached auth header to request", zap.String("url", req.URL.String()), zap.String("type", session.Token.TokenType))

	return nil
}
