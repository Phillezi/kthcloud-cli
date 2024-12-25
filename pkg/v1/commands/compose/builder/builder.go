package builder

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/Phillezi/kthcloud-cli/internal/update"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/build"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/cicd"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/models/service"
	"github.com/sirupsen/logrus"
)

// build a service
func Build(serviceName string, service *service.Service, yesToAll bool) error {
	if _, err := build.HasDockerCommands(); err != nil {
		return err
	}

	if service.Build == nil {
		return errors.New("no build config present")
	}

	exists := false
	deploymentID := ""

	for !exists {
		onCicdNotConfigured := func(baseDir string) {
			if !cicd.FileExists(baseDir, ".") {
				logrus.Warnln("Build context specifies a directory that doesnt exist:", baseDir)
			}
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
		}
		deploymentID, err := GetCICDDeploymentID(service.Build.Context, onCicdNotConfigured)
		if err != nil {
			return err
		}

		if exists = client.Get().DeploymentExists(deploymentID); !exists {
			logrus.Errorln("cicd deployment specified in .kthcloud dir doesnt exist")
			contextPath := service.Build.Context
			if contextPath == "" {
				contextPath = "."
			}
			wd, _ := os.Getwd()
			fullpath := path.Join(wd, contextPath)
			onCicdNotConfigured(fullpath)
		}
	}

	conf, err := cicd.GetGHACIConf(deploymentID)
	if err != nil {
		return err
	}

	username, password, tag, err := cicd.ExtractSecrets(conf)
	if err != nil {
		return err
	}

	logrus.Debugln("starting build and push of", tag)
	err = build.RunBuildPushCommands(username, password, tag, service.Build.Dockerfile, service.Build.Context)
	if err != nil {
		return err
	}
	logrus.Debugln("done with build and push of", tag)

	return nil
}
