package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/kthcloud/cli/internal/app"
	"github.com/kthcloud/cli/pkg/deploy"
	"github.com/kthcloud/cli/pkg/session"
	"github.com/kthcloud/cli/pkg/ui/renderer"
	"github.com/kthcloud/cli/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var getJobsCmd = &cobra.Command{
	Use:   "job",
	Short: "Get job",
	Aliases: []string{
		"jobs",
	},
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

		params := &deploy.GetV2JobsParams{
			All: utils.PtrOf(viper.GetBool("all")),
		}
		if userIDFilter := viper.GetString("by-user-id"); userIDFilter != "" {
			params.UserId = &userIDFilter
		}
		if typeFilter := viper.GetString("type"); typeFilter != "" {
			params.Type = &typeFilter
		}
		if statusFilter := viper.GetString("status"); statusFilter != "" {
			params.Status = &statusFilter
		}

		rttStart := time.Now()
		r, err := a.Deploy().GetV2JobsWithResponse(ctx, params, deploy.GetV2JobsJSONRequestBody{})
		// Note: not true rtt, it includes some processing overhead too
		rtt := time.Since(rttStart)
		if err != nil {
			if errors.Is(err, session.ErrLoginRequired) {
				zap.L().Fatal("Login is required, please run the login command")
			}
			zap.L().Fatal("Error on request", zap.Error(err))
		}

		obj, err := deploy.HandleAndAssert[*[]deploy.BodyJobRead](r, "get")
		if err != nil {
			zap.L().Fatal("Error on handle", zap.Error(err))
		}

		if err := renderer.New().Render(obj, renderer.WithOutput(renderer.OutputFromString(viper.GetString("output")))); err != nil {
			zap.L().Fatal("Error on render", zap.Error(err))
		}

		if viper.GetBool("stats") {
			deploymentsCount := 0
			if obj != nil {
				deploymentsCount = len(*obj)
			}
			fmt.Fprintf(os.Stderr, "\n[stats] jobs: %d, rtt: %s, status: %d\n",
				deploymentsCount, rtt, r.StatusCode())
		}
	},
}

func init() {
	getCmd.AddCommand(getJobsCmd)

	getJobsCmd.PersistentFlags().StringP("type", "t", "", "Get by type")
	getJobsCmd.PersistentFlags().StringP("status", "s", "", "Get by status")

	viper.BindPFlags(getJobsCmd.PersistentFlags())
}
