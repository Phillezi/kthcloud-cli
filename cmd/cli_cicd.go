package cmd

import (
	cicd "github.com/Phillezi/kthcloud-cli/pkg/commands/cicd/init"
	"github.com/kthcloud/go-deploy/pkg/log"
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
		if err := cicd.New(cicd.CommandOpts{
			SaveSecrets: &save,
		}).Run(); err != nil {
			log.Errorln(err)
			return
		}
	},
}

func init() {
	cicdCmd.Flags().BoolP("save-secrets", "j", false, "Save secrets in json file")
	rootCmd.AddCommand(cicdCmd)
}
