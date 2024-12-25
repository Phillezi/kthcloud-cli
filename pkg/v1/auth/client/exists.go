package client

import (
	"go-deploy/dto/v2/body"

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

func (c *Client) DeploymentExistsByName(name string) (bool, string) {
	for try := 0; try < 2; try++ {
		depls, err := c.Deployments()
		if err != nil {
			logrus.Errorln("error occurred when trying to get deployments", err)
			return false, ""
		}
		for _, depl := range depls {
			if depl.Name == name {
				return true, depl.ID
			}
		}
		logrus.Debugln("deployment with name:", name, ". Not found, dropping cache and trying again to ensure this")
		c.DropDeploymentsCache()
	}
	return false, ""
}

// returns existsWithSameName, idOfDeploymentWithTheName, matchesFilter
func (c *Client) DeploymentExistsByNameWFilter(name string, filter func(depl body.DeploymentRead) bool) (bool, string, bool) {
	for try := 0; try < 2; try++ {
		depls, err := c.Deployments()
		if err != nil {
			logrus.Errorln("error occurred when trying to get deployments", err)
			return false, "", false
		}
		for _, depl := range depls {
			if depl.Name == name && filter(depl) {
				return true, depl.ID, false
			} else if depl.Name == name {
				// depl with the same name exists but it doesnt match the filter
				return true, depl.ID, false
			}
		}
		logrus.Debugln("deployment with name:", name, ". Not found, dropping cache and trying again to ensure this")
		c.DropDeploymentsCache()
	}
	return false, "", false
}
