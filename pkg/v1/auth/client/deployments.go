package client

import (
	"errors"
	"go-deploy/dto/v2/body"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/resources"
)

func (c *Client) Deployments() ([]body.DeploymentRead, error) {
	if c.Session.Resources.User != nil && !c.Session.Resources.Deployments.IsExpired() {
		return c.Session.Resources.Deployments.Data, nil
	}

	req := c.client.R()
	req.SetAuthToken(c.Session.Token.AccessToken)

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
		ExpiresIn: 1 * time.Hour,
	}

	return deployments, nil
}

func (c *Client) DropDeploymentsCache() {
	if c.Session.Resources != nil {
		c.Session.Resources.DropVmsCache()
	}
}
