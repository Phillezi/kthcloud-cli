package down

import (
	"context"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/jobs"
	"github.com/Phillezi/kthcloud-cli/pkg/response"
	"github.com/Phillezi/kthcloud-cli/pkg/scheduler"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/briandowns/spinner"
	"github.com/kthcloud/go-deploy/dto/v2/body"
	"github.com/sirupsen/logrus"
)

func (c *Command) down() error {

	scheduleContext, cancelScheduler := context.WithCancel(c.ctx)
	sched := scheduler.NewSched(scheduleContext)

	go sched.Start()
	defer cancelScheduler()

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)

	s.Color("blue")
	s.Start()
	defer s.Stop()

	jobIDs := make(map[string]string)

	c.client.DropDeploymentsCache()
	depls, err := c.client.Deployments()
	if err != nil {
		logrus.Fatal(err)
	}

	deploymentMap := make(map[string]*body.DeploymentRead)
	for _, depl := range depls {
		deploymentMap[depl.Name] = &depl
	}

	for name := range c.compose.Source.Services {
		if deployment, exists := deploymentMap[name]; exists {
			if !c.all && deployment.Image == nil {
				logrus.Infoln("Skipping deletion of deployment:", deployment.Name, ". Since it is a custom deployment (cicd)\n\nUse:\n\t--all\n\nTo remove CICD deployments too")
				continue
			}

			jobIDs[name] = sched.AddJob(scheduler.NewJob(func(ctx context.Context, cancelCallback func()) error {
				resp, err := c.client.Remove(deployment)
				if err != nil {
					logrus.Error(err)
					return err
				}
				err = response.IsError(resp.String())
				if err != nil {
					logrus.Error(err)
					return err
				}
				job, err := util.ProcessResponse[body.DeploymentDeleted](resp.String())
				if err != nil {
					logrus.Error(err)
					return err
				}
				return jobs.From(job).Track(
					c.client,
					ctx, deployment.Name,
					time.Millisecond*500,
					cancelCallback,
				)
			}, func() {}))
		}
	}
	if err := jobs.MonitorJobStates(jobIDs, sched, s); err != nil {
		logrus.Debugln("erroccurred")
		s.Color("red")
	} else {
		logrus.Debugln("alldone")
		s.Color("green")
	}

	return nil
}
