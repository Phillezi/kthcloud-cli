package main

import (
	"kthcloud-cli/pkg/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "kthcloud-cli",
	Short: "CLI app for interacting with kthclouds REST API",
	Long: `   __    __    __         __                __             __   _ 
  / /__ / /_  / /  ____  / / ___  __ __ ___/ / ____ ____  / /  (_)
 /  '_// __/ / _ \/ __/ / / / _ \/ // // _  / /___// __/ / /  / / 
/_/\_\ \__/ /_//_/\__/ /_/  \___/\_,_/ \_,_/       \__/ /_/  /_/  
                                                                  `,
}

func init() {

	cobra.OnInitialize(config.InitConfig)

	// Persistent flags
	rootCmd.PersistentFlags().String("loglevel", "info", "Set the logging level (info, warn, error, debug)")
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))

	rootCmd.PersistentFlags().String("api-url", "https://api.cloud.cbh.kth.se/deploy", "Base URL of the API")
	viper.BindPFlag("api-url", rootCmd.PersistentFlags().Lookup("api-url"))
}
