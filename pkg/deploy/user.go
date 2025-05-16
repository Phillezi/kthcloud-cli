package deploy

import (
	"errors"
	"fmt"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/resources"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/kthcloud/go-deploy/dto/v2/body"
	"github.com/spf13/viper"
)

func (c *Client) User() (*body.UserRead, error) {
	if !c.HasValidSession() {
		return nil, errors.New("no active session, log in first")
	}
	if c.session.Resources != nil &&
		c.session.Resources.User != nil &&
		!c.session.Resources.User.IsExpired() {
		return c.session.Resources.User.Data, nil
	}
	if c.session.ID == nil {
		return nil, fmt.Errorf("session id is nil")
	}

	req := c.client.R()

	resp, err := req.Get("/v2/users/" + *c.session.ID)
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

	if c.session.Resources == nil {
		c.session.Resources = &resources.Resources{}
	}
	c.session.Resources.User = &resources.CachedResource[*body.UserRead]{
		Data:      user,
		CachedAt:  time.Now(),
		ExpiresIn: viper.GetDuration("resource-cache-duration"),
	}

	return user, nil
}

func (c *Client) DropUserCache() {
	if c.session != nil && c.session.Resources != nil {
		c.session.Resources.DropUserCache()
	}
}
