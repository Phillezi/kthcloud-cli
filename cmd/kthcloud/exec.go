package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"github.com/kthcloud/cli/internal/constants"
	"github.com/kthcloud/cli/internal/defaults"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var execCmd = &cobra.Command{
	Use:  "exec [deployment] <args...>",
	Args: cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			return
		}

		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
		defer cancel()

		args[0] = fmt.Sprintf("%s@%s", args[0], viper.GetString(constants.ViperSSHHost))

		ssh := exec.CommandContext(ctx, viper.GetString(constants.ViperSSHBinary), args...)

		ssh.Stdin = os.Stdin
		ssh.Stdout = os.Stdout
		ssh.Stderr = os.Stderr

		err := ssh.Run()
		if err != nil {
			if ctx.Err() == context.Canceled {
				return
			}

			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}

			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)

	rootCmd.Flags().String("ssh-host", defaults.DefaultSSHHost, "ssh host to use for exec")
	viper.BindPFlag(constants.ViperSSHHost, rootCmd.Flags().Lookup("ssh-host"))

	rootCmd.Flags().String("ssh-binary", defaults.DefaultSSHBinary, "ssh binary to use for exec")
	viper.BindPFlag(constants.ViperSSHBinary, rootCmd.Flags().Lookup("ssh-binary"))
}
