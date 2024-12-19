package cmd

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/config"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "kthcloud",
	Short: "CLI app for interacting with kthclouds REST API",
	Long: `   __    __    __         __                __             __   _ 
  / /__ / /_  / /  ____  / / ___  __ __ ___/ / ____ ____  / /  (_)
 /  '_// __/ / _ \/ __/ / / / _ \/ // // _  / /___// __/ / /  / / 
/_/\_\ \__/ /_//_/\__/ /_/  \___/\_,_/ \_,_/       \__/ /_/  /_/  
                                                                  `,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		level := viper.GetString("loglevel")
		lvl, err := logrus.ParseLevel(level)
		if err != nil {
			logrus.Warnf("Invalid log level %s, falling back to INFO", level)
			lvl = logrus.InfoLevel
		}
		logrus.SetLevel(lvl)

		logrus.Debugf("Logging level set to %s", lvl)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "See the version of the binary",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version: " + viper.GetString("release"))
	},
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "See who you are",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.Get()
		if !c.HasValidSession() {
			fmt.Println("I dont know...")
			return
		}
		user, err := c.User()
		if err != nil {
			fmt.Println("I dont know...")
			return
		}
		fmt.Println("Name:\t" + strings.Split(user.FirstName, " ")[0] + " " + user.LastName + "\n\nEmail:\t" + user.Email + "\nRole:\t" + user.Role.Name)

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
	viper.BindPFlag("api-token", rootCmd.PersistentFlags().Lookup("api-token"))

	rootCmd.Flags().StringP("zone", "z", "", "The preferred zone to use")
	viper.BindPFlag("zone", rootCmd.PersistentFlags().Lookup("zone"))

	rootCmd.Flags().StringP("session-path", "s", path.Join(config.GetConfigPath(), "session.json"), "The filepath where the session should be loaded and saved to")
	viper.BindPFlag("session-path", rootCmd.PersistentFlags().Lookup("session-path"))

	viper.SetDefault("session-path", path.Join(config.GetConfigPath(), "session.json"))

	rootCmd.Flags().DurationP("resource-cache-duration", "c", 60*time.Second, "How long resources should be cached when possible")
	viper.BindPFlag("resource-cache-duration", rootCmd.PersistentFlags().Lookup("resource-cache-duration"))

	rootCmd.AddCommand(versionCmd)

	rootCmd.AddCommand(whoamiCmd)

}
