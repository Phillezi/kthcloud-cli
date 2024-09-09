package main

import (
	"fmt"
	"log"
	"path/filepath"

	"kthcloud-cli/internal/compose"

	"github.com/spf13/cobra"
)

var composeCmd = &cobra.Command{
	Use:   "compose",
	Short: "Manage Docker Compose files",
}

var composeParseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse a docker-compose.yaml or docker-compose.yml file",
	Run: func(cmd *cobra.Command, args []string) {
		// Look for docker-compose.yaml or docker-compose.yml
		composeFile, err := findComposeFile()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		// Parse the file
		services, err := compose.ParseComposeFile(composeFile)
		if err != nil {
			log.Fatalf("Failed to parse compose file: %v", err)
		}

		// Output the parsed data
		for name, service := range services {
			fmt.Printf("Service: %s\n", name)
			fmt.Printf("Environment Variables: %v\n", service.Environment)
			fmt.Printf("Ports: %v\n", service.Ports)
			fmt.Printf("Volumes: %v\n", service.Volumes)
			fmt.Println("----------------------------")
		}
	},
}

var composeUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Start Docker Compose services",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting Docker Compose services...")
		// Implement logic to start Docker Compose services
	},
}

var composeDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop Docker Compose services",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Stopping Docker Compose services...")
		// Implement logic to stop Docker Compose services
	},
}

func init() {
	// Register subcommands with the main compose command
	composeCmd.AddCommand(composeParseCmd)
	composeCmd.AddCommand(composeUpCmd)
	composeCmd.AddCommand(composeDownCmd)

	// Register the compose command in root
	rootCmd.AddCommand(composeCmd)
}

// Helper function to find the compose file
func findComposeFile() (string, error) {
	// Search for docker-compose.yaml or docker-compose.yml
	files := []string{"docker-compose.yaml", "docker-compose.yml"}
	for _, file := range files {
		matches, err := filepath.Glob(file)
		if err != nil {
			return "", err
		}
		if len(matches) > 0 {
			return matches[0], nil
		}
	}
	return "", fmt.Errorf("docker-compose.yaml or docker-compose.yml not found")
}
