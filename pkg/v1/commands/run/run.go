package run

import (
	"context"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/jobs"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/response"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/models/options"
	"github.com/briandowns/spinner"
	"github.com/kthcloud/go-deploy/dto/v2/body"
	"github.com/sirupsen/logrus"
)

func Run(opts *options.DeploymentOptions) {
	ctx, cancelRun := context.WithCancel(context.Background())
	var creationDone bool
	var cancelled bool
	done := make(chan bool)

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)

	util.SetupSignalHandler(done, func() {
		cancelled = true
		cancelRun()
		<-ctx.Done()
		s.Stop()
		if creationDone && opts.RemoveOnExit {
			logrus.Infoln("not implemented yet")
		} else if !creationDone {
			logrus.Infoln("Cancelling creation of deployment")
		}
		s.Start()
	})

	if opts.InteractiveLogs {
		defer func() {
			if creationDone && !cancelled {
				logrus.Infoln("not implemented yet")
				<-done
			}
		}()
	}
	c := client.Get()
	if c == nil || !c.HasValidSession() {
		logrus.Fatal("no valid session, log in and try again")
	}

	depl, err := opts.ToDeploymentCreate()
	if err != nil {
		logrus.Fatal(err)
	}

	s.Start()

	resp, err := c.Create(depl)
	if err != nil {
		logrus.Fatal(err)
	}
	err = response.IsError(resp.String())
	if err != nil {
		logrus.Error(err)
	}
	job, err := util.ProcessResponse[body.DeploymentCreated](resp.String())
	if err != nil {
		logrus.Error(err)
	}
	err = jobs.Track(context.TODO(), opts.ContainerName, job, time.Millisecond*500, func() {
		var found *body.DeploymentRead
		func() {
			for i := 0; i < 2; i++ {
				depls, err := c.Deployments()
				if err != nil {
					logrus.Fatal(err)
				}

				for _, depl := range depls {
					if depl.Name == opts.ContainerName {
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

		jobs.TrackDel(opts.ContainerName, rmJob, time.Millisecond*500)
	})

	if err != nil {
		logrus.Debugln("error occurred", err)
		logrus.Error(err)
		s.Color("red")
	} else {
		logrus.Debugln("done")
		creationDone = true
		s.Color("green")
	}
}
