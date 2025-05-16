package deploy

import (
	"context"
	"sync"

	"github.com/Phillezi/kthcloud-cli/pkg/auth"
	"github.com/Phillezi/kthcloud-cli/pkg/defaults"
	"github.com/Phillezi/kthcloud-cli/pkg/filebrowser"
	"github.com/Phillezi/kthcloud-cli/pkg/session"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	once     sync.Once
	instance *Client
)

type Client struct {
	ctx context.Context

	baseURL string

	session *session.Session

	client *resty.Client

	// child clients
	authClient    *auth.Client
	storageClient *filebrowser.Client
}

func new(opts ...ClientOpts) *Client {
	var opt ClientOpts
	if len(opts) > 0 {
		opt = opts[0]
	}
	c := &Client{
		ctx: context.Background(),

		baseURL: util.PtrOr(opt.BaseURL, defaults.DefaultDeployBaseURL),

		session: opt.Session,
		client:  util.Or(opt.Client, resty.New()),
	}

	if c.session == nil {
		sess, err := session.Load(viper.GetString("session-path"))
		if err != nil || sess.IsExpired() {
			// TODO: try to refresh token here later
			sess = nil
			logrus.Warn(err)
		}
		c.session = sess
	}

	c.client.SetBaseURL(c.baseURL)
	if c.session != nil {
		c.client.SetAuthToken(c.session.Token.AccessToken)
	}

	return c
}

func GetInstance(opts ...ClientOpts) *Client {
	once.Do(func() {
		instance = new(opts...)
	})
	return instance
}

func (c *Client) WithContext(ctx context.Context) *Client {
	c.ctx = ctx
	return c
}

func (c *Client) WithAuthClient(authClientOpts ...auth.ClientOpts) *Client {
	if len(authClientOpts) <= 0 {
		authClientOpts = append(authClientOpts, auth.ClientOpts{})
	}
	authClientOpts[0].Client = c.client
	authClientOpts[0].Session = c.session
	c.authClient = auth.New(authClientOpts...).WithContext(c.ctx)
	return c
}

func (c *Client) WithStorageClient(filebrowserClientOpts ...filebrowser.ClientOpts) *Client {
	c.storageClient = filebrowser.New(filebrowserClientOpts...).WithContext(c.ctx)
	return c
}
