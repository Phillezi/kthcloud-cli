package cmd

import (
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/build"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build and push to the deployment specified in .kthcloud/DEPLOYMENT",
	Run: func(cmd *cobra.Command, args []string) {
		build.Build()
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
