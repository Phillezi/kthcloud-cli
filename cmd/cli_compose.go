package cmd

import (
	"fmt"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/Phillezi/kthcloud-cli/internal/compose"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var composeCmd = &cobra.Command{
	Use:   "compose",
	Short: "Deploy and manage Docker Compose projects on the cloud",
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

		compose.PrintServices(services)
	},
}

var composeUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Deploy compose configuration to cloud",
	Run: func(cmd *cobra.Command, args []string) {
		// Look for docker-compose.yaml or docker-compose.yml
		composeFile, err := findComposeFile()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		compose.Up(composeFile)
	},
}

var composeDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop compose configuration to cloud",
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatal("not implemented yet...")
	},
}

func init() {
	composeUpCmd.Flags().BoolP("detached", "d", false, "doesn't do anything, just here for parity with Docker Compose up")
	viper.BindPFlag("detatched", composeUpCmd.Flags().Lookup("detatched"))

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
