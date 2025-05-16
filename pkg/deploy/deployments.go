package deploy

import (
	"errors"
	"time"

	"github.com/kthcloud/go-deploy/dto/v2/body"

	"github.com/Phillezi/kthcloud-cli/pkg/resources"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/spf13/viper"
)

func (c *Client) Deployments() ([]body.DeploymentRead, error) {
	if !c.HasValidSession() {
		return nil, errors.New("no active session, log in first")
	}
	if c.session.Resources != nil &&
		c.session.Resources.Deployments != nil &&
		!c.session.Resources.Deployments.IsExpired() {
		return c.session.Resources.Deployments.Data, nil
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

	if c.session.Resources == nil {
		c.session.Resources = &resources.Resources{}
	}
	c.session.Resources.Deployments = &resources.CachedResource[[]body.DeploymentRead]{
		Data:      deployments,
		CachedAt:  time.Now(),
		ExpiresIn: viper.GetDuration("resource-cache-duration"),
	}

	return deployments, nil
}

func (c *Client) DropDeploymentsCache() {
	if c.session != nil && c.session.Resources != nil {
		c.session.Resources.DropDeploymentsCache()
	}
}
