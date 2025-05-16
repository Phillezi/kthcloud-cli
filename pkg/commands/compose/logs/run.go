package logs

import (
	"context"
	"fmt"

	"github.com/Phillezi/kthcloud-cli/pkg/logs"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/kthcloud/go-deploy/dto/v2/body"
	"github.com/sirupsen/logrus"
)

func (c *Command) Run() error {
	if c.client == nil {
		return fmt.Errorf("client is nil")
	}
	if c.compose == nil {
		return fmt.Errorf("compose is nil")
	}

	c.client.DropDeploymentsCache()
	depls, err := c.client.Deployments()
	if err != nil {
		logrus.Error(err)
		return err
	}

	deploymentMap := make(map[string]*body.DeploymentRead)
	for _, depl := range depls {
		deploymentMap[depl.Name] = &depl
	}

	var deploymentsToLog []*body.DeploymentRead
	if len(c.services) > 0 {
		for name := range c.compose.Services {
			if !util.Contains(c.services, name) {
				continue
			}
			if depl, exists := deploymentMap[name]; exists {
				deploymentsToLog = append(deploymentsToLog, depl)
			}
		}
	} else {
		for name := range c.compose.Services {
			if depl, exists := deploymentMap[name]; exists {
				deploymentsToLog = append(deploymentsToLog, depl)
			}
		}
	}

	if len(deploymentsToLog) == 0 {
		return fmt.Errorf("no instances to log")
	}

	token, err := c.client.Token()
	if err != nil {
		return err
	}

	key, _ := c.client.ApiKey()

	conns := logs.CreateConns(
		deploymentsToLog,
		c.client.BaseURL(),
		token,
		key,
	)

	logCtx, cancelLogCtx := context.WithCancel(c.ctx)
	defer cancelLogCtx()
	logger := logs.New(conns, logCtx)
	go logger.Start()
	<-c.ctx.Done()
	cancelLogCtx()
	logger.Stop()

	return nil
}
