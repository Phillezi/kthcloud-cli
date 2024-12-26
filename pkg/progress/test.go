package progress

import (
	"context"
	"go-deploy/dto/v2/body"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/scheduler"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/jobs"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/parser"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Test() {
	ctx, cancelStop := context.WithCancel(context.Background())
	done := make(chan bool)
	scheduleContext, cancelScheduler := context.WithCancel(ctx)
	sched := scheduler.NewSched(scheduleContext)

	util.SetupSignalHandler(done, func() {
		sched.CancelJobsBlock()
		cancelStop()
		<-ctx.Done()
	})

	composeInstance, err := parser.GetCompose()
	if err != nil {
		logrus.Fatal(err)
	}

	go sched.Start()
	defer cancelScheduler()

	jobMap := make(map[string]*scheduler.Job)
	jobIDs := make(map[string]string)

	deployments, dependencies := composeInstance.ToDeploymentsWDeps()
	for _, deployment := range deployments {
		jobMap[deployment.Name] = scheduler.NewJob(func(ctx context.Context, cancelCallback func()) error {
			return jobs.From(&body.DeploymentCreated{
				ID:    uuid.NewString(),
				JobID: uuid.NewString(),
			}).MockTrack(ctx, deployment.Name, time.Millisecond*500, cancelCallback)
		}, func() {})

	}

	for depl, deps := range dependencies {
		job := jobMap[depl]
		for _, dep := range deps {
			if depJob, ok := jobMap[dep]; !ok {
				logrus.Errorln("A dependency was found with no associated job")
			} else if depJob != nil {
				job.After(depJob)
			}
		}
		logrus.Debugf("%s has  %d dependencies\n", depl, len(job.Dependencies))

		jobIDs[depl] = sched.AddJob(job)
	}

	trackerT := New(sched)

	if err := trackerT.TrackJobs(); err != nil {
		logrus.Debugln("erroccurred")
	} else {
		logrus.Debugln("alldone")
	}
}
