package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/kthcloud/cli/internal/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var composeConfigCmd = &cobra.Command{
	Use: "config",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
		defer cancel()

		a := app.FromViper(ctx, zap.L(), *viper.GetViper())

		proj, err := a.Compose("", "")
		if err != nil {
			zap.L().Fatal("Err conv", zap.Error(err))
		}

		dat, err := yaml.Marshal(proj)
		if err != nil {
			zap.L().Fatal("Err marshal", zap.Error(err))
		}

		_, _ = fmt.Fprintln(os.Stdout, string(dat))
	},
}

func init() {
	composeCmd.AddCommand(composeConfigCmd)
}
