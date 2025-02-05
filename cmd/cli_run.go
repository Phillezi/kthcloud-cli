package cmd

import (
	"strings"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/run"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/models/options"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	envVars         []string
	ports           []string
	volumes         []string
	containerName   string
	interactiveLogs bool
	removeOnExit    bool
	visibility      string
	cores           float64
	ram             float64
	replicas        int
	healthCheck     string
)

var runCmd = &cobra.Command{
	Use:   "run [flags] IMAGE [COMMAND] [ARGS...]",
	Short: "Create a deployment",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Extract image from args
		image := args[0]
		entrypoint := ""
		if len(args) > 1 {
			entrypoint = strings.Join(args[1:], " ")
		}

		if removeOnExit && !interactiveLogs {
			logrus.Warn("Won't be able to determine when to kill the container if not interactive. --rm must be used with -i.")
			removeOnExit = false
		}

		if containerName == "" {
			containerName = util.GenerateRandomName()
		} else if !util.IsValidDeplName(containerName) {
			logrus.Fatal("Invalid name")
		}

		options := &options.DeploymentOptions{
			Image:           image,
			Entrypoint:      util.OrNil(entrypoint),
			ContainerName:   containerName,
			InteractiveLogs: interactiveLogs,
			RemoveOnExit:    removeOnExit,
			Visibility:      visibility,
			Cores:           cores,
			Ram:             ram,
			Replicas:        replicas,
			HealthCheck:     healthCheck,
			EnvVars:         envVars,
			Ports:           ports,
			Volumes:         volumes,
		}

		run.Run(options)
	},
}

func init() {
	runCmd.Flags().StringVarP(&containerName, "name", "n", "", "Assign a name to the container")
	runCmd.Flags().BoolVarP(&interactiveLogs, "interactive", "i", false, "Keep the container in interactive mode")
	runCmd.Flags().BoolVarP(&removeOnExit, "rm", "", false, "Automatically remove the container when exiting interactive mode. Must be used with interactive mode enabled")
	runCmd.Flags().StringArrayVarP(&envVars, "env", "e", []string{}, "Set environment variables (e.g., -e KEY=VALUE)")
	runCmd.Flags().StringArrayVarP(&ports, "publish", "p", []string{}, "Publish container ports (e.g., -p 8080:80)")
	runCmd.Flags().StringArrayVarP(&volumes, "volume", "v", []string{}, "Mount volumes (e.g., -v /host/path:/container/path)")

	// KTHCloud specific
	runCmd.Flags().StringVar(&visibility, "visibility", "public", "Set container visibility (private, public, or auth)")
	runCmd.Flags().Float64Var(&cores, "cores", 0.2, "Number of CPU cores to allocate (default: 0.2)")
	runCmd.Flags().Float64Var(&ram, "ram", 0.5, "Amount of RAM to allocate (default: 0.5 GB)")
	runCmd.Flags().IntVar(&replicas, "replicas", 1, "Number of container replicas (default: 1)")
	runCmd.Flags().StringVar(&healthCheck, "health-check", "/healthz", "Path for health check (default: /healthz)")

	rootCmd.AddCommand(runCmd)
}
