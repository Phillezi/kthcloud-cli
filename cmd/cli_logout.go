package cmd

import (
	"github.com/Phillezi/kthcloud-cli/internal/options"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from kthcloud",
	Run: func(cmd *cobra.Command, args []string) {
		if err := options.DefaultClient().Auth().Logout(); err != nil {
			log.Errorln(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
