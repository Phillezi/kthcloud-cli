package down

import (
	"fmt"
)

func (c *Command) Run() error {
	if c.client == nil {
		return fmt.Errorf("client is nil")
	}
	if c.compose == nil {
		return fmt.Errorf("compose is nil")
	}

	if err := c.down(); err != nil {
		return err
	}

	return nil
}
