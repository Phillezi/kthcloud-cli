package main

import (
	"kthcloud-cli/pkg/config"

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
}
