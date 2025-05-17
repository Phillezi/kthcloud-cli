package cmd

import (
	"github.com/Phillezi/kthcloud-cli/internal/interrupt"
	"github.com/Phillezi/kthcloud-cli/internal/options"
	"github.com/Phillezi/kthcloud-cli/pkg/commands/build"
	"github.com/kthcloud/go-deploy/pkg/log"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build and push to the deployment specified in .kthcloud/DEPLOYMENT",
	Long: `
Allows you to build and push your custom deployment locally using docker.

It uses the file ".kthcloud/DEPLOYMENT" to figure out which deployment your cwd maps to. This file contains the UUID of the deployment.

[!NOTE]: Requires you to have docker with buildx installed.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := build.New(build.CommandOpts{
			Client: options.DefaultClient(),
		}).WithContext(interrupt.GetInstance().Context()).Run(); err != nil {
			log.Errorln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
