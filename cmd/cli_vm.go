package cmd

import (
	"github.com/Phillezi/kthcloud-cli/internal/interrupt"
	"github.com/Phillezi/kthcloud-cli/internal/options"
	"github.com/Phillezi/kthcloud-cli/pkg/commands/vm/ps"
	"github.com/Phillezi/kthcloud-cli/pkg/commands/vm/ssh"
	"github.com/kthcloud/go-deploy/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var vmCmd = &cobra.Command{
	Use:   "vm",
	Short: "Commands related to VMs",
}

var vmPsCmd = &cobra.Command{
	Use:   "ps",
	Short: "List VMs",
	Run: func(cmd *cobra.Command, args []string) {
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			logrus.Fatal(err)
		}
		if err := ps.New(ps.CommandOpts{
			Client: options.DefaultClient(),
			All:    &all,
		}).WithContext(interrupt.GetInstance().Context()).Run(); err != nil {
			log.Errorln(err)
			return
		}
	},
}

var vmSSHCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Connect to vm",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			logrus.Fatal(err)
		}
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			logrus.Fatal(err)
		}
		if err := ssh.New(ssh.CommandOpts{
			Client: options.DefaultClient(),
			ID:     &id,
			Name:   &name,
		}).WithContext(interrupt.GetInstance().Context()).Run(); err != nil {
			log.Errorln(err)
			return
		}
	},
}

func init() {
	vmPsCmd.Flags().BoolP("all", "a", false, "Show all")
	vmCmd.AddCommand(vmPsCmd)

	vmSSHCmd.Flags().StringP("name", "n", "", "Specify VM name to connect to")
	vmSSHCmd.Flags().StringP("id", "i", "", "Specify VM ID to connect to")
	vmCmd.AddCommand(vmSSHCmd)

	rootCmd.AddCommand(vmCmd)
}
