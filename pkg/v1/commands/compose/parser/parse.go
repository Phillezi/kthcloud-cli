package parser

import (
	"errors"
	"log"
	"os"

	"github.com/Phillezi/kthcloud-cli/pkg/v1/models/compose"
)

func GetCompose() (*compose.Compose, error) {
	files := []string{"docker-compose.yml", "docker-compose.yaml"}

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
		log.Println("No docker-compose file found.")
		return nil, errors.New("no docker-compose file found in current directory")
	}

	composeInstance, err := compose.New(filePath)
	if err != nil {
		return nil, err
	}

	return composeInstance, nil
}
