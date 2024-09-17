package cmd

import (
	"log"

	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/ps"
	"github.com/spf13/cobra"
)

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List deployments",
	Run: func(cmd *cobra.Command, args []string) {
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			log.Fatal(err)
		}
		ps.Ps(all)
	},
}

func init() {
	psCmd.Flags().BoolP("all", "a", false, "Show all")
	rootCmd.AddCommand(psCmd)
}
