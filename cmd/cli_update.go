package cmd

import (
	"github.com/Phillezi/kthcloud-cli/internal/update"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var autoApprove bool
var interactive bool

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update to latest version",
	Run: func(cmd *cobra.Command, args []string) {
		updated, err := update.Update(autoApprove, interactive)
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
	updateCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Select a specific version")
}
