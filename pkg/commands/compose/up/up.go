package up

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/builder"
	"github.com/Phillezi/kthcloud-cli/pkg/jobs"
	"github.com/Phillezi/kthcloud-cli/pkg/response"
	"github.com/Phillezi/kthcloud-cli/pkg/scheduler"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/briandowns/spinner"
	"github.com/go-resty/resty/v2"
	"github.com/kthcloud/go-deploy/dto/v2/body"
	"github.com/sirupsen/logrus"
)

func (c *Command) up() error {
	scheduleContext, cancelScheduler := context.WithCancel(c.ctx)
	sched := scheduler.NewSched(scheduleContext)

	go sched.Start()
	defer cancelScheduler()

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)

	s.Color("blue")
	s.Start()
	defer s.Stop()

	jobMap := make(map[string]*scheduler.Job)
	jobIDs := make(map[string]string, 1)

	for _, deployment := range c.compose.Deployments {
		job := scheduler.NewJob(func(ctx context.Context, cancelCallback func()) error {
			var resp *resty.Response
			var err error
			service := c.compose.Source.Services[deployment.Name]

			deplExistsWithSameName, deplWithSameNameID, deplWithSameNameHasSameImage := c.client.DeploymentExistsByNameWFilter(deployment.Name, func(depl body.DeploymentRead) bool {
				if util.NotNilOrEmpty(depl.Image) && util.NotNilOrEmpty(deployment.Image) {
					return *depl.Image == *deployment.Image
				}
				return util.IsEmptyOrNil(depl.Image) && util.IsEmptyOrNil(deployment.Image)
			})

			var mode string
			if !deplExistsWithSameName {
				resp, err = c.client.Create(deployment)
				mode = "create"
			} else if service.Build != nil {
				deplID, erro := builder.GetCICDDeploymentID(service.Build.Context, nil)
				if erro != nil {
					logrus.Error(erro)
					return erro
				}
				if !c.client.DeploymentExists(deplID) {
					logrus.Error("deployment for build of ", deployment.Name, " does not exist\n\tre-run with \"--build ", deployment.Name, "\" to ensure cicd deployment exists")
					return errors.New("cicd deployment doesnt exist")
				}
				updateDepl := util.DeploymentCreateToUpdate(deployment)
				resp, err = c.client.Update(&updateDepl, deplID)
				mode = "update"
			} else if deplWithSameNameHasSameImage {
				logrus.Debugln("deployment found with the specified service name and image, will update it to match the specification")
				updateDepl := util.DeploymentCreateToUpdate(deployment)
				resp, err = c.client.Update(&updateDepl, deplWithSameNameID)
				mode = "update"
			} else {
				return errors.New("service " + deployment.Name + " is not unique, (you already have a deployment with that name, and it doesnt have the same image)")
			}
			if err != nil {
				logrus.Error(err)
				return err
			}
			if resp.IsError() {
				logrus.Errorln("failed to  ", mode, " deployment ", deployment.Name, " status: ", resp.Status())

				logrus.Errorln("response body:", string(resp.Body()))

				cancelCallback()
				return err
			}
			err = response.IsError(resp.String())
			if err != nil {
				logrus.Error(err)
				return err
			}

			job, err := func() (any, error) {
				switch mode {
				case "create":
					return util.ProcessResponse[body.DeploymentCreated](resp.String())
				case "update":
					return util.ProcessResponse[body.DeploymentUpdated](resp.String())
				default:
					return nil, errors.ErrUnsupported
				}
			}()

			if err != nil {
				logrus.Error(err)
				return err
			}
			return jobs.From(job).Track(c.client, ctx, deployment.Name, time.Millisecond*500, cancelCallback)
		}, func() {
			logrus.Debugln("removing depl")
			var found *body.DeploymentRead
			func() {
				for range 2 {
					depls, err := c.client.Deployments()
					if err != nil {
						logrus.Fatal(err)
					}

					for _, depl := range depls {
						if depl.Name == deployment.Name {
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
			logrus.Debugln("tracking removal of depl")
			jobs.From(rmJob).Track(c.client, c.ctx, deployment.Name, time.Millisecond*500, nil)
		})
		jobMap[deployment.Name] = job

	}

	for depl, deps := range c.compose.Dependencies {
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
		return fmt.Errorf("err")
	}

	logrus.Debugln("alldone")
	c.creationDone = true
	s.Color("green")

	return nil
}
