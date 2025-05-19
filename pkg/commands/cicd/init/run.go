package init

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/Phillezi/kthcloud-cli/pkg/commands/cicd"
	"github.com/Phillezi/kthcloud-cli/pkg/file"
	"github.com/Phillezi/kthcloud-cli/pkg/git"
	"github.com/Phillezi/kthcloud-cli/pkg/github"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func (c *Command) Run() error {
	if c.client == nil {
		return fmt.Errorf("client is nil")
	}

	rootdir, upstreamURL, err := git.GetGitRepoInfo()
	if err != nil {
		logrus.Fatal("Error when getting repo info:", err)
		return err
	}
	repoConfDir := path.Join(rootdir, ".kthcloud")
	ghConfDir := path.Join(rootdir, ".github")
	wfConfDir := path.Join(ghConfDir, "workflows")

	var id string

	if !file.Exists(repoConfDir, "DEPLOYMENT") {
		name, err := util.GetNameFromUser()
		if err != nil {
			logrus.Fatal("Error when getting name from user:", err)
			return err
		}
		id, err = cicd.CreateEmptyDeployment(c.client, context.Background(), name)
		if err != nil {
			logrus.Fatal("Error when creating empty deployment:", err)
			return err
		}
	} else {
		id, err = file.Read(repoConfDir, "DEPLOYMENT")
		if err != nil {
			logrus.Fatal("Error when trying to get deployment id", err)
			return err
		}
	}

	if id != "" {
		err = file.Create(repoConfDir, "DEPLOYMENT", id)
		if err != nil {
			logrus.Fatal("Error when trying to save deployment id", err)
			return err
		}
	}

	conf, err := c.client.CiConfig(id)
	if err != nil {
		logrus.Fatal("Error when getting CI config:", err)
		return err
	}

	gha, err := github.ToModel(conf)
	if err != nil {
		return err
	}

	username, password, tag, err := github.ExtractSecrets(gha)
	if err != nil {
		logrus.Fatal("Error when extracting secrets:", err)
		return err
	}

	if c.saveSecrets {

		secrets := map[string]string{
			"username": username,
			"password": password,
			"tag":      tag,
		}

		secretsJson, err := json.Marshal(secrets)
		if err != nil {
			log.Fatalf("failed to marshal JSON: %v", err)
		}

		err = file.Create(repoConfDir, ".gitignore", "secrets.json")
		if err != nil {
			logrus.Fatal("Error when trying to add gitignore for secrets.json", err)
			return err
		}

		err = file.Create(repoConfDir, "secrets.json", string(secretsJson))
		if err != nil {
			logrus.Fatal("Error when trying to add secrets.json", err)
			return err
		}

	}

	yamlConf, err := yaml.Marshal(gha)
	if err != nil {
		logrus.Fatalf("error marshalling YAML: %v", err)
		return err
	}

	err = file.Create(wfConfDir, "kthcloud.yml", string(yamlConf))
	if err != nil {
		logrus.Fatal("Error when trying save cicd config", err)
		return err
	}

	if strings.Contains(upstreamURL, "github.com") {
		github.PromptUserAddSecrets(upstreamURL, username, password, tag)
	}

	return nil
}
