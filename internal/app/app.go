package app

import (
	"context"
	"fmt"
	"net/http"

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

	loginServerAddress string

	sessionKey         string
	sessionService     string
	sessionFallbackDir string

	oauth2Conf *oauth2.Config

	session session.Auth

	loginServer *auth.Server

	deploy deploy.ClientWithResponsesInterface

	l *zap.Logger
}

func New(ctx context.Context, opts ...Option) *App {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	a := App{
		ctx: ctx,

		deployAPIBaseURL: cfg.DeployAPIBaseURL,

		keycloakURL:      cfg.KeycloakBaseURL,
		keycloakClientID: cfg.KeycloakClientID,
		keycloakRealm:    cfg.KeycloakRealm,

		loginServerAddress: cfg.LoginServerAddress,

		oauth2Conf: cfg.Oauth2Config,

		sessionKey:         cfg.SessionKey,
		sessionService:     cfg.SessionService,
		sessionFallbackDir: cfg.SessionFallbackDir,

		session: cfg.Session,

		l: cfg.Logger,
	}

	if a.oauth2Conf == nil {
		a.oauth2Conf = keycloak.Config(
			a.keycloakClientID,
			a.keycloakURL,
			fmt.Sprintf(
				"http://%s/callback",
				a.loginServerAddress,
			),
			a.keycloakRealm,
		)
	}

	if a.session == nil {
		a.session = session.NewManager(
			session.WithContext(ctx),
			session.WithLogger(a.l.Named("session")),
			session.WithFallbackStoreDir(a.sessionFallbackDir),
			session.WithService(a.sessionService),
			session.WithOAuth2Config(a.oauth2Conf),
			session.WithSessionKey(a.sessionKey),
		)
	}

	if a.loginServer == nil {
		a.loginServer = auth.NewServer(
			auth.WithOAuth2Config(a.oauth2Conf),
			auth.WithLogger(a.l.Named("auth")),
			auth.WithAddr(a.loginServerAddress),
		)
	}

	if a.deploy == nil {
		dc, err := deploy.NewClientWithResponses(
			a.deployAPIBaseURL,
			deploy.WithRequestEditorFn(a.session.AuthMiddleware),
		)
		if err != nil {
			// FIXME: handle me nicer
			panic(err)
		}
		a.deploy = dc

	}

	return &a
}

func (a *App) Deploy() deploy.ClientWithResponsesInterface {
	return a.deploy
}

func (a *App) SessionMiddleware() func(ctx context.Context, req *http.Request) error {
	return a.session.AuthMiddleware
}

func (a *App) Session() session.Auth {
	return a.session
}
