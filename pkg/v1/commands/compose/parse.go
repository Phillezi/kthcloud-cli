package compose

import (
	"encoding/json"
	"fmt"

	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/parser"
	"github.com/sirupsen/logrus"
)

func Parse(jsonOnly bool) {
	composeInstance, err := parser.GetCompose()
	if err != nil {
		logrus.Fatal(err)
	}

	if !jsonOnly {
		fmt.Println("Parsed Compose file:")
		fmt.Println(composeInstance.String() + "\n")

		fmt.Println("kthcloud deployments:")
	}
	deployments := composeInstance.ToDeployments()
	for _, deployment := range deployments {
		data, err := json.MarshalIndent(deployment, "", "  ")
		if err != nil {
			logrus.Fatalf("Error marshalling deployment to JSON: %v", err)
		}
		fmt.Println(string(data))
	}
}
