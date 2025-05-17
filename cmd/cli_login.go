package cmd

import (
	"github.com/Phillezi/kthcloud-cli/internal/options"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to kthcloud using Keycloak and retrieve the authentication token",
	Long: `
Lets you log in to kthcloud by opening a local server on ":3000" and opening your browser to kthclouds keycloak server with the redirect uri set to the local http server so the accesstoken can be retrieved.

On Linux make sure you have "xdg-open" (should be included by default in most distros) and a browser for this to work.`,
	Example: "kthcloud login",
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
