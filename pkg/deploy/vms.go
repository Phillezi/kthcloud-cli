package deploy

import (
	"errors"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/resources"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/kthcloud/go-deploy/dto/v2/body"
	"github.com/spf13/viper"
)

func (c *Client) Vms() ([]body.VmRead, error) {
	if !c.HasValidSession() {
		return nil, errors.New("no active session, log in first")
	}
	if c.session.Resources != nil &&
		c.session.Resources.Vms != nil &&
		!c.session.Resources.Vms.IsExpired() {
		return c.session.Resources.Vms.Data, nil
	}

	req := c.client.R()

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

	if c.session.Resources == nil {
		c.session.Resources = &resources.Resources{}
	}
	c.session.Resources.Vms = &resources.CachedResource[[]body.VmRead]{
		Data:      vms,
		CachedAt:  time.Now(),
		ExpiresIn: viper.GetDuration("resource-cache-duration"),
	}

	return vms, nil
}

func (c *Client) DropVmsCache() {
	if c.session != nil && c.session.Resources != nil {
		c.session.Resources.DropVmsCache()
	}
}
