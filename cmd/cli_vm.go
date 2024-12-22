package cmd

import (
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/vm/connect"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/vm/ps"
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
		ps.List(all)
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
		connect.SSH(id, name)
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
