package client

import (
	"errors"
	"go-deploy/dto/v2/body"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/resources"
)

func (c *Client) User() (*body.UserRead, error) {
	if c.Session.Resources.User != nil && !c.Session.Resources.User.IsExpired() {
		return c.Session.Resources.User.Data, nil
	}

	req := c.client.R()
	req.SetAuthToken(c.Session.Token.AccessToken)

	resp, err := req.Get("/v2/users/" + *c.Session.ID)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, errors.New("request to get user returned with status: " + resp.Status())
	}

	user, err := util.ProcessResponse[body.UserRead](resp.String())
	if err != nil {
		return nil, err
	}

	if c.Session.Resources == nil {
		c.Session.Resources = &resources.Resources{}
	}
	c.Session.Resources.User = &resources.CachedResource[*body.UserRead]{
		Data:      user,
		CachedAt:  time.Now(),
		ExpiresIn: 1 * time.Hour,
	}

	return user, nil
}

func (c *Client) DropUserCache() {
	if c.Session.Resources != nil {
		c.Session.Resources.DropUserCache()
	}
}
