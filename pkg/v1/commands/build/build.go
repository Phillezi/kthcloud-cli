package build

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/cicd"
	"github.com/sirupsen/logrus"
)

func Build() {
	if _, err := HasDockerCommands(); err != nil {
		logrus.Fatal(err)
	}

	rootdir, _, err := cicd.GetGitRepoInfo()
	if err != nil {
		logrus.Fatal("Error when getting repo info:", err)
		return
	}
	repoConfDir := path.Join(rootdir, ".kthcloud")
	if !cicd.FileExists(repoConfDir, "DEPLOYMENT") {
		logrus.Fatalln("No deployment created for this repo\nRun the cicd command first to generate it")
	}

	id, err := cicd.ReadFile(repoConfDir, "DEPLOYMENT")
	if err != nil {
		logrus.Fatal("Error when trying to get deployment id", err)
		return
	}

	conf, err := cicd.GetGHACIConf(id)
	if err != nil {
		logrus.Fatal("Error when getting CI config:", err)
		return
	}

	username, password, tag, err := cicd.ExtractSecrets(conf)
	if err != nil {
		logrus.Fatal("Error when extracting secrets:", err)
		return
	}

	err = RunBuildPushCommands(username, password, tag, "", "")
	if err != nil {
		logrus.Fatal("Error running build and push command:", err)
		return
	}
}

func loginToRegistry(registry, username, password string) error {
	cmd := exec.Command("docker", "login", registry, "-u", username, "-p", password)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func buildPush(tag, dockerfile, context string) error {
	cmdParts := []string{
		"buildx", "build", "--platform=linux/amd64", "--tag=" + tag, "--push",
	}

	if dockerfile != "" {
		cmdParts = append(cmdParts, "--dockerfile="+dockerfile)
	}

	if context != "" {
		cmdParts = append(cmdParts, context)
	} else {
		cmdParts = append(cmdParts, ".")
	}

	cmd := exec.Command("docker", cmdParts...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func RunBuildPushCommands(username, password, tag, dockerfile, context string) error {
	registry := strings.Split(tag, "/")[0]
	if err := loginToRegistry(registry, username, password); err != nil {
		return err
	}
	return buildPush(tag, dockerfile, context)
}

func HasDockerCommands() (bool, error) {
	if _, err := exec.LookPath("docker"); err != nil {
		return false, fmt.Errorf("docker is not installed or not in PATH")
	}

	cmd := exec.Command("docker", "buildx", "version")
	err := cmd.Run()
	if err != nil {
		return false, fmt.Errorf("docker buildx is not installed or not in PATH")
	}

	return true, nil
}
