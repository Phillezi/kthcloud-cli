package client

import (
	"errors"
	"go-deploy/dto/v2/body"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/resources"
	"github.com/spf13/viper"
)

func (c *Client) Deployments() ([]body.DeploymentRead, error) {
	if c.Session.Resources != nil && c.Session.Resources.Deployments != nil && !c.Session.Resources.Deployments.IsExpired() {
		return c.Session.Resources.Deployments.Data, nil
	}

	req := c.client.R()

	resp, err := req.Get("/v2/deployments")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, errors.New("request to get deployments returned with status: " + resp.Status())
	}

	deployments, err := util.ProcessResponseArr[body.DeploymentRead](resp.String())
	if err != nil {
		return nil, err
	}

	if c.Session.Resources == nil {
		c.Session.Resources = &resources.Resources{}
	}
	c.Session.Resources.Deployments = &resources.CachedResource[[]body.DeploymentRead]{
		Data:      deployments,
		CachedAt:  time.Now(),
		ExpiresIn: viper.GetDuration("resource-cache-duration"),
	}

	return deployments, nil
}

func (c *Client) DropDeploymentsCache() {
	if c.Session.Resources != nil {
		c.Session.Resources.DropVmsCache()
	}
}