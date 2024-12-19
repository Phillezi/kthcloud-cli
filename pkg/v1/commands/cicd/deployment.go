package cicd

import (
	"context"
	"go-deploy/dto/v2/body"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/jobs"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/response"
	"github.com/sirupsen/logrus"
)

func createDeployment(ctx context.Context, name string) (string, error) {
	c := client.Get()
	resp, err := c.Create(&body.DeploymentCreate{Name: name})
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	err = response.IsError(resp.String())
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	job, err := util.ProcessResponse[body.DeploymentCreated](resp.String())
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	err = jobs.Track(ctx, name, job, time.Millisecond*500, func() {
		var found *body.DeploymentRead
		func() {
			for i := 0; i < 2; i++ {
				depls, err := c.Deployments()
				if err != nil {
					logrus.Fatal(err)
				}

				for _, depl := range depls {
					if depl.Name == name {
						found = &depl
						return
					}
				}
				c.DropDeploymentsCache()
			}
		}()
		if found == nil {
			return
		}

		resp, err := c.Remove(found)
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

		jobs.TrackDel(name, rmJob, time.Millisecond*500)
	})
	return job.ID, err
}
