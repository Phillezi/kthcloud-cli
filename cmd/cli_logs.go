package cmd

import (
	"github.com/Phillezi/kthcloud-cli/internal/interrupt"
	"github.com/Phillezi/kthcloud-cli/internal/options"
	"github.com/Phillezi/kthcloud-cli/pkg/commands/logs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs [services]",
	Short: "Get logs from deployments you specify",
	Long: `
Allows you to get logs from specified deployments.`,
	Args:    cobra.MinimumNArgs(1),
	Example: "kthcloud logs myapp myotherapp",
	Run: func(cmd *cobra.Command, args []string) {
		if err := logs.New(logs.CommandOpts{
			Client:   options.DefaultClient(),
			Services: args,
		}).WithContext(interrupt.GetInstance().Context()).Run(); err != nil {
			logrus.Errorln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
