package cmd

import (
	"github.com/Phillezi/kthcloud-cli/internal/interrupt"
	"github.com/Phillezi/kthcloud-cli/internal/options"
	"github.com/Phillezi/kthcloud-cli/pkg/commands/run"
	"github.com/Phillezi/kthcloud-cli/pkg/defaults"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a container",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		i, err := cmd.Flags().GetBool("interactive")
		if err != nil {
			logrus.Errorln("Error parsing flag 'interactive':", err)
			return
		}
		t, err := cmd.Flags().GetBool("tty")
		if err != nil {
			logrus.Errorln("Error parsing flag 'tty':", err)
			return
		}
		rm, err := cmd.Flags().GetBool("rm")
		if err != nil {
			logrus.Errorln("Error parsing flag 'rm':", err)
			return
		}
		d, err := cmd.Flags().GetBool("detatch")
		if err != nil {
			logrus.Errorln("Error parsing flag 'detatch':", err)
			return
		}
		e, err := cmd.Flags().GetStringToString("env")
		if err != nil {
			logrus.Errorln("Error parsing flag 'env':", err)
			return
		}
		p, err := cmd.Flags().GetIntSlice("port")
		if err != nil {
			logrus.Errorln("Error parsing flag 'port':", err)
			return
		}
		v, err := cmd.Flags().GetString("visibility")
		if err != nil {
			logrus.Errorln("Error parsing flag 'visibility':", err)
			return
		}
		m, err := cmd.Flags().GetFloat64("memory")
		if err != nil {
			logrus.Errorln("Error parsing flag 'memory':", err)
			return
		}
		c, err := cmd.Flags().GetFloat64("cores")
		if err != nil {
			logrus.Errorln("Error parsing flag 'cores':", err)
			return
		}
		r, err := cmd.Flags().GetInt("replicas")
		if err != nil {
			logrus.Errorln("Error parsing flag 'replicas':", err)
			return
		}
		n, err := cmd.Flags().GetString("name")
		if err != nil {
			logrus.Errorln("Error parsing flag 'name':", err)
			return
		}

		if err := run.New(run.CommandOpts{
			Client:      options.DefaultClient(),
			Image:       &args[0],
			Interactive: &i,
			TTY:         &t,
			Remove:      &rm,
			Detatch:     &d,
			Envs:        e,
			Port:        p,
			Visibility:  &v,
			Memory:      &m,
			Cores:       &c,
			Replicas:    &r,
			Name:        &n,
		}).WithContext(interrupt.GetInstance().Context()).Run(); err != nil {
			logrus.Errorln(err)
		}
	},
}

func init() {
	runCmd.Flags().BoolP("interactive", "i", defaults.DefaultDeploymentRunInteractive, "Enable interactive mode")
	runCmd.Flags().BoolP("tty", "t", defaults.DefaultDeploymentRunTTY, "Allocate a pseudo-TTY")
	runCmd.Flags().Bool("rm", defaults.DefaultDeploymentRunRemove, "Remove deployment after exit")
	runCmd.Flags().BoolP("detatch", "d", defaults.DefaultDeploymentRunDetatch, "Run in background and detatch")
	runCmd.Flags().StringToStringP("env", "e", nil, "Set environment variables")
	runCmd.Flags().IntSliceP("port", "p", nil, "Expose ports")
	runCmd.Flags().StringP("visibility", "v", defaults.DefaultDeploymentVisibility, "Set deployment visibility (public, private, auth)")
	runCmd.Flags().Float64P("memory", "m", defaults.DefaultDeploymentMemory, "Memory in GB")
	runCmd.Flags().Float64P("cores", "c", defaults.DefaultDeploymentCores, "CPU cores")
	runCmd.Flags().IntP("replicas", "r", defaults.DefaultDeploymentReplicas, "Number of replicas")
	runCmd.Flags().StringP("name", "n", "", "Deployment name")

	rootCmd.AddCommand(runCmd)
}
