package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/kthcloud/cli/internal/app"
	"github.com/kthcloud/cli/internal/constants"
	"github.com/kthcloud/cli/pkg/deploy"
	"github.com/kthcloud/cli/pkg/session"
	"github.com/kthcloud/cli/pkg/ui/renderer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var getDeploymentsCmd = &cobra.Command{
	Use:   "deployment",
	Short: "Get deployments",
	Aliases: []string{
		"deployments",
	},
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

		params := &deploy.GetV2DeploymentsParams{
			All:    new(viper.GetBool("all")),
			Shared: new(viper.GetBool("shared")),
		}
		if userIDFilter := viper.GetString("by-user-id"); userIDFilter != "" {
			params.UserId = &userIDFilter
		}

		rttStart := time.Now()
		r, err := a.Deploy().GetV2DeploymentsWithResponse(ctx, params)
		// Note: not true rtt, it includes some processing overhead too
		rtt := time.Since(rttStart)
		if err != nil {
			if errors.Is(err, session.ErrLoginRequired) {
				zap.L().Fatal("Login is required, please run the login command")
			}
			zap.L().Fatal("Error on request", zap.Error(err))
		}

		obj, err := deploy.HandleAndAssert[*[]deploy.BodyDeploymentRead](r, "get")
		if err != nil {
			zap.L().Fatal("Error on handle", zap.Error(err))
		}

		deployments := make([]renderer.DeploymentLike, len(*obj))
		for i, deployment := range *obj {
			deployments[i] = renderer.DeploymentAdapter{
				ID:     *deployment.Id,
				Name:   *deployment.Name,
				Owner:  *deployment.OwnerId,
				Status: *deployment.Status,
			}
		}

		if err := renderer.New().Render(deployments, renderer.WithOutput(renderer.OutputFromString(viper.GetString("output")))); err != nil {
			zap.L().Fatal("Error on render", zap.Error(err))
		}

		if viper.GetBool("stats") {
			deploymentsCount := 0
			if obj != nil {
				deploymentsCount = len(*obj)
			}
			fmt.Fprintf(os.Stderr, "\n[stats] deployments: %d, rtt: %s, status: %d\n",
				deploymentsCount, rtt, r.StatusCode())
		}
	},
}

func init() {
	getCmd.AddCommand(getDeploymentsCmd)

	getDeploymentsCmd.PersistentFlags().BoolP("shared", "s", false, "Get shared deployments")

	viper.BindPFlags(getDeploymentsCmd.PersistentFlags())
}
