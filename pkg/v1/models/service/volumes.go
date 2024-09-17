package service

import (
	"go-deploy/dto/v2/body"
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
