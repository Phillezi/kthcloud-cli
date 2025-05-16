package up

import (
	"fmt"

	"github.com/Phillezi/kthcloud-cli/internal/interrupt"
	"github.com/Phillezi/kthcloud-cli/internal/update"
	"github.com/Phillezi/kthcloud-cli/pkg/commands/compose/logs"
	"github.com/Phillezi/kthcloud-cli/pkg/commands/compose/stop"
	"github.com/sirupsen/logrus"
)

func (c *Command) Run() error {
	if c.client == nil {
		return fmt.Errorf("client is nil")
	}
	if c.compose == nil {
		return fmt.Errorf("compose is nil")
	}

	if err := c.volumes(); err != nil {
		return err
	}

	if err := c.build(); err != nil {
		return err
	}

	if err := c.up(); err != nil {
		return err
	}

	if !c.detached && c.creationDone {
		interrupt.GetInstance().AddShutdownHook(func() {
			// todo: make this work
			resp, err := update.PromptYesNo("Do you want to stop deployments")
			if err != nil {
				return
			}
			if resp {
				stop.New(stop.CommandOpts{
					Client:  c.client,
					Compose: c.compose,
				}).WithContext(c.ctx).Run()
			}
		})
	}

	if !c.detached && c.creationDone && !c.cancelled {
		logrus.Debug("Starting logger")
		logs.New(logs.CommandOpts{
			Client:   c.client,
			Compose:  c.compose,
			Services: c.services,
		}).WithContext(c.ctx).Run()
	}

	return nil
}
