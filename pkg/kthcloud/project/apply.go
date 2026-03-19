package project

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kthcloud/cli/internal/body"
	"github.com/kthcloud/cli/pkg/deploy"
	"github.com/kthcloud/cli/pkg/scheduler"
)

func (proj *Project) Apply(ctx context.Context, client deploy.ClientWithResponsesInterface) error {
	sched := scheduler.New(ctx)

	if err := scheduler.Batch(sched, proj.Sevices, func(key string, dependent Service) scheduler.Job {
		return &scheduler.FuncJob{
			RunFunc: func(ctx context.Context) error {
				if dependent.Image == nil || *dependent.Image == "" {
					// build & push custom
				}

				r, err := body.Reader(dependent.Body())
				if err != nil {
					return err
				}
				// create / update deployment
				resp, err := client.PostV2DeploymentsWithBodyWithResponse(ctx,
					"application/json",
					r,
				)
				if err != nil {
					return err
				}

				_ = resp

				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(1 * time.Second):
					return nil
				}
			},
			RevertFunc: func(ctx context.Context) error {
				// revent create / update
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(1 * time.Second):
					return nil
				}
			},
		}
	}); err != nil {
		return err
	}

	if err := sched.Start(); err != nil && err != context.Canceled {
		return err
	}

	for stat := range sched.Status() {
		for _, status := range stat {
			fmt.Fprintf(os.Stderr, "%s: %s %v\n", status.ID, status.State, status.Error)
		}
	}

	return nil
}
