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
	Long: `
Compose lets you describe multiple deployments in yaml files similary to how you would do it with docker compose.

The services described in the yaml files can be brought up on the cloud, you can bring them down, stop them or see logs from them.`,
}

var composeParseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse a docker-compose.yaml or docker-compose.yml file",
	Long: `
A simple command that parses a compose file into kthclouds create deployment json format.

The compose file can be specified with the --file flag.

Adding the --json flag will return pure json as response.`,
	Example: "kthcloud compose parse",
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
	Long: `
This command allows you to bring up a local compose definition to kthcloud.

It will automatically detect if you have custom deployments in your config that needs to be built.
[!NOTE]: To keep track of custom deployments a file in "./.kthcloud/DEPLOYMENT" is made (relative to the build context of the service). This file contains the UUID of the deployment that the service maps to.

To rebuild services you can apply the --build flag to rebuild all services, or specify the ones you want to rebuild.

Volumes can be defined and managed, but it is done in a hacky way. Since the storage (filebrowser) instance that gets deployed for your user is behind a oauth2 proxy that only uses cookies for authentication. The current solution is as mention very hacky as it will try to scrape your browser for these cookies to authenticate. This only works on Linux and MacOS, but it can be unreliable since browser companies dont want you to be able too scrape stuff like this.

Default behaviour of this command will after creating all the deployments setup SSE log streams from all created deployments to your terminal. This can be skipped by adding the -d flag.`,
	Example: `kthcloud compose up \
	--file ./somesubdir/example.compose.yml \
	--try-volumes`,
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
	Long: `
This command will bring down all services specified in the compose file. By default custom deployments wont be brought down, this can be changed with the -a flag.`,
	Example: "kthcloud compose --file ./compose.yaml down -a",
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
	Long: `
This command will disable all services specified in the compose file. This is done my setting their replicas to zero.

[!NOTE]: At the time of writing this there is currently a but that when trying to re-enable these services later requires you to change more than just setting replicas back to 1 or whatever value it was before.`,
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
	Long: `
This command allows you to get all logs from the services specified in the compose file.`,
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
