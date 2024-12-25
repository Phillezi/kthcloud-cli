package client

import (
	"github.com/sirupsen/logrus"
)

func (c *Client) DeploymentExists(id string) bool {

	resp, err := c.client.R().
		Get(c.baseURL + "/v2/deployments/" + id)

	if err != nil {
		logrus.Errorln("error checking if deployment with id", id, "exists:", err)
		return false
	}

	if resp.StatusCode() == 404 {
		return false
	}

	return true
}
