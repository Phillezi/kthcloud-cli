package parse

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

func (c *Command) Run() error {
	if c.client == nil {
		return fmt.Errorf("client is nil")
	}
	if c.compose == nil {
		return fmt.Errorf("compose is nil")
	}

	if !c.json {
		fmt.Println("Parsed Compose file:")
		fmt.Println(c.compose.String() + "\n")

		fmt.Println("kthcloud deployments:")
	}
	deployments := c.compose.ToDeployments()
	for _, deployment := range deployments {
		data, err := json.MarshalIndent(deployment, "", "  ")
		if err != nil {
			logrus.Errorf("Error marshalling deployment to JSON: %v", err)
			return err
		}
		fmt.Println(string(data))
	}

	return nil
}
