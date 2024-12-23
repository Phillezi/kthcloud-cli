package builder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/cicd"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/models/compose"
)

func GetBuildsRequired(compose compose.Compose) map[string]bool {
	var wg sync.WaitGroup
	var mu sync.Mutex
	needsRebuildMap := make(map[string]bool)
	for name, service := range compose.Services {
		if service.Build != nil {
			id, err := getCICDDeploymentID(service.Build.Context, func(baseDir string) {})
			if err != nil {
				mu.Lock()
				needsRebuildMap[name] = true
				mu.Unlock()
				continue
			}
			conf, err := cicd.GetGHACIConf(id)
			username, password, tag, err := cicd.ExtractSecrets(conf)
			wg.Add(1)
			go HasDockerImage(username, password, tag, func() {
				mu.Lock()
				needsRebuildMap[name] = true
				mu.Unlock()
				wg.Done()
			})

		}
	}

	wg.Wait()
	return needsRebuildMap
}

// tag contains the registry
func HasDockerImage(username, password, tag string, onNotExists func()) (bool, error) {
	parts := strings.Split(tag, "/")
	if len(parts) < 2 {
		return false, fmt.Errorf("invalid tag format: %s", tag)
	}

	registry := parts[0]
	repoAndTag := strings.Join(parts[1:], "/")

	repoParts := strings.SplitN(repoAndTag, ":", 2)
	if len(repoParts) != 2 {
		return false, fmt.Errorf("tag must include an image and tag (e.g., repository/image:tag): %s", repoAndTag)
	}

	repository := repoParts[0]
	imageTag := repoParts[1]

	url := fmt.Sprintf("https://%s/v2/%s/manifests/%s", registry, repository, imageTag)

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to contact registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		if onNotExists != nil {
			onNotExists()
		}
		return false, nil
	} else {
		var errorMessage map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorMessage)
		return false, fmt.Errorf("unexpected response from registry: %s (%d) - %v", resp.Status, resp.StatusCode, errorMessage)
	}
}
