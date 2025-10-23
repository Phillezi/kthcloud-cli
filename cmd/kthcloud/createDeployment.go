package main

import (
	"os"
	"os/signal"

	"github.com/kthcloud/cli/internal/app"
	"github.com/kthcloud/cli/pkg/deploy"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var createDeploymentCmd = &cobra.Command{
	Use: "deployment",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
		defer cancel()

		a := app.New(ctx, app.WithKeycloakOptions(
			viper.GetString("keycloak-client-id"),
			viper.GetString("keycloak-base-url"),
			viper.GetString("keycloak-realm"),
		),
			app.WithSessionKey(viper.GetString("session-key")),
			app.WithLogger(zap.L()),
		)

		a.Deploy().PostV2DeploymentsWithResponse(ctx, deploy.PostV2DeploymentsJSONRequestBody{
			
		})

	},
}

func init() {
	createCmd.AddCommand(createDeploymentCmd)
}
