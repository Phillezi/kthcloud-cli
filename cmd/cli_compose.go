package cmd

import (
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var composeCmd = &cobra.Command{
	Use:   "compose",
	Short: "Deploy and manage Docker Compose projects on the cloud",
}

var composeParseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse a docker-compose.yaml or docker-compose.yml file",
	Run: func(cmd *cobra.Command, args []string) {
		compose.Parse()
	},
}

var composeUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Deploy compose configuration to cloud",
	Run: func(cmd *cobra.Command, args []string) {
		detached, _ := cmd.Flags().GetBool("detached")
		tryToCreateVolumes, _ := cmd.Flags().GetBool("try-volumes")
		compose.Up(detached, tryToCreateVolumes)
	},
}
var composeDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop compose configuration to cloud",
	Run: func(cmd *cobra.Command, args []string) {
		compose.Down()
	},
}
var composeLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Get logs from deployments in the compose file",
	Run: func(cmd *cobra.Command, args []string) {
		compose.Logs()
	},
}

var testSMAuthCmd = &cobra.Command{
	Use:    "sm check",
	Short:  "Test authentication against storage manager",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		storage.Check()
	},
}

func init() {
	composeUpCmd.Flags().BoolP("try-volumes", "", false, "Try uploading local files and dirs that should be mounted on the deployment.\nIf enabled it will \"steal\" cookies from your browser to authenticate.")
	composeUpCmd.Flags().BoolP("detached", "d", false, "Run detached, default behaviour attaches logs from the deployments")
	viper.BindPFlag("detached", composeUpCmd.Flags().Lookup("detached"))

	// Register subcommands with the main compose command
	composeCmd.AddCommand(composeParseCmd)
	composeCmd.AddCommand(composeUpCmd)
	composeCmd.AddCommand(composeDownCmd)
	composeCmd.AddCommand(composeLogsCmd)
	composeCmd.AddCommand(testSMAuthCmd)

	// Register the compose command in root
	rootCmd.AddCommand(composeCmd)
}
