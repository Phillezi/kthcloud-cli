package compose

import (
	"context"
	"encoding/json"
	"fmt"
	"go-deploy/dto/v2/body"
	"sync"
	"time"

	"github.com/Phillezi/kthcloud-cli/internal/update"
	"github.com/Phillezi/kthcloud-cli/pkg/scheduler"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/jobs"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/logs"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/parser"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/response"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/storage"
	"github.com/briandowns/spinner"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Up(detached, tryToCreateVolumes bool) {
	ctx, cancelUp := context.WithCancel(context.Background())
	var creationDone bool
	var cancelled bool
	done := make(chan bool)

	scheduleContext, cancelScheduler := context.WithCancel(ctx)
	sched := scheduler.NewSched(scheduleContext)

	util.SetupSignalHandler(done, func() {
		sched.CancelJobsBlock()
		cancelled = true
		cancelUp()
		<-ctx.Done()
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

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Color("blue")
	s.Start()
	defer s.Stop()

	jobIDs := make(map[string]string, 1)

	deployments := composeInstance.ToDeployments()
	for _, deployment := range deployments {

		job := scheduler.NewJob(func(ctx context.Context, callback func(cArg interface{})) error {
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
			return jobs.Track(ctx, deployment.Name, job, time.Millisecond*500, s, func() {
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

				delResp, err := c.Remove(found)
				if err != nil {
					logrus.Fatal(err)
				}
				err = response.IsError(resp.String())
				if err != nil {
					logrus.Fatal(err)
				}
				rmJob, err := util.ProcessResponse[body.DeploymentDeleted](delResp.String())
				if err != nil {
					logrus.Fatal(err)
				}
				logrus.Debugln("tracking removal of depl")
				jobs.TrackDel(deployment.Name, rmJob, time.Millisecond*500, s)
			})
		}, func(cArg interface{}) {}, nil)
		jobIDs[deployment.Name] = sched.AddJob(job)
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

func Parse() {
	composeInstance, err := parser.GetCompose()
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Println("Parsed Compose file:")
	fmt.Println(composeInstance.String() + "\n")

	fmt.Println("kthcloud deployments:")
	deployments := composeInstance.ToDeployments()
	for _, deployment := range deployments {
		data, err := json.MarshalIndent(deployment, "", "  ")
		if err != nil {
			logrus.Fatalf("Error marshalling deployment to JSON: %v", err)
		}
		fmt.Println(string(data))
	}
}

func Down() {
	composeInstance, err := parser.GetCompose()
	if err != nil {
		logrus.Fatal(err)
	}
	c := client.Get()
	if !c.HasValidSession() {
		logrus.Fatal("login")
	}

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
		if depl, exists := deploymentMap[name]; exists {
			c.Remove(depl)
		}
	}

	var wg sync.WaitGroup

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Color("blue")
	s.Start()
	defer s.Stop()

	for name := range composeInstance.Services {
		if deployment, exists := deploymentMap[name]; exists {
			resp, err := c.Remove(deployment)
			if err != nil {
				logrus.Fatal(err)
			}
			err = response.IsError(resp.String())
			if err != nil {
				logrus.Fatal(err)
			}
			job, err := util.ProcessResponse[body.DeploymentDeleted](resp.String())
			if err != nil {
				logrus.Errorln(resp.String())
				logrus.Fatal(err)
			}
			jobs.TrackDeploymentDeletionW(deployment.Name, job, &wg, s)
		}
	}
	wg.Wait()
	s.Color("green")
	s.Stop()
}

func Logs() {
	composeInstance, err := parser.GetCompose()
	if err != nil {
		logrus.Fatal(err)
	}
	c := client.Get()
	if !c.HasValidSession() {
		logrus.Fatal("login")
	}

	c.DropDeploymentsCache()
	depls, err := c.Deployments()
	if err != nil {
		logrus.Fatal(err)
	}

	deploymentMap := make(map[string]*body.DeploymentRead)
	for _, depl := range depls {
		deploymentMap[depl.Name] = &depl
	}

	var deploymentsToLog []*body.DeploymentRead
	for name := range composeInstance.Services {
		if depl, exists := deploymentMap[name]; exists {
			deploymentsToLog = append(deploymentsToLog, depl)
		}
	}

	if len(deploymentsToLog) == 0 {
		logrus.Fatal("no instances to log")
	}

	key := ""
	token := ""
	if c.Session != nil && c.Session.ApiKey != nil {
		key = c.Session.ApiKey.Key
	}
	if c.Session != nil {
		token = c.Session.Token.AccessToken
	}

	conns := logs.CreateConns(
		deploymentsToLog,
		viper.GetString("api-url"),
		token,
		key,
	)

	logger := logs.New(conns, context.Background())
	logger.Start()
}
