package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/kthcloud/cli/internal/app"
	"github.com/kthcloud/cli/pkg/deploy"
	"github.com/kthcloud/cli/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [deployment IDs...]",
	Short: "Delete deployments",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()

		a := app.New(ctx,
			app.WithKeycloakOptions(
				viper.GetString("keycloak-client-id"),
				viper.GetString("keycloak-base-url"),
				viper.GetString("keycloak-realm"),
			),
			app.WithSessionKey(viper.GetString("session-key")),
			app.WithLogger(zap.L()),
		)

		var wg sync.WaitGroup
		resultCh := make(chan *deploy.BodyJobRead, len(args))
		errCh := make(chan error, len(args))

		for _, deploymentID := range args {
			wg.Go(func() {
				func(id string) {
					ctxx, cancelx := context.WithTimeout(ctx, 60*time.Second)
					defer cancelx()

					resp, err := a.Deploy().DeleteV2DeploymentsDeploymentId(ctxx, id, deploy.DeleteV2DeploymentsDeploymentIdJSONRequestBody{})
					if err != nil {
						zap.L().Error("failed to delete deployment", zap.String("id", id), zap.Error(err))
						errCh <- err
						resultCh <- nil
						return
					}
					if resp == nil || resp.Body == nil {
						errCh <- fmt.Errorf("nil response for %s", id)
						resultCh <- nil
						return
					}
					defer resp.Body.Close()

					if resp.StatusCode >= 400 {
						errCh <- fmt.Errorf("bad status %d for %s", resp.StatusCode, id)
						resultCh <- nil
						return
					}

					var delJob deploy.BodyDeploymentCreated
					if err := json.NewDecoder(resp.Body).Decode(&delJob); err != nil {
						errCh <- fmt.Errorf("failed to decode delete response for %s: %w", id, err)
						resultCh <- nil
						return
					}

					if delJob.JobId == nil {
						errCh <- fmt.Errorf("delete job ID missing for %s", id)
						resultCh <- nil
						return
					}

					// Poll job until finished
					var jobRead deploy.BodyJobRead
					for {
						jobResp, err := a.Deploy().GetV2JobsJobId(ctxx, *delJob.JobId, deploy.GetV2JobsJobIdJSONRequestBody{})
						if err != nil {
							errCh <- fmt.Errorf("failed to get job %s: %w", *delJob.JobId, err)
							resultCh <- nil
							return
						}
						if jobResp == nil || jobResp.Body == nil {
							errCh <- fmt.Errorf("nil response for job %s", *delJob.JobId)
							resultCh <- nil
							return
						}
						defer jobResp.Body.Close()

						if jobResp.StatusCode >= 400 {
							errCh <- fmt.Errorf("bad status %d for job %s", jobResp.StatusCode, *delJob.JobId)
							resultCh <- nil
							return
						}

						if err := json.NewDecoder(jobResp.Body).Decode(&jobRead); err != nil {
							errCh <- fmt.Errorf("failed to decode job %s: %w", *delJob.JobId, err)
							resultCh <- nil
							return
						}

						if jobRead.Status != nil && (*jobRead.Status == "completed" || *jobRead.Status == "failed") {
							break
						}
						time.Sleep(1 * time.Second) // poll interval
					}

					resultCh <- &jobRead
					zap.L().Info("delete job completed", zap.String("deploymentID", id), zap.String("jobID", utils.DerefOrZero(jobRead.Id)))
				}(deploymentID)
			})
		}

		for r := range resultCh {
			if r == nil {
				continue
			}
			if r.LastError != nil {
				fmt.Printf("Deployment ID: %s, Job ID: %s, Status: %s LastError: %s\n",
					utils.DerefOrZero(r.Id), utils.DerefOrZero(r.Status), utils.DerefOrZero(r.Status), utils.DerefOrZero(r.LastError))
			} else {
				fmt.Printf("Deployment ID: %s, Job ID: %s, Status: %s\n",
					utils.DerefOrZero(r.Id), utils.DerefOrZero(r.Status), utils.DerefOrZero(r.Status))
			}
		}

		wg.Wait()
		close(resultCh)
		close(errCh)

		if len(errCh) > 0 {
			fmt.Fprintf(os.Stderr, "Encountered %d errors during deletion\n", len(errCh))
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
