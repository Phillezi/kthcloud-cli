package stop

import (
	"context"
	"fmt"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/jobs"
	"github.com/Phillezi/kthcloud-cli/pkg/response"
	"github.com/Phillezi/kthcloud-cli/pkg/scheduler"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/briandowns/spinner"
	"github.com/kthcloud/go-deploy/dto/v2/body"
	"github.com/sirupsen/logrus"
)

func (c *Command) stop() error {
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

	if len(c.services) > 0 {
		c.stopSpecifiedServices(jobIDs, deploymentMap, sched.AddJob)
	} else {
		c.stopAllServices(jobIDs, deploymentMap, sched.AddJob)
	}

	if err := jobs.MonitorJobStates(jobIDs, sched, s); err != nil {
		logrus.Debugln("erroccurred")
		s.Color("red")
		return fmt.Errorf("error stopping services")
	}
	logrus.Debugln("alldone")
	s.Color("green")

	return nil
}

func (c *Command) stopAllServices(
	jobIDs map[string]string,
	deploymentMap map[string]*body.DeploymentRead,
	addJob func(job *scheduler.Job,
	) string) error {

	for name := range c.compose.Source.Services {
		if deployment, exists := deploymentMap[name]; exists {
			sjob := scheduler.NewJob(func(ctx context.Context, cancelCallback func()) error {
				disableDepl := &body.DeploymentUpdate{
					Replicas: util.IntPointer(0),
				}
				resp, err := c.client.Update(disableDepl, deployment.ID)
				if err != nil {
					logrus.Error(err)
					return err
				}
				err = response.IsError(resp.String())
				if err != nil {
					logrus.Error(err)
					return err
				}
				job, err := util.ProcessResponse[body.DeploymentUpdated](resp.String())
				if err != nil {
					logrus.Error(err)
					return err
				}
				return jobs.From(job).Track(
					c.client,
					ctx,
					deployment.Name,
					time.Millisecond*500,
					cancelCallback,
				)
			}, func() {})
			jobIDs[name] = addJob(sjob)
		}
	}

	return nil
}

func (c *Command) stopSpecifiedServices(
	jobIDs map[string]string,
	deploymentMap map[string]*body.DeploymentRead,
	addJob func(job *scheduler.Job) string,
) error {
	if len(c.services) <= 0 {
		return fmt.Errorf("no specified services")
	}

	for name := range c.compose.Source.Services {
		if !util.Contains(c.services, name) {
			continue
		}
		if deployment, exists := deploymentMap[name]; exists {
			sjob := scheduler.NewJob(func(ctx context.Context, cancelCallback func()) error {
				disableDepl := &body.DeploymentUpdate{
					Replicas: util.IntPointer(0),
				}
				resp, err := c.client.Update(disableDepl, deployment.ID)
				if err != nil {
					logrus.Error(err)
					return err
				}
				err = response.IsError(resp.String())
				if err != nil {
					logrus.Error(err)
					return err
				}
				job, err := util.ProcessResponse[body.DeploymentUpdated](resp.String())
				if err != nil {
					logrus.Error(err)
					return err
				}
				return jobs.From(job).Track(
					c.client,
					ctx,
					deployment.Name,
					time.Millisecond*500,
					cancelCallback,
				)
			}, func() {})
			jobIDs[name] = addJob(sjob)
		}
	}

	return nil
}
