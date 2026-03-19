package main

import (
	"os"
	"os/signal"

	"github.com/kthcloud/cli/internal/app"
	"github.com/kthcloud/cli/pkg/out"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var composeConfigCmd = &cobra.Command{
	Use: "config",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
		defer cancel()

		a := app.FromViper(ctx, zap.L(), *viper.GetViper())

		proj, err := a.Compose(viper.GetStringSlice("compose.file"), viper.GetString("compose.projectname"))
		if err != nil {
			zap.L().Fatal("Err conv", zap.Error(err))
		}
		if err := out.Write(proj, viper.GetString("compose.format"), os.Stdout); err != nil {
			zap.L().Fatal("Err marshal", zap.Error(err))
		}
	},
}

func init() {
	composeCmd.AddCommand(composeConfigCmd)

	composeConfigCmd.Flags().String("format", out.FormatYAML, "Output format")
	viper.BindPFlag("compose.format", composeConfigCmd.Flags().Lookup("format"))
}
