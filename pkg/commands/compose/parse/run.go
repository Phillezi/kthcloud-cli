package parse

import (
	"encoding/json"
	"fmt"
)

func (c *Command) Run() error {
	if c.client == nil {
		return fmt.Errorf("client is nil")
	}
	if c.compose == nil {
		return fmt.Errorf("compose is nil")
	}

	json, err := json.MarshalIndent(c.compose.Deployments, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(json))

	return nil
}
