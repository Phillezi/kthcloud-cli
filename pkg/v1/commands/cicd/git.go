package cicd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetGitRepoInfo() (repoRoot string, upstreamURL string, err error) {

	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir, _ = os.Getwd()
	repoRootBytes, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("failed to get repo root: %v", err)
	}
	repoRoot = strings.TrimSpace(string(repoRootBytes))

	cmd = exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = repoRoot
	upstreamURLBytes, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("failed to get upstream URL: %v", err)
	}
	upstreamURL = strings.TrimSpace(string(upstreamURLBytes))

	return repoRoot, upstreamURL, nil
}
