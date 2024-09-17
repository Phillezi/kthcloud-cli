package cmd

import (
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose"
	"github.com/spf13/cobra"
)

var composeCmd = &cobra.Command{
	Use:   "compose",
	Short: "Temp compose test",
	Run: func(cmd *cobra.Command, args []string) {
		compose.Up()
	},
}

func init() {
	rootCmd.AddCommand(composeCmd)
}
