package jobs

import (
	"go-deploy/dto/v2/body"
	"sync"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	"github.com/briandowns/spinner"
	log "github.com/sirupsen/logrus"
)

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