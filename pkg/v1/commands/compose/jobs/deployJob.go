package jobs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kthcloud/go-deploy/dto/v2/body"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	"github.com/sirupsen/logrus"
)

type DeploymentJob struct {
	Type  string
	ID    string  `json:"id"`
	JobID *string `json:"jobId,omitempty"`
}

func From(job interface{}) *DeploymentJob {
	switch v := job.(type) {
	case *body.DeploymentCreated:
		return &DeploymentJob{
			Type:  "Created",
			ID:    v.ID,
			JobID: &v.JobID,
		}
	case *body.DeploymentUpdated:
		return &DeploymentJob{
			Type:  "Updated",
			ID:    v.ID,
			JobID: v.JobID,
		}
	case *body.DeploymentDeleted:
		return &DeploymentJob{
			Type:  "Deleted",
			ID:    v.ID,
			JobID: &v.JobID,
		}
	default:
		logrus.Debugln("cant create DeploymentJob from this type")
		return nil
	}
}

func FromCreated(job *body.DeploymentCreated) *DeploymentJob {
	return &DeploymentJob{
		Type:  "Created",
		ID:    job.ID,
		JobID: &job.JobID,
	}
}

func FromUpdated(job *body.DeploymentUpdated) *DeploymentJob {
	return &DeploymentJob{
		Type:  "Updated",
		ID:    job.ID,
		JobID: job.JobID,
	}
}

func FromDeleted(job *body.DeploymentDeleted) *DeploymentJob {
	return &DeploymentJob{
		Type:  "Deleted",
		ID:    job.ID,
		JobID: &job.JobID,
	}
}

func (job *DeploymentJob) Track(ctx context.Context, deploymentName string, tickerInterval time.Duration, onCancel func()) error {
	if job.JobID == nil {
		logrus.Debugln("no jobID to track")
	}

	c := client.Get().Client()
	ticker := time.NewTicker(tickerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logrus.Debugf("deployment %s was cancelled\n", deploymentName)
			if onCancel != nil {
				onCancel()
			}
			return nil
		case <-ticker.C:
			jobResp, err := c.R().Get("/v2/jobs/" + *job.JobID)
			if err != nil {
				return fmt.Errorf("failed to get job status for deployment %s: %w", deploymentName, err)
			}

			jobStatus, err := util.ProcessResponse[body.JobRead](jobResp.String())
			if err != nil {
				return fmt.Errorf("error processing job status for deployment %s: %w", deploymentName, err)
			}

			switch jobStatus.Status {
			case "finished":
				logrus.Debugf("Deployment %s %s successfully", deploymentName, strings.ToLower(job.Type))
				return nil
			case "terminated":
				logrus.Debugf("Job for deployment %s was terminated", deploymentName)
				return nil
			}

			if jobStatus.LastError != nil {
				return fmt.Errorf("failed to %s deployment %s: %s", strings.ToLower(strings.TrimSuffix(job.Type, "d")), deploymentName, *jobStatus.LastError)
			}
		}
	}
}
