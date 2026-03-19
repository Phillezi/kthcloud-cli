package main

import (
	"context"
	"errors"
	"os"
	"os/signal"

	"github.com/kthcloud/cli/internal/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var composeUpCmd = &cobra.Command{
	Use: "up",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
		defer cancel()

		a := app.FromViper(ctx, zap.L(), *viper.GetViper())

		proj, err := a.Compose(viper.GetStringSlice("compose.file"), viper.GetString("compose.projectname"))
		if err != nil {
			zap.L().Fatal("Err conv", zap.Error(err))
		}

		if err := proj.Build(ctx, a.Deploy(), viper.GetStringSlice("compose.build")); err != nil {
			zap.L().Fatal("Error building", zap.Error(err))
		}

		if err := proj.Apply(ctx, a.Deploy()); err != nil {
			zap.L().Fatal("Error applying", zap.Error(err))
		}

		if !viper.GetBool("compose.detach") {
			if err := proj.Logs(ctx); err != nil && !errors.Is(err, context.Canceled) {
				zap.L().Fatal("Error getting logs", zap.Error(err))
			}
		}
	},
}

func init() {
	composeCmd.AddCommand(composeUpCmd)

	composeUpCmd.Flags().BoolP("detach", "d", false, "Detach after creation")
	viper.BindPFlag("compose.detach", composeUpCmd.Flags().Lookup("detach"))

	composeUpCmd.Flags().StringSlice("build", nil, "Services to build none means all")
	viper.BindPFlag("compose.build", composeUpCmd.Flags().Lookup("build"))
}
