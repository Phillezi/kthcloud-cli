package cmd

import (
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/cicd"
	"github.com/spf13/cobra"
)

var cicdCmd = &cobra.Command{
	Use:   "cicd",
	Short: "Generate CICD for gh repo",
	Run: func(cmd *cobra.Command, args []string) {
		cicd.CICDInit()
	},
}

func init() {
	rootCmd.AddCommand(cicdCmd)
}
