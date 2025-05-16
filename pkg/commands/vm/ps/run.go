package ps

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func (c *Command) Run() error {
	if c.client == nil {
		return fmt.Errorf("client is nil")
	}

	vms, err := c.client.Vms()
	if err != nil {
		logrus.Fatal("could not get vms:", err)
		return err
	}

	c.renderVmsTable(vms, c.all)

	return nil
}
