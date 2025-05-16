package docker

import (
	"fmt"
	"os/exec"
	"strings"
)

func MakeDockerContext(tag string) (string, error) {
	tagParts := strings.Split(tag, "/")
	ctxName := tagParts[len(tagParts)-1]

	// Check if the context already exists
	cmd := exec.Command("docker", "context", "inspect", ctxName)
	err := cmd.Run()
	if err == nil {
		// Context already exists
		return ctxName, nil
	}

	cmd = exec.Command(
		"docker", "context", "create", ctxName,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to create Docker context: %s\n%s", err, string(output))
	}

	return ctxName, nil
}

func GetDockerContext(tag string) (string, error) {
	tagParts := strings.Split(tag, "/")
	ctxName := tagParts[len(tagParts)-1]

	cmd := exec.Command("docker", "context", "inspect", ctxName)
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return ctxName, nil
}
