package builder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/cicd"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/models/compose"
	"github.com/sirupsen/logrus"
)

func GetBuildsRequired(compose compose.Compose) map[string]bool {
	logrus.Traceln("builder.GetBuildsRequired")
	var wg sync.WaitGroup
	var mu sync.Mutex
	var erroccurred bool
	needsRebuildMap := make(map[string]bool)
	for name, service := range compose.Services {
		if service.Build != nil {
			id, err := GetCICDDeploymentID(service.Build.Context, func(baseDir string) {})
			if err != nil {
				mu.Lock()
				needsRebuildMap[name] = true
				mu.Unlock()
				continue
			}
			conf, err := cicd.GetGHACIConf(id)
			username, password, tag, err := cicd.ExtractSecrets(conf)
			wg.Add(1)
			go func() {
				HasDockerImage(username, password, tag, func() {
					mu.Lock()
					needsRebuildMap[name] = true
					mu.Unlock()
				}, func() {
					// TODO: handle better
					erroccurred = true
					logrus.Errorln("error could not check docker image status on registry for", name)
				})
				wg.Done()
			}()

		}
	}

	wg.Wait()
	if erroccurred {
		return nil
	}
	return needsRebuildMap
}

// tag contains the registry
func HasDockerImage(username, password, tag string, onNotExists func(), onError func()) (bool, error) {
	logrus.Traceln("builder.HasDockerImage")
	parts := strings.Split(tag, "/")
	if len(parts) < 2 {
		return false, fmt.Errorf("invalid tag format: %s", tag)
	}

	registry := parts[0]
	repoAndTag := strings.Join(parts[1:], "/")

	repoParts := strings.SplitN(repoAndTag, ":", 2)
	if len(repoParts) == 1 {
		repoParts = append(repoParts, "latest")
	}
	if len(repoParts) != 2 {
		logrus.Errorln("repoparts", repoParts)
		onError()
		return false, fmt.Errorf("tag must include an image and tag (e.g., repository/image:tag): %s", repoAndTag)
	}

	repository := repoParts[0]
	imageTag := repoParts[1]

	url := fmt.Sprintf("https://%s/v2/%s/manifests/%s", registry, repository, imageTag)
	logrus.Debugln("Check docker image url:", url)

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
		logrus.Debugln("registry has image")
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
