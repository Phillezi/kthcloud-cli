package cmd

import (
	"github.com/Phillezi/kthcloud-cli/internal/interrupt"
	"github.com/Phillezi/kthcloud-cli/internal/options"
	"github.com/Phillezi/kthcloud-cli/pkg/commands/vm/ps"
	"github.com/Phillezi/kthcloud-cli/pkg/commands/vm/ssh"
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
	Long: `
This command lets you list your VMs that are running, adding the -a or --all flag will list all VMs and wont filter to only the ones with resourceRunning status.`,
	Run: func(cmd *cobra.Command, args []string) {
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			logrus.Fatal(err)
		}
		if err := ps.New(ps.CommandOpts{
			Client: options.DefaultClient(),
			All:    &all,
		}).WithContext(interrupt.GetInstance().Context()).Run(); err != nil {
			logrus.Errorln(err)
			return
		}
	},
}

var vmSSHCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Connect to vm",
	Long: `
This command will let you ssh into a VM on kthcloud. It can be used by specifying the VMs name, id or if you dont specify a VM it will: if you only have one VM running, select that one, otherwise it will bring up a TUI selector that allows you to select the VM you want to ssh into.

What does it achieve? You dont have to look up the connectionstring yourself. In short it is basically just a wrapper around the ssh executable on your system that gets the connectionstring for you.

[!NOTE]: This requires you to have "ssh" installed on your machine.`,
	Example: "kthcloud vm ssh --name foo",
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
			logrus.Errorln(err)
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
