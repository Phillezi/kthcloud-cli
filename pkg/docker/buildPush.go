package docker

import (
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

func buildPush(tag, dockerfile, context string) error {
	cmdParts := []string{
		"buildx", "build", "--platform=linux/amd64", "--tag=" + tag, "--push",
	}

	if dockerfile != "" {
		dockerfilePath := path.Join(context, dockerfile)
		logrus.Debugln("using dockerfile:", dockerfilePath)
		cmdParts = append(cmdParts, "--file="+dockerfilePath)
	}

	if context != "" {
		logrus.Debugln("using context:", context)
		cmdParts = append(cmdParts, context)
	} else {
		cmdParts = append(cmdParts, ".")
	}

	cmd := exec.Command("docker", cmdParts...)

	dockerctx, err := GetDockerContext(tag)
	if err != nil {
		return err
	}

	cmd.Env = append(os.Environ(), "DOCKER_CONTEXT="+dockerctx)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func RunBuildPushCommands(username, password, tag, dockerfile, context string) error {
	registry := strings.Split(tag, "/")[0]
	if err := LoginToRegistry(registry, username, password, tag); err != nil {
		return err
	}
	return buildPush(tag, dockerfile, context)
}
