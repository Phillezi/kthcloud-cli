package app

import (
	"strings"

	"github.com/kthcloud/cli/internal/defaults"
	"github.com/kthcloud/cli/pkg/session"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type Config struct {
	DeployAPIBaseURL string
	DeployAPIToken   string

	KeycloakClientID     string
	KeycloakClientSecret string
	KeycloakBaseURL      string
	KeycloakRealm        string

	LoginServerAddress string

	Oauth2Config *oauth2.Config

	SessionKey         string
	SessionService     string
	SessionFallbackDir string

	Session session.Auth

	Logger *zap.Logger
}

func (cfg *Config) Merge(config Config) {
	if config.DeployAPIBaseURL != "" {
		cfg.DeployAPIBaseURL = config.DeployAPIBaseURL
	}
	if config.DeployAPIToken != "" {
		cfg.DeployAPIToken = config.DeployAPIToken
	}

	if config.KeycloakClientID != "" {
		cfg.KeycloakClientID = config.KeycloakClientID
	}
	if config.KeycloakClientSecret != "" {
		cfg.KeycloakClientSecret = config.KeycloakClientSecret
	}
	if config.KeycloakBaseURL != "" {
		cfg.KeycloakBaseURL = config.KeycloakBaseURL
	}
	if config.KeycloakRealm != "" {
		cfg.KeycloakRealm = config.KeycloakRealm
	}

	if config.LoginServerAddress != "" {
		cfg.LoginServerAddress = config.LoginServerAddress
	}

	if config.Oauth2Config != nil {
		cfg.Oauth2Config = config.Oauth2Config
	}

	if config.SessionKey != "" {
		cfg.SessionKey = config.SessionKey
	}
	if config.SessionService != "" {
		cfg.SessionService = config.SessionService
	}
	if config.SessionFallbackDir != "" {
		cfg.SessionFallbackDir = config.SessionFallbackDir
	}

	if config.Session != nil {
		cfg.Session = config.Session
	}

	if config.Logger != nil {
		cfg.Logger = config.Logger
	}
}

func DefaultConfig() Config {
	return Config{
		DeployAPIBaseURL: defaults.DefaultDeployAPIBaseURL,

		KeycloakClientID:     defaults.DefaultKeycloakClientID,
		KeycloakBaseURL:      defaults.DefaultKeycloakBaseURL,
		KeycloakRealm:        defaults.DefaultKeycloakRealm,
		KeycloakClientSecret: defaults.DefaultKeycloakClientSecret,

		LoginServerAddress: defaults.DefaultLoginServerAddress,

		Oauth2Config: nil,

		SessionKey:         defaults.DefaultKeystoreSessionKey,
		SessionService:     defaults.DefaultKeystoreServiceName,
		SessionFallbackDir: defaults.DefaultKeystoreFallbackDir,

		Session: nil,

		Logger: zap.NewNop(),
	}
}

type Option func(cfg *Config)

func WithConfig(config Config) Option {
	return func(cfg *Config) {
		cfg.Merge(config)
	}
}

func WithKeycloakOptions(clientID, baseURL, realm string) Option {
	return func(app *Config) {
		app.KeycloakClientID = clientID
		app.KeycloakBaseURL = baseURL
		app.KeycloakRealm = realm
	}
}

func WithOAuth2Config(conf *oauth2.Config) Option {
	return func(a *Config) {
		a.Oauth2Config = conf
	}
}

// Useful if you want to be able to use different users for
// different things
func WithSessionKey(sessionKey string) Option {
	return func(a *Config) {
		a.SessionKey = sessionKey
	}
}

// Use api token based session
func WithAPITokenSession(token string) Option {
	// FIXME: usse DeployAPIToken
	if strings.TrimSpace(token) == "" {
		return func(_ *Config) {}
	}
	return func(app *Config) {
		app.Session = session.APITokenSession(token)
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(a *Config) {
		a.Logger = logger
	}
}
