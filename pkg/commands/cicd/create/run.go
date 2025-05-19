package create

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/commands/cicd"
	"github.com/Phillezi/kthcloud-cli/pkg/file"
	"github.com/Phillezi/kthcloud-cli/pkg/git"
	"github.com/Phillezi/kthcloud-cli/pkg/github"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func (c *Command) Run() error {
	if c.client == nil {
		return fmt.Errorf("client is nil")
	}
	if c.deploymentName == "" {
		return fmt.Errorf("deployment name is required")
	}

	gitDir, upstreamURL, err := git.GetGitRepoInfoFrom(c.rootDir)
	if err != nil {
		gitDir = c.rootDir
	}
	repoConfDir := path.Join(c.rootDir, ".kthcloud")
	ghConfDir := path.Join(gitDir, ".github")
	wfConfDir := path.Join(ghConfDir, "workflows")

	var id string

	if !file.Exists(repoConfDir, "DEPLOYMENT") {
		id, err = cicd.CreateEmptyDeployment(c.client, context.Background(), c.deploymentName)
		if err != nil {
			logrus.Fatal("Error when creating empty deployment:", err)
			return err
		}
	} else {
		id, err = file.Read(repoConfDir, "DEPLOYMENT")
		if err != nil {
			logrus.Fatal("Error when trying to get deployment id", err)
			return err
		}
		if !c.client.DeploymentExists(id) {
			logrus.Debugln("file exists but contains ID of deployment that doesnt")
			id, err = cicd.CreateEmptyDeployment(c.client, context.Background(), c.deploymentName)
			if err != nil {
				logrus.Fatal("Error when creating empty deployment:", err)
				return err
			}
		}
	}

	if id != "" {
		err = file.Create(repoConfDir, "DEPLOYMENT", id)
		if err != nil {
			logrus.Fatal("Error when trying to save deployment id", err)
			return err
		}
	}

	if !c.createWorkFlow {
		return nil
	}

	conf, err := c.client.CiConfig(id)
	if err != nil {
		logrus.Fatal("Error when getting CI config:", err)
		return err
	}

	ghaConf, err := github.ToModel(conf)
	if err != nil {
		return err
	}

	username, password, tag, err := github.ExtractSecrets(ghaConf)
	if err != nil {
		logrus.Fatal("Error when trying remove secrets from workflow", err)
		return err
	}

	yamlConf, err := yaml.Marshal(ghaConf)
	if err != nil {
		logrus.Fatalf("error marshalling YAML: %v", err)
		return err
	}

	filename := "kthcloud.yml"
	if c.deploymentName != "" {
		filename = "kthcloud-" + c.deploymentName + "-cicd.yml"
	}
	for i := 0; file.Exists(wfConfDir, filename); i++ {
		filename = fmt.Sprintf("kthcloud-%s-cicd-%d.yml", c.deploymentName, i)
	}

	err = file.Create(wfConfDir, filename, string(yamlConf))
	if err != nil {
		logrus.Fatal("Error when trying save cicd config", err)
		return err
	}

	if strings.Contains(upstreamURL, "github.com") {
		github.PromptUserAddSecrets(upstreamURL, username, password, tag)
	}

	for range 5 {
		if file.Exists(repoConfDir, "DEPLOYMENT") {
			break
		}
		logrus.Debugln("waiting on fs...")
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}
