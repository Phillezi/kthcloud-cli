package compose

import (
	"fmt"
	"log"
	"os"

	"github.com/Phillezi/kthcloud-cli/pkg/v1/models/compose"
)

func Up() {
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
		return
	}

	composeInstance, err := compose.New(filePath)
	if err != nil {
		log.Printf("Error creating Compose instance: %v", err)
		return
	}

	log.Printf("Successfully created Compose instance from file: %s", filePath)

	fmt.Printf("Compose instance: %+v\n", composeInstance.String())
}
