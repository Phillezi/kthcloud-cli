package cmd

import (
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

func init() {
	vmPsCmd.Flags().BoolP("all", "a", false, "Show all")
	vmCmd.AddCommand(vmPsCmd)

	rootCmd.AddCommand(vmCmd)
}
