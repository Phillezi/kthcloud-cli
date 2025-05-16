package cmd

import (
	"github.com/Phillezi/kthcloud-cli/internal/interrupt"
	"github.com/Phillezi/kthcloud-cli/internal/options"
	"github.com/Phillezi/kthcloud-cli/pkg/commands/ps"
	"github.com/Phillezi/kthcloud-cli/pkg/deploy"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List deployments",
	Run: func(cmd *cobra.Command, args []string) {
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			logrus.Errorln(err)
			return
		}
		if err := ps.New(ps.CommandOpts{
			All: &all,
			Client: deploy.GetInstance(
				options.DeployOpts(),
			).WithContext(
				interrupt.GetInstance().Context(),
			),
		}).WithContext(
			interrupt.GetInstance().Context(),
		).Run(); err != nil {
			logrus.Errorln(err)
			return
		}
	},
}

func init() {
	psCmd.Flags().BoolP("all", "a", false, "Show all")
	rootCmd.AddCommand(psCmd)
}
