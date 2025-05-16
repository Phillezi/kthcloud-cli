package deploy

import (
	"github.com/Phillezi/kthcloud-cli/pkg/auth"
	"github.com/Phillezi/kthcloud-cli/pkg/filebrowser"
	"github.com/Phillezi/kthcloud-cli/pkg/session"
)

func (c *Client) Auth() *auth.Client {
	if c.authClient == nil {
		c.authClient = auth.New(auth.ClientOpts{
			Client:  c.client,
			Session: c.session,
		}).WithContext(c.ctx)
	}
	return c.authClient
}

func (c *Client) Login() (sess *session.Session, err error) {
	s, err := c.Auth().Login()
	if err != nil {
		return s, err
	}
	c.session = s
	c.client.SetAuthToken(c.session.Token.AccessToken)
	return s, err
}

func (c *Client) Storage() *filebrowser.Client {
	if c.storageClient == nil {
		c.storageClient = filebrowser.New().WithContext(c.ctx)
	}
	return c.storageClient
}
