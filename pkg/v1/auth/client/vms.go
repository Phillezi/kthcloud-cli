package client

import (
	"errors"
	"go-deploy/dto/v2/body"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/resources"
)

func (c *Client) Vms() ([]body.VmRead, error) {
	if c.Session.Resources.User != nil && !c.Session.Resources.Vms.IsExpired() {
		return c.Session.Resources.Vms.Data, nil
	}

	req := c.client.R()
	req.SetAuthToken(c.Session.Token.AccessToken)

	resp, err := req.Get("/v2/vms")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, errors.New("request to get vms returned with status: " + resp.Status())
	}

	vms, err := util.ProcessResponseArr[body.VmRead](resp.String())
	if err != nil {
		return nil, err
	}

	if c.Session.Resources == nil {
		c.Session.Resources = &resources.Resources{}
	}
	c.Session.Resources.Vms = &resources.CachedResource[[]body.VmRead]{
		Data:      vms,
		CachedAt:  time.Now(),
		ExpiresIn: 1 * time.Hour,
	}

	return vms, nil
}

func (c *Client) DropVmsCache() {
	if c.Session.Resources != nil {
		c.Session.Resources.DropVmsCache()
	}
}
