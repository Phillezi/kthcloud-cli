package filebrowser

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/defaults"
	"github.com/Phillezi/kthcloud-cli/pkg/session"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Client struct {
	ctx context.Context

	keycloakBaseURL string
	filebrowserURL  string

	requestTimeout time.Duration

	jar     *cookiejar.Jar
	session *session.Session

	client *http.Client

	token string
}

func New(opts ...ClientOpts) *Client {
	logrus.Traceln("pkg.filebrowser.client.go New")
	var opt ClientOpts
	if len(opts) > 0 {
		opt = opts[0]
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		logrus.Error("failed to create cookiejar")
	}

	c := &Client{
		ctx: util.PtrOr(opt.Context, context.Background()),

		keycloakBaseURL: util.PtrOr(opt.KeycloakBaseURL, defaults.DefaultKeycloakBaseURL),
		filebrowserURL:  util.PtrOr(opt.FilebrowserURL, defaults.DefaultStorageManagerProxy),

		requestTimeout: util.PtrOr(opt.RequestTimeout, defaults.DefaultRequestTimeout),

		session: util.Or(opt.Session),

		client: util.Or(opt.Client, &http.Client{
			Jar: jar,
		}),

		jar: jar,
	}

	c.client.Timeout = c.requestTimeout

	if c.session != nil {
		sess, err := session.Load(viper.GetString("session-path"))
		if err != nil || sess.IsExpired() {
			// TODO: try to refresh token here later
			sess = nil
		}
		c.session = sess
	}

	logrus.Debugln("Created filebrowser.Client with: filebrowserURL: " + c.filebrowserURL + " session: " + func() string {
		if c.session != nil {
			return "true"
		}
		return "false"
	}())

	return c
}

func (c *Client) WithContext(ctx context.Context) *Client {
	c.ctx = ctx
	return c
}

func (c *Client) WithFilebrowserURL(url string) *Client {
	c.filebrowserURL = url
	return c
}
