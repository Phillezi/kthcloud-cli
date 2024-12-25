package compose

import (
	"context"
	"errors"
	"go-deploy/dto/v2/body"
	"time"

	"github.com/Phillezi/kthcloud-cli/internal/update"
	"github.com/Phillezi/kthcloud-cli/pkg/scheduler"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/builder"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/jobs"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/parser"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/response"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/storage"
	"github.com/briandowns/spinner"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

func Up(detached, tryToCreateVolumes, buildAll, nonInteractive bool, servicesToBuild []string) {
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
			resp, err := update.PromptYesNo("Do you want to stop deployments")
			if err != nil {
				return
			}
			if resp {
				Stop()
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

	if buildAll {
		logrus.Debugln("buildAll is true")
		for n, s := range composeInstance.Services {
			if s.Build != nil {
				if err := builder.Build(n, s, nonInteractive); err != nil {
					logrus.Fatalln("Could not build service:", n, "Error:", err)
				}
				logrus.Debugln("build of", n, "is done!")
			}
		}
	} else if servicesToBuild != nil && len(servicesToBuild) > 0 {
		logrus.Debugln("services to build are specified")
		for n, s := range composeInstance.Services {
			if s.Build != nil {
				if util.Contains(servicesToBuild, n) {
					if err := builder.Build(n, s, nonInteractive); err != nil {
						logrus.Fatalln("Could not build service:", n, "Error:", err)
					}
					logrus.Debugln("build of", n, "is done!")
				}
			}
		}
	}

	buildsReq, err := builder.GetBuildsRequired(*composeInstance)
	if err != nil {
		logrus.Fatalln("Error getting builds required:", err)
	}
	for n, needsBuild := range buildsReq {
		if needsBuild {
			if err := builder.Build(n, composeInstance.Services[n], nonInteractive); err != nil {
				logrus.Fatalln("Could not build service:", n, "Error:", err)
			}
			logrus.Debugln("build of", n, "is done!")
		}
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
			var resp *resty.Response
			var err error
			service := composeInstance.Services[deployment.Name]

			deplExistsWithSameName, deplWithSameNameID, deplWithSameNameHasSameImage := client.Get().DeploymentExistsByNameWFilter(deployment.Name, func(depl body.DeploymentRead) bool {
				if util.NotNilOrEmpty(depl.Image) && util.NotNilOrEmpty(deployment.Image) {
					return *depl.Image == *deployment.Image
				}
				return util.IsEmptyOrNil(depl.Image) && util.IsEmptyOrNil(deployment.Image)
			})

			var mode string

			if !deplExistsWithSameName {
				resp, err = c.Create(deployment)
				mode = "create"
			} else if service.Build != nil {
				deplID, erro := builder.GetCICDDeploymentID(service.Build.Context, nil)
				if erro != nil {
					logrus.Error(erro)
					return erro
				}
				if !c.DeploymentExists(deplID) {
					logrus.Error("deployment for build of ", deployment.Name, " does not exist\n\tre-run with \"--build ", deployment.Name, "\" to ensure cicd deployment exists")
					return errors.New("cicd deployment doesnt exist")
				}
				updateDepl := util.DeploymentCreateToUpdate(deployment)
				resp, err = c.Update(&updateDepl, deplID)
				mode = "update"
			} else if deplWithSameNameHasSameImage {
				logrus.Debugln("deployment found with the specified service name and image, will update it to match the specification")
				updateDepl := util.DeploymentCreateToUpdate(deployment)
				resp, err = c.Update(&updateDepl, deplWithSameNameID)
				mode = "update"
			} else {
				return errors.New("service " + deployment.Name + " is not unique, (you already have a deployment with that name, and it doesnt have the same image)")
			}
			if err != nil {
				logrus.Error(err)
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
			return jobs.From(job).Track(ctx, deployment.Name, time.Millisecond*500, cancelCallback)
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
			jobs.From(rmJob).Track(ctx, deployment.Name, time.Millisecond*500, nil)
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
