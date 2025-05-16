package docker

import (
	"fmt"
	"os/exec"
)

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
