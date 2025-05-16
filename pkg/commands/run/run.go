package run

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/jobs"
	"github.com/Phillezi/kthcloud-cli/pkg/logs"
	"github.com/Phillezi/kthcloud-cli/pkg/response"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/kthcloud/go-deploy/dto/v2/body"
	"github.com/sirupsen/logrus"
)

func (c *Command) Run() error {
	if c.client == nil {
		return fmt.Errorf("client is nil")
	}

	if strings.TrimSpace(c.name) == "" {
		c.name = GenerateRandomName(3, 30)
	}

	if errs := c.validate(); len(errs) > 0 {
		for err := range errs {
			logrus.Error(err)
		}
		if len(errs) == 1 {
			return errs[0]
		}
		return fmt.Errorf("errors occured when validating")
	}

	resp, err := c.client.Create(ConvertToDeploymentCreate(
		c.name,
		c.envs,
		c.port,
		c.visibility,
		c.image,
		c.memory,
		c.cores,
		c.replicas,
	))
	if err != nil {
		return err
	}
	err = response.IsError(resp.String())
	if err != nil {
		logrus.Error(err)
		return err
	}
	job, err := util.ProcessResponse[body.DeploymentCreated](resp.String())
	if err != nil {
		logrus.Error(err)
		return err
	}

	removeDepl := func() {
		var found *body.DeploymentRead
		func() {
			for range 2 {
				depls, err := c.client.Deployments()
				if err != nil {
					logrus.Fatal(err)
				}

				for _, depl := range depls {
					if depl.Name == c.name {
						found = &depl
						return
					}
				}
				c.client.DropDeploymentsCache()
			}
		}()
		if found == nil {
			return
		}

		resp, err := c.client.Remove(found)
		if err != nil {
			logrus.Fatal(err)
		}
		err = response.IsError(resp.String())
		if err != nil {
			logrus.Fatal(err)
		}
		rmJob, err := util.ProcessResponse[body.DeploymentDeleted](resp.String())
		if err != nil {
			logrus.Fatal(err)
		}
		jobs.From(rmJob).Track(c.client, c.ctx, c.name, time.Millisecond*500, nil)
	}

	err = jobs.From(job).Track(c.client, c.ctx, c.name, time.Millisecond*500, removeDepl)
	if err != nil {
		return err
	}

	if c.remove {
		defer removeDepl()
	}

	if c.interactive || c.tty {
		token, _ := c.client.Token()
		apiKey, _ := c.client.ApiKey()
		conn := logs.CreateConn(c.name, job.ID, c.client.BaseURL(), token, apiKey)
		logCtx, cancelLogCtx := context.WithCancel(c.ctx)
		defer cancelLogCtx()
		logger := logs.New([]*logs.SSEConnection{conn}, logCtx)

		go logger.Start()
		<-c.ctx.Done()
		cancelLogCtx()
		logger.Stop()

	} else {
		fmt.Println(job.ID)
	}

	return nil
}
