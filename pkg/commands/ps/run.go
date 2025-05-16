package ps

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func (c *Command) Run() error {
	if c.client == nil {
		return fmt.Errorf("client is nil")
	}

	c.client.DropDeploymentsCache()
	depls, err := c.client.Deployments()
	if err != nil {
		logrus.Fatal(err)
		return err
	}

	renderDeplsTable(depls, c.all)

	return nil
}
