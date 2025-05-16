package cicd

import (
	"context"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/deploy"
	"github.com/Phillezi/kthcloud-cli/pkg/jobs"
	"github.com/Phillezi/kthcloud-cli/pkg/response"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/kthcloud/go-deploy/dto/v2/body"

	"github.com/sirupsen/logrus"
)

func CreateEmptyDeployment(client *deploy.Client, ctx context.Context, name string) (string, error) {
	resp, err := client.Create(&body.DeploymentCreate{Name: name})
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

	err = jobs.From(job).Track(client, ctx, name, time.Millisecond*500, func() {
		var found *body.DeploymentRead
		func() {
			for range 2 {
				depls, err := client.Deployments()
				if err != nil {
					logrus.Fatal(err)
				}

				for _, depl := range depls {
					if depl.Name == name {
						found = &depl
						return
					}
				}
				client.DropDeploymentsCache()
			}
		}()
		if found == nil {
			return
		}

		resp, err := client.Remove(found)
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
		jobs.From(rmJob).Track(client, ctx, name, time.Millisecond*500, nil)
	})
	return job.ID, err
}
