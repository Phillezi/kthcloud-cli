package builder

import (
	"errors"
	"fmt"

	"github.com/Phillezi/kthcloud-cli/internal/update"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/build"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/cicd"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/models/service"
)

// build a service
func Build(serviceName string, service *service.Service, yesToAll bool) error {
	if _, err := build.HasDockerCommands(); err != nil {
		return err
	}

	if service.Build == nil {
		return errors.New("no build config present")
	}

	deploymentID, err := getCICDDeploymentID(service.Build.Context, func(baseDir string) {
		if !yesToAll {
			yesBuild, _ := update.PromptYesNo("No existing deployment specified in " + baseDir + "/.kthcloud/DEPLOYMENT" + "\nDo you wish to create a cicd deployment for " + serviceName + "?")
			if !yesBuild {
				fmt.Println("Ok, wont build " + serviceName)
				return
			}
			yesAddWorkflow, _ := update.PromptYesNo("Do you want to add a Github workflow for CICD on " + serviceName + "?")
			cicd.Create(baseDir, yesAddWorkflow, serviceName)
		} else {
			cicd.Create(baseDir, true, serviceName)
		}
	})
	if err != nil {
		return err
	}

	conf, err := cicd.GetGHACIConf(deploymentID)
	if err != nil {
		return err
	}

	username, password, tag, err := cicd.ExtractSecrets(conf)
	if err != nil {
		return err
	}

	err = build.RunBuildPushCommands(username, password, tag, service.Build.Dockerfile, service.Build.Context)
	if err != nil {
		return err
	}

	return nil
}
