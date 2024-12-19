package compose

import (
	"context"
	"go-deploy/dto/v2/body"
	"time"

	"github.com/Phillezi/kthcloud-cli/internal/update"
	"github.com/Phillezi/kthcloud-cli/pkg/scheduler"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/jobs"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/parser"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/response"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/storage"
	"github.com/briandowns/spinner"
	"github.com/sirupsen/logrus"
)

func Up(detached, tryToCreateVolumes bool) {
	ctx, cancelUp := context.WithCancel(context.Background())
	var creationDone bool
	var cancelled bool
	done := make(chan bool)

	scheduleContext, cancelScheduler := context.WithCancel(ctx)
	sched := scheduler.NewSched(scheduleContext)

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)

	util.SetupSignalHandler(done, func() {
		sched.CancelJobsBlock()
		cancelled = true
		cancelUp()
		<-ctx.Done()
		s.Stop()
		if creationDone && !detached {
			resp, err := update.PromptYesNo("Do you want to terminate deployments")
			if err != nil {
				return
			}
			if resp {
				Down()
			}
		} else if !creationDone {
			logrus.Infoln("Cancelling creation of deployments")
		}
		s.Start()
	})

	if !detached {
		defer func() {
			if creationDone && !cancelled {
				logrus.Debug("Starting logger")
				go Logs()
				<-done
			}
		}()
	}

	composeInstance, err := parser.GetCompose()
	if err != nil {
		logrus.Fatal(err)
	}

	c := client.Get()
	if !c.HasValidSession() {
		logrus.Fatal("no valid session, log in and try again")
	}

	if tryToCreateVolumes {
		_, err = storage.CreateVolumes(c, composeInstance)
		if err != nil {
			logrus.Fatal(err)
		}
	} else {
		logrus.Infoln("Skipping volume creation from local structure")
		logrus.Infoln("If enabled it will \"steal\" cookies from your browser to authenticate")
		logrus.Infoln("use --try-volumes to try")
	}

	go sched.Start()
	defer cancelScheduler()

	s.Color("blue")
	s.Start()
	defer s.Stop()

	jobMap := make(map[string]*scheduler.Job)
	jobIDs := make(map[string]string, 1)

	deployments, dependencies := composeInstance.ToDeploymentsWDeps()
	for _, deployment := range deployments {

		job := scheduler.NewJob(func(ctx context.Context, cancelCallback func()) error {
			resp, err := c.Create(deployment)
			if err != nil {
				logrus.Error(err)
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
			return jobs.Track(ctx, deployment.Name, job, time.Millisecond*500, s, cancelCallback)
		}, func() {
			logrus.Debugln("removing depl")
			var found *body.DeploymentRead
			func() {
				for i := 0; i < 2; i++ {
					depls, err := c.Deployments()
					if err != nil {
						logrus.Fatal(err)
					}

					for _, depl := range depls {
						if depl.Name == deployment.Name {
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
			logrus.Debugln("tracking removal of depl")
			jobs.TrackDel(deployment.Name, rmJob, time.Millisecond*500, s)
		})
		jobMap[deployment.Name] = job

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

	if err := jobs.MonitorJobStates(jobIDs, sched, s); err != nil {
		logrus.Debugln("erroccurred")
		s.Color("red")
	} else {
		logrus.Debugln("alldone")
		creationDone = true
		s.Color("green")
	}
}
