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
	Long: `
This command lets you update or change your version of the cli. It fetched the github api to find out if the current version is the latest. If your version is older than the latest release it will prompt you to let you select if you want to update or not. This can be auto-accepted with the -y option.

This command also lets you change back to a older version using the -i option. It will bring up a TUI that lists all the available versions where you can select an older version. This can be useful if I accidentally break something that previously worked ;).`,
	Example: "kthcloud update -i",
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
