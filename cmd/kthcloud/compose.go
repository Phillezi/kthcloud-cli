package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/kthcloud/cli/internal/app"
	"github.com/kthcloud/cli/internal/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var composeCmd = &cobra.Command{
	Use: "compose",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
		defer cancel()

		a := app.New(ctx, app.WithKeycloakOptions(
			viper.GetString(constants.ViperKeycloakClientId),
			viper.GetString(constants.ViperKeycloakBaseURL),
			viper.GetString(constants.ViperKeycloakRealm),
		),
			app.WithSessionKey(viper.GetString(constants.ViperSessionKey)),
			app.WithAPITokenSession(viper.GetString(constants.ViperDeployAPIToken)),
			app.WithLogger(zap.L()),
		)

		proj, err := a.Compose("", "")
		if err != nil {
			zap.L().Fatal("Err conv", zap.Error(err))
		}

		dat, err := yaml.Marshal(proj)
		if err != nil {
			zap.L().Fatal("Err marshal", zap.Error(err))
		}

		fmt.Fprintln(os.Stdout, string(dat))
	},
}

func init() {
	rootCmd.AddCommand(composeCmd)
}
