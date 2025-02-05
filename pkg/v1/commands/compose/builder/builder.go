package builder

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/kthcloud/go-deploy/models/model"

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
	var err error

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
		deploymentID, err = GetCICDDeploymentID(service.Build.Context, onCicdNotConfigured)
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
			logrus.Debugln("cicd configured")
		}
	}

	var errGettingConf error = errors.New("tmp")
	var username string
	var password string
	var tag string
	var conf *model.GithubActionConfig
	maxRetries := 10

	for try := 0; errGettingConf != nil && try < maxRetries; try++ {
		conf, errGettingConf = cicd.GetGHACIConf(deploymentID)
		if errGettingConf != nil {
			logrus.Infof("could not get GHA config for cicd deployment, retrying in 500ms, retry [%d]\n", try)
			time.Sleep(500 * time.Millisecond)
			continue
			//return errGettingConf
		}

		username, password, tag, errGettingConf = cicd.ExtractSecrets(conf)
		if errGettingConf != nil {
			logrus.Infof("could not extract secrets from GHA config, retrying in 500ms, retry [%d]\n", try)
			if try < maxRetries {
				time.Sleep(500 * time.Millisecond)
			}
			//return errGettingConf
		}
	}
	if errGettingConf != nil {
		return errGettingConf
	}
	if username == "" || password == "" || tag == "" {
		return errors.New("username, password or tag is empty")
	}

	logrus.Debugln("starting build and push of", tag)
	err = build.RunBuildPushCommands(username, password, tag, service.Build.Dockerfile, service.Build.Context)
	if err != nil {
		return err
	}
	logrus.Debugln("done with build and push of", tag)

	return nil
}
