package parser

import (
	"errors"
	"log"
	"os"

	"github.com/Phillezi/kthcloud-cli/pkg/v1/models/compose"
	"github.com/spf13/viper"
)

func GetCompose() (*compose.Compose, error) {
	files := []string{"kthcloud.docker-compose.yml", "kthcloud.docker-compose.yaml", "docker-compose.yml", "docker-compose.yaml"}
	specifiedFile := viper.GetString("file")
	if specifiedFile != "" {
		files = []string{specifiedFile}
	}

	var filePath string
	var fileFound bool

	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			filePath = file
			fileFound = true
			break
		}
	}

	if !fileFound {
		if specifiedFile != "" {
			log.Println("No file " + specifiedFile + " found in current directory.")
		} else {
			log.Println("No docker-compose file found.")

		}
		return nil, errors.New("no docker-compose file found in current directory")
	}

	composeInstance, err := compose.New(filePath)
	if err != nil {
		return nil, err
	}

	return composeInstance, nil
}
