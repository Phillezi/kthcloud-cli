package main

import (
	"fmt"
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

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "See the version of the binary",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version: " + viper.GetString("release"))
	},
}

func init() {

	cobra.OnInitialize(config.InitConfig)

	// Persistent flags
	rootCmd.PersistentFlags().String("loglevel", "info", "Set the logging level (info, warn, error, debug)")
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))

	rootCmd.PersistentFlags().String("api-url", "https://api.cloud.cbh.kth.se/deploy", "Base URL of the API")
	viper.BindPFlag("api-url", rootCmd.PersistentFlags().Lookup("api-url"))

	rootCmd.Flags().StringP("api-token", "x", "", "kthcloud api token")
	viper.BindPFlag("api-token", loginCmd.Flags().Lookup("api-token"))

	rootCmd.Flags().StringP("zone", "z", "", "The preferred zone to use")
	viper.BindPFlag("zone", loginCmd.Flags().Lookup("zone"))

	rootCmd.Flags().StringP("session-path", "s", "session.json", "The filepath where the session should be loaded and saved to")
	viper.BindPFlag("session-path", loginCmd.Flags().Lookup("session-path"))

	viper.SetDefault("session-path", "session.json")

	rootCmd.AddCommand(versionCmd)

}
