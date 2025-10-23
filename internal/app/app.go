package app

import (
	"context"
	"fmt"

	"github.com/kthcloud/cli/internal/defaults"
	"github.com/kthcloud/cli/pkg/auth"
	"github.com/kthcloud/cli/pkg/deploy"
	"github.com/kthcloud/cli/pkg/keycloak"
	"github.com/kthcloud/cli/pkg/session"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type App struct {
	ctx context.Context

	deployAPIBaseURL string

	keycloakURL      string
	keycloakClientID string
	keycloakRealm    string

	loginServerPort string

	sessionKey         string
	sessionService     string
	sessionFallbackDir string

	oauth2Conf *oauth2.Config

	session session.Manager

	loginServer *auth.Server

	deploy deploy.ClientInterface

	l *zap.Logger
}

func New(ctx context.Context, opts ...Option) *App {
	a := App{
		ctx:              ctx,
		deployAPIBaseURL: defaults.DefaultDeployAPIBaseURL,

		keycloakURL:      defaults.DefaultKeycloakBaseURL,
		keycloakClientID: defaults.DefaultKeycloakClientID,
		keycloakRealm:    defaults.DefaultKeycloakRealm,

		loginServerPort: defaults.DefaultLoginServerPort,

		sessionKey:         defaults.DefaultKeystoreSessionKey,
		sessionService:     defaults.DefaultKeystoreServiceName,
		sessionFallbackDir: defaults.DefaultKeystoreFallbackDir,

		l: zap.NewNop(),
	}

	for _, opt := range opts {
		opt(&a)
	}

	if a.oauth2Conf == nil {
		a.oauth2Conf = keycloak.Config(a.keycloakClientID, a.keycloakURL, fmt.Sprintf("http://localhost:%s/callback", a.loginServerPort), a.keycloakRealm)
	}

	if a.session == nil {
		a.session = session.NewManager(
			session.WithContext(ctx),
			session.WithLogger(a.l.Named("session")),
			session.WithFallbackStoreDir(a.sessionFallbackDir),
			session.WithService(a.sessionService),
			session.WithOAuth2Config(a.oauth2Conf),
		)
	}

	if a.loginServer == nil {
		a.loginServer = auth.NewServer(
			auth.WithOAuth2Config(a.oauth2Conf),
			auth.WithLogger(a.l.Named("auth")),
		)
	}

	if a.deploy == nil {
		dc, err := deploy.NewClientWithResponses(a.deployAPIBaseURL, deploy.WithRequestEditorFn(a.session.AuthMiddleware))
		if err != nil {
			// TODO: handle me nicer
			panic(err)
		}
		a.deploy = dc

	}

	return &a
}

func (a *App) Deploy() deploy.ClientInterface {
	return a.deploy
}
