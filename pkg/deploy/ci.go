package deploy

import (
	"errors"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/kthcloud/go-deploy/dto/v2/body"
)

func (c *Client) CiConfig(deploymentID string) (*body.CiConfig, error) {
	if !c.HasValidSession() {
		return nil, errors.New("no active session, log in first")
	}
	req := c.client.R()

	resp, err := req.Get("/v2/deployments/" + deploymentID + "/ciConfig")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, errors.New("request to deployment ci-config returned with status: " + resp.Status())
	}

	ciConf, err := util.ProcessResponse[body.CiConfig](resp.String())
	if err != nil {
		return nil, err
	}

	return ciConf, nil
}
