package main

import (
	"os"
	"os/signal"

	"github.com/google/uuid"
	"github.com/kthcloud/cli/internal/app"
	"github.com/kthcloud/cli/internal/constants"
	"github.com/kthcloud/cli/internal/defaults"
	"github.com/kthcloud/cli/pkg/deploy"
	"github.com/kthcloud/cli/pkg/logs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var logCmd = &cobra.Command{
	Use: "log",
	Aliases: []string{
		"logs",
	},
	Args: cobra.MinimumNArgs(1),
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

		var seen map[string]struct{} = make(map[string]struct{}, len(args))

		var deployments []*deploy.BodyDeploymentRead = make([]*deploy.BodyDeploymentRead, 0, len(args))
		for _, idOrName := range args {
			if _, found := seen[idOrName]; found {
				zap.L().Warn("Duplicate entry, skipping...", zap.String("id", idOrName))
				continue
			}
			seen[idOrName] = struct{}{}
			_, err := uuid.Parse(idOrName)
			if err != nil {
				zap.L().Warn("Invalid id, skipping...", zap.Error(err))
				continue
			}
			r, err := a.Deploy().GetV2DeploymentsDeploymentIdWithResponse(ctx, idOrName, a.SessionMiddleware())
			if err != nil {
				zap.L().Warn("Error getting deployment, skipping...", zap.Error(err))
				continue
			}
			d, err := deploy.HandleAndAssert[*deploy.BodyDeploymentRead](r, "get")
			if err != nil {
				zap.L().Error("Error on handle", zap.Error(err))
				continue
			}
			deployments = append(deployments, d)
		}

		l := logs.New(
			logs.WithContext(ctx),
			logs.WithAPIURL(defaults.DefaultDeployAPIBaseURL),
			logs.WithSession(a.Session()),
			logs.WithLogger(zap.L().Named("logs")),
		)

		if err := l.Subscribe(deployments...); err != nil {
			zap.L().Error("Error subscribing", zap.Error(err))
		}

		if err := l.Consume(os.Stderr); err != nil {
			zap.L().Error("Error consuming logs", zap.Error(err))
		}

	},
}

func init() {
	rootCmd.AddCommand(logCmd)

	//FIME: not used yet
	logCmd.PersistentFlags().BoolP("follow", "f", true, "Follow")

	viper.BindPFlags(logCmd.PersistentFlags())
}
