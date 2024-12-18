package compose

import (
	"context"
	"go-deploy/dto/v2/body"

	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/logs"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/parser"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Logs() {
	composeInstance, err := parser.GetCompose()
	if err != nil {
		logrus.Fatal(err)
	}
	c := client.Get()
	if !c.HasValidSession() {
		logrus.Fatal("login")
	}

	c.DropDeploymentsCache()
	depls, err := c.Deployments()
	if err != nil {
		logrus.Fatal(err)
	}

	deploymentMap := make(map[string]*body.DeploymentRead)
	for _, depl := range depls {
		deploymentMap[depl.Name] = &depl
	}

	var deploymentsToLog []*body.DeploymentRead
	for name := range composeInstance.Services {
		if depl, exists := deploymentMap[name]; exists {
			deploymentsToLog = append(deploymentsToLog, depl)
		}
	}

	if len(deploymentsToLog) == 0 {
		logrus.Fatal("no instances to log")
	}

	key := ""
	token := ""
	if c.Session != nil && c.Session.ApiKey != nil {
		key = c.Session.ApiKey.Key
	}
	if c.Session != nil {
		token = c.Session.Token.AccessToken
	}

	conns := logs.CreateConns(
		deploymentsToLog,
		viper.GetString("api-url"),
		token,
		key,
	)

	logger := logs.New(conns, context.Background())
	logger.Start()
}
