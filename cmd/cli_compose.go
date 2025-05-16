package cmd

import (
	"github.com/Phillezi/kthcloud-cli/internal/interrupt"
	"github.com/Phillezi/kthcloud-cli/internal/options"
	"github.com/Phillezi/kthcloud-cli/pkg/commands/compose/down"
	"github.com/Phillezi/kthcloud-cli/pkg/commands/compose/logs"
	"github.com/Phillezi/kthcloud-cli/pkg/commands/compose/parse"
	"github.com/Phillezi/kthcloud-cli/pkg/commands/compose/stop"
	"github.com/Phillezi/kthcloud-cli/pkg/commands/compose/up"
	"github.com/Phillezi/kthcloud-cli/pkg/parser"
	"github.com/Phillezi/kthcloud-cli/pkg/storage"
	"github.com/kthcloud/go-deploy/pkg/log"
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
		compose, err := parser.GetCompose()
		if err != nil {
			logrus.Error(err)
			return
		}
		if err := parse.New(parse.CommandOpts{
			Client:  options.DefaultClient(),
			Compose: compose,
			Json:    &json,
		}).WithContext(interrupt.GetInstance().Context()).Run(); err != nil {
			logrus.Errorln(err)
			return
		}
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

		compose, err := parser.GetCompose()
		if err != nil {
			logrus.Errorln(err)
			return
		}
		if err := up.New(up.CommandOpts{
			Client:  options.DefaultClient(),
			Compose: compose,
			//Services: ,
			ServicesToBuild: servicesToBuild,
			BuildAll:        &buildAll,
			Detatched:       &detached,
			TryVolumes:      &tryToCreateVolumes,
			NonInteractive:  &nonInteractive,
		}).WithContext(interrupt.GetInstance().Context()).Run(); err != nil {
			log.Errorln(err)
			return
		}
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
		compose, err := parser.GetCompose()
		if err != nil {
			logrus.Errorln(err)
			return
		}
		if err := down.New(down.CommandOpts{
			Client:  options.DefaultClient(),
			Compose: compose,
			All:     &all,
		}).WithContext(interrupt.GetInstance().Context()).Run(); err != nil {
			logrus.Errorln(err)
			return
		}
	},
}
var composeStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop compose configuration to cloud",
	Run: func(cmd *cobra.Command, args []string) {
		compose, err := parser.GetCompose()
		if err != nil {
			logrus.Errorln(err)
			return
		}
		if err := stop.New(stop.CommandOpts{
			Client:  options.DefaultClient(),
			Compose: compose,
		}).WithContext(interrupt.GetInstance().Context()).Run(); err != nil {
			logrus.Errorln(err)
			return
		}
	},
}
var composeLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Get logs from deployments in the compose file",
	Run: func(cmd *cobra.Command, args []string) {
		compose, err := parser.GetCompose()
		if err != nil {
			logrus.Errorln(err)
			return
		}
		if err := logs.New(logs.CommandOpts{
			Client:  options.DefaultClient(),
			Compose: compose,
		}).WithContext(interrupt.GetInstance().Context()).Run(); err != nil {
			logrus.Errorln(err)
			return
		}
	},
}

var testSMAuthCmd = &cobra.Command{
	Use:    "sm check",
	Short:  "Test authentication against storage manager",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		storage.Check(options.DefaultClient())
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
