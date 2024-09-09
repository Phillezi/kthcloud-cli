package main

import (
	"kthcloud-cli/internal/api"
	"kthcloud-cli/pkg/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "kthcloud-cli",
	Short: "CLI app for interacting with a REST API",
	Long:  `CLI application that uses logrus for logging, cobra for CLI, and viper for configuration.`,
}

func init() {

	cobra.OnInitialize(config.InitConfig)

	// Persistent flags
	rootCmd.PersistentFlags().String("loglevel", "info", "Set the logging level (info, warn, error, debug)")
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))

	rootCmd.PersistentFlags().String("api-url", "https://api.cloud.cbh.kth.se/deploy", "Base URL of the API")
	viper.BindPFlag("api-url", rootCmd.PersistentFlags().Lookup("api-url"))

	rootCmd.AddCommand(apiCmd)
}

var apiCmd = &cobra.Command{
	Use:   "get [resource]",
	Short: "Fetch data from the API with authorization",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Set logging level from config
		logLevel, err := log.ParseLevel(viper.GetString("loglevel"))
		if err != nil {
			log.Fatal(err)
		}
		log.SetLevel(logLevel)

		resource := args[0]
		apiURL := viper.GetString("api-url")

		// Retrieve the auth token from viper
		token := viper.GetString("auth-token")
		if token == "" {
			log.Fatal("No authentication token found. Please log in first.")
		}

		// Use the token in the request header
		client := api.NewClient(apiURL, token)

		resp, err := client.FetchResource(resource, "GET")
		if err != nil {
			log.Fatalf("Failed to fetch resource: %v", err)
		}

		log.Infof("Response: %s", resp)
	},
}
