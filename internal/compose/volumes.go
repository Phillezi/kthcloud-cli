package compose

import (
	"errors"
	"go-deploy/dto/v2/body"
	"kthcloud-cli/internal/api"
	"kthcloud-cli/internal/model"
	"strings"

	log "github.com/sirupsen/logrus"
)

// TODO: Check if "volume" is a file, since docker compose allows this but kthcloud does not
func ToVolumes(volumes []string, projectRoot string) []body.Volume {
	var parsedVolumes []body.Volume

	for _, volume := range volumes {
		parts := strings.Split(volume, ":")
		v := body.Volume{
			Name:       "kth-cli-generated",
			ServerPath: toServerPath(parts[0], projectRoot),
			AppPath:    parts[0],
		}

		if len(parts) > 1 {
			v.AppPath = parts[1]
		}

		if len(parts) > 2 {
			log.Warnln("extra volume info, omitting:", strings.Join(parts[1:], " "))
		}

		parsedVolumes = append(parsedVolumes, v)
	}

	return parsedVolumes
}

func toServerPath(path string, root string) string {
	if root == "" {
		return path
	}
	return root + "/" + path
}

func CreateVolume(session *model.Session, services map[string]model.Service) (string, error) {
	projectDir := model.Hash(services)
	if session.User == nil {
		err := session.FetchUser()
		if err != nil {
			return "", err
		}
	}
	storageURL := session.User.StorageURL
	if storageURL == nil {
		return "", errors.New("user does not have a storageURL")
	}

	if session.AuthSession == nil || session.AuthSession.Token == "" || session.AuthSession.IsExpired() {
		return "", errors.New("volume creation requires authentication, please log in")
	}

	// might need to set "X-Auth" header for authentication here
	client := api.NewClient(*storageURL, session.AuthSession.Token)
	resp, err := client.Req("api/resources/"+projectDir+"/?override=false", "POST", nil)
	if err != nil {
		return "", err
	}
	if resp.IsError() {
		return "", errors.New(resp.String())
	}

	return projectDir, nil
}
