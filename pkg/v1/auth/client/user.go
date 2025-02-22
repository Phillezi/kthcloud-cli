package client

import (
	"errors"
	"time"

	"github.com/kthcloud/go-deploy/dto/v2/body"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/resources"
	"github.com/spf13/viper"
)

func (c *Client) User() (*body.UserRead, error) {
	if c.Session == nil {
		return nil, errors.New("no active session, log in first")
	}
	if c.Session.Resources != nil &&
		c.Session.Resources.User != nil &&
		!c.Session.Resources.User.IsExpired() {
		return c.Session.Resources.User.Data, nil
	}

	req := c.client.R()

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
		ExpiresIn: viper.GetDuration("resource-cache-duration"),
	}

	return user, nil
}

func (c *Client) DropUserCache() {
	if c.Session != nil && c.Session.Resources != nil {
		c.Session.Resources.DropUserCache()
	}
}
