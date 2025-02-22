package cicd

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/kthcloud/go-deploy/dto/v2/body"
	"github.com/kthcloud/go-deploy/models/model"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/response"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func GetGHACIConf(id string) (*model.GithubActionConfig, error) {
	if id == "" {
		return nil, errors.New("id cant be empty")
	}
	c := client.Get()

	r := c.Client().R()
	resp, err := r.Get("/v2/deployments/" + id + "/ciConfig")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() >= 300 {
		err := fmt.Errorf("received unexpected status code %d: %s", resp.StatusCode(), resp.Status())
		logrus.WithFields(logrus.Fields{
			"status_code": resp.StatusCode(),
			"response":    resp.String(),
			"called with": fmt.Sprintf("id:\"%s\"", id),
		}).Error("get ciConfig failed")
		return nil, err
	}

	err = response.IsError(resp.String())
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	ciConf, err := util.ProcessResponse[body.CiConfig](resp.String())
	if err != nil {
		return nil, fmt.Errorf("error processing ci config for deployment %s: %w", id, err)
	}

	logrus.Debugln("Got GHA config:", ciConf)

	var config model.GithubActionConfig

	err = yaml.Unmarshal([]byte(ciConf.Config), &config)
	if err != nil {
		log.Fatalf("error unmarshalling YAML: %v", err)
		return nil, err
	}

	return &config, nil
}

func ExtractSecrets(config *model.GithubActionConfig) (username, password, tag string, err error) {
	if config != nil && config.Jobs.Docker.Steps != nil {
		for i, step := range config.Jobs.Docker.Steps {
			if step.With.Password != "" {
				if username == "" {
					username = step.With.Username
				}
				if password == "" {
					password = step.With.Password
				}

				config.Jobs.Docker.Steps[i].With.Password = "${{ secrets.DOCKER_PASSWORD }}"
				config.Jobs.Docker.Steps[i].With.Username = "${{ secrets.DOCKER_USERNAME }}"
			}
			if step.With.Tags != "" {
				if tag == "" {
					tag = step.With.Tags
				}
				config.Jobs.Docker.Steps[i].With.Tags = "${{ secrets.DOCKER_TAG }}"
			}
		}
	} else {
		return "", "", "", errors.New("invalid config")
	}
	return username, password, tag, nil
}

func promptUserAddSecrets(upstreamURL, username, password, tag string) {
	if strings.HasPrefix(upstreamURL, "git@github.com:") {
		upstreamURL = strings.Replace(upstreamURL, "git@github.com:", "https://github.com/", 1)
	}
	if strings.HasSuffix(upstreamURL, ".git") {
		upstreamURL = strings.TrimSuffix(upstreamURL, ".git")
	}

	secretsURL := upstreamURL + "/settings/secrets/actions/new"
	secretsURL = fmt.Sprintf("\u001b]8;;%s\u0007%s\u001b]8;;\u0007", secretsURL, "github")

	fmt.Println("Add workflow secrets")
	fmt.Println("DOCKER_TAG:", tag)
	fmt.Println("DOCKER_USERNAME:", username)
	fmt.Println("DOCKER_PASSWORD:", password)

	fmt.Println()
	fmt.Println("Add your secrets on", secretsURL)
}
