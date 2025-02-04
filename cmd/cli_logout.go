package cmd

import (
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from kthcloud",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.Get()

		err := c.Logout()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
