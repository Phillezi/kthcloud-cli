package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/kthcloud/cli/internal/app"
	"github.com/kthcloud/cli/pkg/deploy"
	"github.com/kthcloud/cli/pkg/session"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var deleteDeploymentCmd = &cobra.Command{
	Use: "deployment [uuid...]",
	Aliases: []string{
		"deployments",
	},
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			return
		}

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

		var wg sync.WaitGroup
		errCh := make(chan string, len(args))

		for _, uuid := range args {
			wg.Go(func() {
				func(uuid string) {
					r, err := a.Deploy().DeleteV2DeploymentsDeploymentIdWithResponse(ctx, uuid, deploy.DeleteV2DeploymentsDeploymentIdJSONBody{})
					if err != nil {
						if errors.Is(err, session.ErrLoginRequired) {
							zap.L().Fatal("Login is required, please run the login command")
						}
						errCh <- fmt.Sprintf("Delete request failed for %s: %v", uuid, err)
						return
					}

					deleteJobResp, err := deploy.HandleAndAssert[*deploy.BodyDeploymentCreated](r, "delete")
					if err != nil {
						errCh <- fmt.Sprintf("Failed handling delete response for %s: %v", uuid, err)
						return
					}

					for {
						resp, err := a.Deploy().GetV2JobsJobIdWithResponse(ctx, *deleteJobResp.JobId, deploy.GetV2JobsJobIdJSONBody{})
						if err != nil {
							errCh <- fmt.Sprintf("Error getting job %s: %v", uuid, err)
							break
						}

						job, err := deploy.HandleAndAssert[*deploy.BodyJobRead](resp, "getJob")
						if err != nil {
							errCh <- fmt.Sprintf("Error handling job %s: %v", uuid, err)
							break
						}

						if job.Status != nil && *job.Status == "finished" {
							fmt.Printf("%s\n", uuid)
							break
						}

						if job.LastError != nil && *job.LastError != "" {
							errCh <- fmt.Sprintf("Job error for %s: %s", uuid, *job.LastError)
							break
						}
					}
				}(uuid)
			})
		}

		wg.Wait()
		close(errCh)

		for e := range errCh {
			zap.L().Error(e)
		}

	},
}

func init() {
	deleteCmd.AddCommand(deleteDeploymentCmd)
}
