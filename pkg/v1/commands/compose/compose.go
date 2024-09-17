package compose

import (
	"encoding/json"
	"fmt"
	"go-deploy/dto/v2/body"
	"sync"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/jobs"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/parser"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/response"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/storage"
	"github.com/briandowns/spinner"
	"github.com/sirupsen/logrus"
)

func Up(tryToCreateVolumes bool) {
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
		logrus.Infoln("Skipping trying to create volumes, no solution to auth on the storage managers proxy")
		logrus.Infoln("use --try-volumes to try anyway")
	}

	var wg sync.WaitGroup

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Color("blue")
	s.Start()
	defer s.Stop()

	deployments := composeInstance.ToDeployments()
	for _, deployment := range deployments {
		resp, err := c.Create(deployment)
		if err != nil {
			logrus.Fatal(err)
		}
		err = response.IsError(resp.String())
		if err != nil {
			logrus.Fatal(err)
		}
		job, err := util.ProcessResponse[body.DeploymentCreated](resp.String())
		if err != nil {
			logrus.Fatal(err)
		}
		jobs.TrackDeploymentCreationW(deployment.Name, job, &wg, s)
	}
	wg.Wait()
	s.Color("green")
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
	logrus.Fatal("not implemented yet...")
}
