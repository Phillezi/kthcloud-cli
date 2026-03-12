package main

import (
	"os"
	"os/signal"

	"github.com/kthcloud/cli/internal/app"
	"github.com/kthcloud/cli/internal/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
		defer cancel()

		if err := app.New(ctx, app.WithKeycloakOptions(
			viper.GetString(constants.ViperKeycloakClientId),
			viper.GetString(constants.ViperKeycloakBaseURL),
			viper.GetString(constants.ViperKeycloakRealm),
		),
			app.WithSessionKey(viper.GetString(constants.ViperSessionKey)),
			app.WithLogger(zap.L()),
		).Login(); err != nil {
			zap.L().Fatal("Error when logging in", zap.Error(err))
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
