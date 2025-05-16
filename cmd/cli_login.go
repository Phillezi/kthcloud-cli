package cmd

import (
	"github.com/Phillezi/kthcloud-cli/internal/options"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to kthcloud using Keycloak and retrieve the authentication token",
	Run: func(cmd *cobra.Command, args []string) {
		c := options.DefaultClient()
		_, err := c.Login()
		if err != nil {
			logrus.Fatal(err)
		}
		if c.HasValidSession() {
			logrus.Info("Logged in")
		} else {
			logrus.Fatal("Could not login")
		}
	},
}

func init() {
	// Add the login command
	rootCmd.AddCommand(loginCmd)
}
