package cicd

import (
	"context"
	"encoding/json"
	"log"
	"path"
	"strings"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func CICDInit() {
	rootdir, upstreamURL, err := GetGitRepoInfo()
	if err != nil {
		logrus.Fatal("Error when getting repo info:", err)
		return
	}
	repoConfDir := path.Join(rootdir, ".kthcloud")
	ghConfDir := path.Join(rootdir, ".github")
	wfConfDir := path.Join(ghConfDir, "workflows")

	var id string

	if !FileExists(repoConfDir, "DEPLOYMENT") {
		name, err := util.GetNameFromUser()
		if err != nil {
			logrus.Fatal("Error when getting name from user:", err)
			return
		}
		id, err = createDeployment(context.Background(), name)
		if err != nil {
			logrus.Fatal("Error when creating empty deployment:", err)
			return
		}
	} else {
		id, err = ReadFile(repoConfDir, "DEPLOYMENT")
		if err != nil {
			logrus.Fatal("Error when trying to get deployment id", err)
			return
		}
	}

	if id != "" {
		err = CreateFile(repoConfDir, "DEPLOYMENT", id)
		if err != nil {
			logrus.Fatal("Error when trying to save deployment id", err)
			return
		}
	}

	conf, err := getGHACIConf(id)
	if err != nil {
		logrus.Fatal("Error when getting CI config:", err)
		return
	}

	username, password, tag, err := extractSecrets(conf)
	if err != nil {
		logrus.Fatal("Error when extracting secrets:", err)
		return
	}

	secrets := map[string]string{
		"username": username,
		"password": password,
		"tag":      tag,
	}

	secretsJson, err := json.Marshal(secrets)
	if err != nil {
		log.Fatalf("failed to marshal JSON: %v", err)
	}

	err = CreateFile(repoConfDir, ".gitignore", "secrets.json")
	if err != nil {
		logrus.Fatal("Error when trying to add gitignore for secrets.json", err)
		return
	}

	err = CreateFile(repoConfDir, "secrets.json", string(secretsJson))
	if err != nil {
		logrus.Fatal("Error when trying to add secrets.json", err)
		return
	}

	yamlConf, err := yaml.Marshal(conf)
	if err != nil {
		logrus.Fatalf("error marshalling YAML: %v", err)
		return
	}

	err = CreateFile(wfConfDir, "kthcloud.yml", string(yamlConf))
	if err != nil {
		logrus.Fatal("Error when trying save cicd config", err)
		return
	}

	if strings.Contains(upstreamURL, "github.com") {
		promptUserAddSecrets(upstreamURL, username, password, tag)
	}
}
