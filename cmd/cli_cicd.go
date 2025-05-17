package cmd

import (
	"github.com/Phillezi/kthcloud-cli/internal/options"
	cicd "github.com/Phillezi/kthcloud-cli/pkg/commands/cicd/init"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cicdCmd = &cobra.Command{
	Use:   "cicd",
	Short: "Generate CICD for gh repo",
	Long: `
This command allows you to setup cicd for your git repo, it creates a custom deployments and adds a workflow for building and pushing to the kthcloud registry on push / merge to main ".github/workflows/".`,
	Run: func(cmd *cobra.Command, args []string) {
		save, err := cmd.Flags().GetBool("save-secrets")
		if err != nil {
			logrus.Fatal(err)
		}
		if err := cicd.New(cicd.CommandOpts{
			Client:      options.DefaultClient(),
			SaveSecrets: &save,
		}).Run(); err != nil {
			logrus.Errorln(err)
			return
		}
	},
}

func init() {
	cicdCmd.Flags().BoolP("save-secrets", "j", false, "Save secrets in json file")
	rootCmd.AddCommand(cicdCmd)
}
