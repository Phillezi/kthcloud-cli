package jobs

import (
	"context"
	"fmt"
	"go-deploy/dto/v2/body"
	"sync"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	"github.com/briandowns/spinner"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func Track(ctx context.Context, deploymentName string, job *body.DeploymentCreated, tickerInterval time.Duration, onCancel func()) error {
	c := client.Get().Client()
	ticker := time.NewTicker(tickerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logrus.Debugf("deployment %s was cancelled\n", deploymentName)
			onCancel()
			return nil
		case <-ticker.C:
			jobResp, err := c.R().Get("/v2/jobs/" + job.JobID)
			if err != nil {
				return fmt.Errorf("failed to get job status for deployment %s: %w", deploymentName, err)
			}

			jobStatus, err := util.ProcessResponse[body.JobRead](jobResp.String())
			if err != nil {
				return fmt.Errorf("error processing job status for deployment %s: %w", deploymentName, err)
			}

			switch jobStatus.Status {
			case "finished":
				log.Debugf("Deployment %s created successfully", deploymentName)
				return nil
			case "terminated":
				log.Debugf("Job for deployment %s was terminated", deploymentName)
				return nil
			}

			if jobStatus.LastError != nil {
				return fmt.Errorf("failed to create deployment %s: %s", deploymentName, *jobStatus.LastError)
			}
		}
	}
}

func TrackDel(deploymentName string, job *body.DeploymentDeleted, tickerInterval time.Duration) error {
	c := client.Get().Client()
	ticker := time.NewTicker(tickerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			jobResp, err := c.R().Get("/v2/jobs/" + job.JobID)
			if err != nil {
				return fmt.Errorf("failed to get job status for deployment %s: %w", deploymentName, err)
			}

			jobStatus, err := util.ProcessResponse[body.JobRead](jobResp.String())
			if err != nil {
				return fmt.Errorf("error processing job status for deployment %s: %w", deploymentName, err)
			}

			switch jobStatus.Status {
			case "finished":
				log.Debugf("Deployment %s deleted successfully", deploymentName)
				return nil
			case "terminated":
				log.Debugf("Job for deployment %s was terminated", deploymentName)
				return nil
			}

			if jobStatus.LastError != nil {
				return fmt.Errorf("failed to delete deployment %s: %s", deploymentName, *jobStatus.LastError)
			}
		}
	}
}

func TrackDelW(ctx context.Context, deploymentName string, job *body.DeploymentDeleted, tickerInterval time.Duration, s *spinner.Spinner, onCancel func()) error {
	c := client.Get().Client()
	ticker := time.NewTicker(tickerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logrus.Debugf("deployment %s was cancelled\n", deploymentName)
			onCancel()
			return nil
		case <-ticker.C:
			jobResp, err := c.R().Get("/v2/jobs/" + job.JobID)
			if err != nil {
				s.Color("red")
				s.Stop()
				return fmt.Errorf("failed to get job status for deployment %s: %w", deploymentName, err)
			}

			jobStatus, err := util.ProcessResponse[body.JobRead](jobResp.String())
			if err != nil {
				s.Color("red")
				s.Stop()
				return fmt.Errorf("error processing job status for deployment %s: %w", deploymentName, err)
			}

			switch jobStatus.Status {
			case "finished":
				s.Color("green")
				s.Stop()
				log.Infof("Deployment %s deleted successfully", deploymentName)
				return nil
			case "terminated":
				s.Color("yellow")
				s.Stop()
				log.Infof("Job for deployment %s was terminated", deploymentName)
				return nil
			}

			if jobStatus.LastError != nil {
				s.Color("red")
				s.Stop()
				return fmt.Errorf("failed to delete deployment %s: %s", deploymentName, *jobStatus.LastError)
			}
		}
	}
}

func TrackDeploymentCreation(deploymentName string, job *body.DeploymentCreated, wg *sync.WaitGroup, s *spinner.Spinner) {
	c := client.Get().Client()

	defer wg.Done()
	for {
		jobResp, err := c.R().Get("/v2/jobs/" + job.JobID)
		if err != nil {
			s.Color("red")
			//s.Lock()
			s.Stop()
			log.Errorf("Failed to get job status for deployment %s: %v", deploymentName, err)
			s.Start()
			//s.Unlock()
			return
		}

		jobStatus, err := util.ProcessResponse[body.JobRead](jobResp.String())
		if err != nil {
			s.Color("red")
			//s.Lock()
			s.Stop()
			log.Errorf("Error processing job status for deployment %s: %v", deploymentName, err)
			s.Start()
			//s.Unlock()
			return
		}

		// Break out of loop when job is complete
		if jobStatus.Status == "finished" {
			s.Color("green")
			//s.Lock()
			s.Stop()
			log.Infof("Deployment %s created successfully", deploymentName)
			s.Start()
			//s.Unlock()
			break
		}
		if jobStatus.Status == "terminated" {
			s.Color("yellow")
			//s.Lock()
			s.Stop()
			log.Infof("Job for deployment %s was terminated", deploymentName)
			s.Start()
			//s.Unlock()
			break
		}
		if jobStatus.LastError != nil {
			s.Color("red")
			//s.Lock()
			s.Stop()
			log.Errorf("Failed to create deployment: %s", deploymentName)
			s.Start()
			//s.Unlock()
			break
		}

		time.Sleep(500 * time.Millisecond)
	}

}

func TrackDeploymentCreationW(deploymentName string, job *body.DeploymentCreated, wg *sync.WaitGroup, s *spinner.Spinner) {
	wg.Add(1)
	go TrackDeploymentCreation(deploymentName, job, wg, s)
}

func TrackDeploymentDeletion(deploymentName string, job *body.DeploymentDeleted, wg *sync.WaitGroup, s *spinner.Spinner) {
	c := client.Get().Client()

	defer wg.Done()
	for {
		jobResp, err := c.R().Get("/v2/jobs/" + job.JobID)
		if err != nil {
			s.Color("red")
			//s.Lock()
			s.Stop()
			log.Errorf("Failed to get job status for deployment %s: %v", deploymentName, err)
			s.Start()
			//s.Unlock()
			return
		}

		jobStatus, err := util.ProcessResponse[body.JobRead](jobResp.String())
		if err != nil {
			s.Color("red")
			//s.Lock()
			s.Stop()
			log.Errorf("Error processing job status for deployment %s: %v", deploymentName, err)
			s.Start()
			//s.Unlock()
			return
		}

		// Break out of loop when job is complete
		if jobStatus.Status == "finished" {
			s.Color("green")
			//s.Lock()
			s.Stop()
			log.Infof("Deployment %s deleted successfully", deploymentName)
			s.Start()
			//s.Unlock()
			break
		}
		if jobStatus.Status == "terminated" {
			s.Color("yellow")
			//s.Lock()
			s.Stop()
			log.Infof("Job for deployment %s was terminated", deploymentName)
			s.Start()
			//s.Unlock()
			break
		}
		if jobStatus.LastError != nil {
			s.Color("red")
			//s.Lock()
			s.Stop()
			log.Errorf("Failed to delete deployment: %s", deploymentName)
			s.Start()
			//s.Unlock()
			break
		}

		time.Sleep(500 * time.Millisecond)
	}

}

func TrackDeploymentDeletionW(deploymentName string, job *body.DeploymentDeleted, wg *sync.WaitGroup, s *spinner.Spinner) {
	wg.Add(1)
	go TrackDeploymentDeletion(deploymentName, job, wg, s)
}
