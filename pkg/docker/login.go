package docker

import (
	"os"
	"os/exec"
)

func LoginToRegistry(registry, username, password, tag string) error {
	cmd := exec.Command("docker", "login", registry, "-u", username, "-p", password)

	dockerctx, err := MakeDockerContext(tag)
	if err != nil {
		return err
	}

	cmd.Env = append(os.Environ(), "DOCKER_CONTEXT="+dockerctx)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
