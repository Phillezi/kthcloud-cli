package jobs

import (
	"context"
	"fmt"
	"go-deploy/dto/v2/body"
	"strings"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/rand"
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

func (m *DeploymentJob) MockTrack(ctx context.Context, deploymentName string, tickerInterval time.Duration, onCancel func()) error {
	// Randomize the sleep duration to simulate random waiting between job checks
	waitTime := time.Duration(rand.Intn(5)+1) * time.Second
	time.Sleep(waitTime)

	// Simulate the process of tracking the deployment job status
	logrus.Debugf("Started tracking deployment: %s", deploymentName)
	ticker := time.NewTicker(tickerInterval)
	ticker2 := time.NewTicker(700 * time.Millisecond)
	defer ticker.Stop()
	defer ticker2.Stop()

	var status string
	var i int

	for {
		select {
		case <-ctx.Done():
			logrus.Debugf("Deployment %s was cancelled\n", deploymentName)
			if onCancel != nil {
				onCancel()
			}
			return nil
		case <-ticker2.C:
			randomNumber := rand.Intn(100)
			if randomNumber < 2 {
				status = "errored"
			} else {
				status = []string{"in-progress", "in-progress", "finished"}[i]
				i++
			}
		case <-ticker.C:
			jobStatus := &body.JobRead{
				Status: status,
			}

			// Process the job status and handle it
			switch jobStatus.Status {
			case "finished":
				logrus.Debugf("Deployment %s completed successfully", deploymentName)
				return nil
			case "terminated":
				logrus.Debugf("Job for deployment %s was terminated", deploymentName)
				return nil
			case "errored":
				// Simulate an error during the job process
				logrus.Debugf("Deployment %s encountered an error", deploymentName)
				return fmt.Errorf("failed to track deployment %s due to an error", deploymentName)
			}

			// If the job is still in progress, log the status
			if jobStatus.Status == "in-progress" {
				logrus.Debugf("Deployment %s is still in progress", deploymentName)
			}
		}
	}
}
