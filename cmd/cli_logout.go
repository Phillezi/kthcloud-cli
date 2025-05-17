package cmd

import (
	"github.com/Phillezi/kthcloud-cli/internal/options"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from kthcloud",
	Long: `
Removes your current session.`,
	Example: "kthcloud logout",
	Run: func(cmd *cobra.Command, args []string) {
		if err := options.DefaultClient().Auth().Logout(); err != nil {
			logrus.Errorln(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
