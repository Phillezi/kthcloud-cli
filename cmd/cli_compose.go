package cmd

import (
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/compose/storage"
	"github.com/sirupsen/logrus"
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
		json, err := cmd.Flags().GetBool("json")
		if err != nil {
			logrus.Fatal(err)
		}
		compose.Parse(json)
	},
}

var composeUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Deploy compose configuration to cloud",
	Run: func(cmd *cobra.Command, args []string) {
		detached, _ := cmd.Flags().GetBool("detached")
		tryToCreateVolumes, _ := cmd.Flags().GetBool("try-volumes")
		servicesToBuild, _ := cmd.Flags().GetStringSlice("build")
		nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
		buildAll := false
		if len(servicesToBuild) == 1 && servicesToBuild[0] == "__all__" {
			logrus.Debugln("build all is set")
			buildAll = true
		}
		compose.Up(detached, tryToCreateVolumes, buildAll, nonInteractive, servicesToBuild)
	},
}
var composeDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Bring down compose configuration to cloud",
	Run: func(cmd *cobra.Command, args []string) {
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			logrus.Fatal(err)
		}
		volumes, err := cmd.Flags().GetBool("volumes")
		if err != nil {
			logrus.Fatal(err)
		}
		compose.Down(all, volumes)
	},
}
var composeStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop compose configuration to cloud",
	Run: func(cmd *cobra.Command, args []string) {
		compose.Stop()
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
	composeCmd.PersistentFlags().StringP("file", "", "", "Specify which docker-compose file to use.")
	viper.BindPFlag("file", composeCmd.PersistentFlags().Lookup("file"))

	composeUpCmd.Flags().BoolP("try-volumes", "", false, "Try uploading local files and dirs that should be mounted on the deployment.\nIf enabled it will \"steal\" cookies from your browser to authenticate.")
	composeUpCmd.Flags().BoolP("detached", "d", false, "Run detached, default behaviour attaches logs from the deployments.")
	viper.BindPFlag("detached", composeUpCmd.Flags().Lookup("detached"))

	composeUpCmd.Flags().Bool("non-interactive", false, "Yes to all options and run non interactively")

	composeUpCmd.Flags().StringSlice("build", nil, "Build services and push to registry, can be followed by service name to specify which should be built")
	buildF := composeUpCmd.Flags().Lookup("build")
	buildF.NoOptDefVal = "__all__"

	composeDownCmd.Flags().BoolP("all", "a", false, "Remove all")
	composeDownCmd.Flags().BoolP("volumes", "v", false, "Remove volumes")

	composeParseCmd.Flags().Bool("json", false, "Specify if output should be only the parsed json")

	// Register subcommands with the main compose command
	composeCmd.AddCommand(composeParseCmd)
	composeCmd.AddCommand(composeUpCmd)
	composeCmd.AddCommand(composeDownCmd)
	composeCmd.AddCommand(composeStopCmd)
	composeCmd.AddCommand(composeLogsCmd)
	composeCmd.AddCommand(testSMAuthCmd)

	// Register the compose command in root
	rootCmd.AddCommand(composeCmd)
}
