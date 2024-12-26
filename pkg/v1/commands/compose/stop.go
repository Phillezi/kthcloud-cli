package compose

import (
	"context"
	"go-deploy/dto/v2/body"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/scheduler"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/jobs"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/parser"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/response"
	"github.com/briandowns/spinner"
	"github.com/sirupsen/logrus"
)

func Stop() {
	ctx, cancelStop := context.WithCancel(context.Background())
	done := make(chan bool)
	scheduleContext, cancelScheduler := context.WithCancel(ctx)
	sched := scheduler.NewSched(scheduleContext)

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)

	util.SetupSignalHandler(done, func() {
		sched.CancelJobsBlock()
		cancelStop()
		<-ctx.Done()
	})

	composeInstance, err := parser.GetCompose()
	if err != nil {
		logrus.Fatal(err)
	}

	c := client.Get()
	if !c.HasValidSession() {
		logrus.Fatal("no valid session, log in and try again")
	}

	go sched.Start()
	defer cancelScheduler()

	s.Color("blue")
	s.Start()
	defer s.Stop()

	jobIDs := make(map[string]string)

	c.DropDeploymentsCache()
	depls, err := c.Deployments()
	if err != nil {
		logrus.Fatal(err)
	}

	deploymentMap := make(map[string]*body.DeploymentRead)
	for _, depl := range depls {
		deploymentMap[depl.Name] = &depl
	}

	for name := range composeInstance.Services {
		if deployment, exists := deploymentMap[name]; exists {
			sjob := scheduler.NewJob(func(ctx context.Context, cancelCallback func()) error {
				disableDepl := &body.DeploymentUpdate{
					Replicas: util.IntPointer(0),
				}
				resp, err := c.Update(disableDepl, deployment.ID)
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
				return jobs.From(job).Track(ctx, deployment.Name, time.Millisecond*500, cancelCallback)
			}, func() {})
			jobIDs[name] = sched.AddJob(sjob)
		}
	}

	if err := jobs.MonitorJobStates(jobIDs, sched, s); err != nil {
		logrus.Debugln("erroccurred")
		s.Color("red")
	} else {
		logrus.Debugln("alldone")
		s.Color("green")
	}
}
