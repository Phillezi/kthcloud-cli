package compose

import (
	"go-deploy/dto/v2/body"
	"sync"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/jobs"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/parser"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/response"
	"github.com/briandowns/spinner"
	"github.com/sirupsen/logrus"
)

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
