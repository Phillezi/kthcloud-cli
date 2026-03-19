package main

import (
	"os"
	"os/signal"

	"github.com/kthcloud/cli/internal/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var composeDownCmd = &cobra.Command{
	Use: "down",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
		defer cancel()

		a := app.FromViper(ctx, zap.L(), *viper.GetViper())

		proj, err := a.Compose(viper.GetStringSlice("compose.file"), viper.GetString("compose.projectname"))
		if err != nil {
			zap.L().Fatal("Err conv", zap.Error(err))
		}

		if err := proj.Destroy(ctx, a.Deploy()); err != nil {
			zap.L().Fatal("Err destroy", zap.Error(err))
		}
	},
}

func init() {
	composeCmd.AddCommand(composeDownCmd)
}
