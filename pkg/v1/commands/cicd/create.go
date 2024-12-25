package cicd

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func Create(rootdir string, createWF bool, name string) {
	gitDir, upstreamURL, err := GetGitRepoInfoFrom(rootdir)
	if err != nil {
		gitDir = rootdir
	}
	repoConfDir := path.Join(rootdir, ".kthcloud")
	ghConfDir := path.Join(gitDir, ".github")
	wfConfDir := path.Join(ghConfDir, "workflows")

	var id string

	if !FileExists(repoConfDir, "DEPLOYMENT") {
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
		if !client.Get().DeploymentExists(id) {
			logrus.Debugln("file exists but contains ID of deployment that doesnt")
			id, err = createDeployment(context.Background(), name)
			if err != nil {
				logrus.Fatal("Error when creating empty deployment:", err)
				return
			}
		}
	}

	if id != "" {
		err = CreateFile(repoConfDir, "DEPLOYMENT", id)
		if err != nil {
			logrus.Fatal("Error when trying to save deployment id", err)
			return
		}
	}

	if !createWF {
		return
	}

	conf, err := GetGHACIConf(id)
	if err != nil {
		logrus.Fatal("Error when getting CI config:", err)
		return
	}

	username, password, tag, err := ExtractSecrets(conf)
	if err != nil {
		logrus.Fatal("Error when trying remove secrets from workflow", err)
		return
	}

	yamlConf, err := yaml.Marshal(conf)
	if err != nil {
		logrus.Fatalf("error marshalling YAML: %v", err)
		return
	}
	filename := "kthcloud.yml"
	if name != "" {
		filename = "kthcloud-" + name + "-cicd.yml"
	}
	for i := 0; FileExists(wfConfDir, filename); i++ {
		filename = fmt.Sprintf("kthcloud-%s-cicd-%d.yml", name, i)
	}

	err = CreateFile(wfConfDir, filename, string(yamlConf))
	if err != nil {
		logrus.Fatal("Error when trying save cicd config", err)
		return
	}

	if strings.Contains(upstreamURL, "github.com") {
		promptUserAddSecrets(upstreamURL, username, password, tag)
	}

	for checks := 0; checks < 5; checks++ {
		if FileExists(repoConfDir, "DEPLOYMENT") {
			break
		}
		logrus.Debugln("waiting on fs...")
		time.Sleep(100 * time.Millisecond)
	}

}
