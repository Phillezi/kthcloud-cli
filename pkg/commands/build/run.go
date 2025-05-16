package build

import (
	"fmt"
	"path"

	"github.com/Phillezi/kthcloud-cli/pkg/docker"
	"github.com/Phillezi/kthcloud-cli/pkg/file"
	"github.com/Phillezi/kthcloud-cli/pkg/git"
	"github.com/Phillezi/kthcloud-cli/pkg/github"
	"github.com/sirupsen/logrus"
)

func (c *Command) Run() error {
	if c.client == nil {
		return fmt.Errorf("client is nil")
	}

	if _, err := docker.HasDockerCommands(); err != nil {
		logrus.Fatal(err)
	}

	rootdir, _, err := git.GetGitRepoInfo()
	if err != nil {
		logrus.Error("Error when getting repo info:", err)
		return err
	}
	repoConfDir := path.Join(rootdir, ".kthcloud")
	if !file.Exists(repoConfDir, "DEPLOYMENT") {
		logrus.Fatalln("No deployment created for this repo\nRun the cicd command first to generate it")
	}

	id, err := file.Read(repoConfDir, "DEPLOYMENT")
	if err != nil {
		logrus.Error("Error when trying to get deployment id", err)
		return err
	}

	conf, err := c.client.CiConfig(id)
	if err != nil {
		logrus.Error("Error when getting CI config:", err)
		return err
	}

	model, err := github.ToModel(conf)
	if err != nil {
		logrus.Error("Error when converting CI config to GithubActions model:", err)
		return err
	}

	username, password, tag, err := github.ExtractSecrets(model)
	if err != nil {
		logrus.Error("Error when extracting secrets:", err)
		return err
	}

	err = docker.RunBuildPushCommands(username, password, tag, "", "")
	if err != nil {
		logrus.Error("Error running build and push command:", err)
		return err
	}
	return nil
}
