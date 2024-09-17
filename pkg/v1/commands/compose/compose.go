package compose

import (
	"encoding/json"
	"fmt"

	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/parser"
	"github.com/sirupsen/logrus"
)

func Up() {
	composeInstance, err := parser.GetCompose()
	if err != nil {
		logrus.Fatal(err)
	}

	c := client.Get()
	if !c.HasValidSession() {
		logrus.Fatal("no valid session, log in and try again")
	}

	deployments := composeInstance.ToDeployments()
	for _, deployment := range deployments {
		resp, err := c.Create(deployment)
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Info(resp.String())
	}
}

func Parse() {
	composeInstance, err := parser.GetCompose()
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Println("Parsed Compose file:")
	fmt.Println(composeInstance.String() + "\n")

	fmt.Println("kthcloud deployments:")
	deployments := composeInstance.ToDeployments()
	for _, deployment := range deployments {
		data, err := json.MarshalIndent(deployment, "", "  ")
		if err != nil {
			logrus.Fatalf("Error marshalling deployment to JSON: %v", err)
		}
		fmt.Println(string(data))
	}
}

func Down() {
	logrus.Fatal("not implemented yet...")
}
