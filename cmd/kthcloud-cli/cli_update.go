package main

import (
	"kthcloud-cli/internal/update"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var autoApprove bool

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update to latest version",
	Run: func(cmd *cobra.Command, args []string) {
		updated, err := update.Update(autoApprove)
		if err != nil {
			log.Fatal(err)
		}
		if updated {
			log.Infoln("Updated")
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVarP(&autoApprove, "yes", "y", false, "Skip prompt, automatically approve")
}
