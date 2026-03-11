package deploy

import (
	"context"
	"fmt"
	"time"

	"github.com/kthcloud/go-deploy/dto/v2/body"
)

type JobType uint

const (
	JobTypeCreateDeployment JobType = iota
)

type Waiter struct {
	deploy ClientWithResponsesInterface
}

func (w *Waiter) Wait(ctx context.Context, id string, statusUpdates chan string) error {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		stat, err := w.deploy.GetV2JobsJobIdWithResponse(ctx, id, nil)
		if err != nil {
			return err
		}

		job, err := HandleAndAssert[body.JobRead](stat)
		if err != nil {
			return err
		}

		if job.FinishedAt != nil {
			if job.LastError != nil {
				return fmt.Errorf("%s", *job.LastError)
			}
			return nil
		}
		if statusUpdates != nil {
			statusUpdates <- job.Status
		}
	}
	return ctx.Err()
}
