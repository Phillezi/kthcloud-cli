package auth

import (
	"context"
	"net/http"
	"path"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/config"
	"github.com/Phillezi/kthcloud-cli/pkg/defaults"
	"github.com/Phillezi/kthcloud-cli/pkg/session"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type Client struct {
	ctx context.Context

	keycloakBaseURL      string
	keycloakClientID     string
	keycloakClientSecret string
	keycloakRealm        string

	redirectHost string
	redirectPath string
	redirectURI  string

	sessionPath string

	requestTimeout time.Duration

	jar     http.CookieJar
	session *session.Session

	client *resty.Client
}

func New(opts ...ClientOpts) *Client {
	var opt ClientOpts
	if len(opts) > 0 {
		opt = opts[0]
	}

	c := &Client{
		ctx: context.Background(),

		keycloakBaseURL:      util.PtrOr(opt.KeycloakBaseURL),
		keycloakClientID:     util.PtrOr(opt.KeycloakClientID),
		keycloakClientSecret: util.PtrOr(opt.KeycloakClientSecret),
		keycloakRealm:        util.PtrOr(opt.KeycloakRealm),

		redirectHost: util.PtrOr(opt.RedirectHost, defaults.DefaultRedirectSchemeHostPort),
		redirectPath: util.PtrOr(opt.RedirectPath, defaults.DefaultRedirectBasePath),

		sessionPath: util.PtrOr(opt.SessionPath, path.Join(config.GetConfigPath(), "session.json")),

		requestTimeout: util.PtrOr(opt.RequestTimeout, defaults.DefaultRequestTimeout),

		session: util.Or(opt.Session),

		client: util.Or(opt.Client, resty.New()),
	}

	c.redirectURI, _ = c.constructRedirectURI()
	c.client.SetCookieJar(c.jar)

	if c.session == nil {
		sess, err := session.Load(c.sessionPath)
		if err != nil || sess.IsExpired() {
			// TODO: try to refresh token here later
			sess = nil
			logrus.Warn(err)
		}
		c.session = sess
	}

	return c
}

func (c *Client) WithContext(ctx context.Context) *Client {
	c.ctx = ctx
	return c
}

func (c *Client) WithRestClient(client *resty.Client) *Client {
	c.client = client
	return c
}

func (c *Client) WithSession(session *session.Session) *Client {
	c.session = session
	return c
}
