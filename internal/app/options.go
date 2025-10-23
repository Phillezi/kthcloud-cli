package app

import (
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type Option func(app *App)

func WithKeycloakOptions(clientID, baseURL, realm string) Option {
	return func(app *App) {
		app.keycloakClientID = clientID
		app.keycloakURL = baseURL
		app.keycloakRealm = realm
	}
}

func WithOAuth2Config(conf *oauth2.Config) Option {
	return func(a *App) {
		a.oauth2Conf = conf
	}
}

// Useful if you want to be able to use different users for
// different things
func WithSessionKey(sessionKey string) Option {
	return func(a *App) {
		a.sessionKey = sessionKey
	}
}

func WithLogger(l *zap.Logger) Option {
	return func(a *App) {
		a.l = l
	}
}
