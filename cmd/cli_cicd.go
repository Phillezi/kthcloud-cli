package cmd

import (
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/cicd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cicdCmd = &cobra.Command{
	Use:   "cicd",
	Short: "Generate CICD for gh repo",
	Run: func(cmd *cobra.Command, args []string) {
		save, err := cmd.Flags().GetBool("save-secrets")
		if err != nil {
			logrus.Fatal(err)
		}
		cicd.CICDInit(save)
	},
}

func init() {
	cicdCmd.Flags().BoolP("save-secrets", "j", false, "Save secrets in json file")
	rootCmd.AddCommand(cicdCmd)
}
