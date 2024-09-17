package cmd

import (
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose"
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
		tryToCreateVolumes, _ := cmd.Flags().GetBool("try-volumes")
		compose.Up(tryToCreateVolumes)
	},
}
var composeDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop compose configuration to cloud",
	Run: func(cmd *cobra.Command, args []string) {
		compose.Down()
	},
}

func init() {
	composeUpCmd.Flags().BoolP("try-volumes", "", false, "Try to create volumes despite auth not working for it yet")
	composeUpCmd.Flags().BoolP("detached", "d", false, "doesn't do anything, just here for parity with Docker Compose up")
	viper.BindPFlag("detached", composeUpCmd.Flags().Lookup("detached"))

	// Register subcommands with the main compose command
	composeCmd.AddCommand(composeParseCmd)
	composeCmd.AddCommand(composeUpCmd)
	composeCmd.AddCommand(composeDownCmd)

	// Register the compose command in root
	rootCmd.AddCommand(composeCmd)
}
