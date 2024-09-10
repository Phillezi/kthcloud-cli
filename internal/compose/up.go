package compose

import (
	"errors"
	"go-deploy/dto/v2/body"
	"kthcloud-cli/internal/model"
	"kthcloud-cli/pkg/util"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Up(filename string) error {
	services, err := ParseComposeFile(filename)
	if err != nil {
		log.Errorln(err)
	}

	// load the session from the session.json file
	session, err := model.Load(viper.GetString("session-path"))
	if err != nil {
		log.Fatalln("No active session. Please log in")
	}
	if session.AuthSession.IsExpired() {
		log.Fatalln("Session is expired. Please log in again")
	}
	session.SetupClient()

	projectDir, err := CreateVolume(session, services)
	if err != nil {
		log.Fatalln(err)
	}

	var wg sync.WaitGroup

	for key, service := range services {
		resp, err := session.Client.Req("/v2/deployments", "POST", serviceToDepl(service, key, projectDir))
		if err != nil {
			log.Errorln("error: ", err, " response: ", resp)
			return err
		}
		if resp.IsError() {
			return errors.New("could not create deployment: " + key)
		}
		job, err := util.ProcessResponse[body.DeploymentCreated](resp.String())
		if err != nil {
			return err
		}

		wg.Add(1)
		go func(jobId string, serviceKey string) {
			defer wg.Done()
			log.Infof("Tracking job for deployment: %s", serviceKey)

			for {
				jobResp, err := session.Client.Req("/v2/jobs/"+jobId, "GET", nil)
				if err != nil {
					log.Errorf("Failed to get job status for deployment %s: %v", serviceKey, err)
					return
				}

				jobStatus, err := util.ProcessResponse[body.JobRead](jobResp.String())
				if err != nil {
					log.Errorf("Error processing job status for deployment %s: %v", serviceKey, err)
					return
				}

				//log.Infof("Job status for deployment %s: %v", serviceKey, jobStatus)

				// Break out of loop when job is complete
				if jobStatus.Status == "finished" {
					log.Infof("Deployment %s created successfully", serviceKey)
					break
				}
				if jobStatus.Status == "terminated" {
					log.Infof("Job for deployment %s was terminated", serviceKey)
					break
				}
				if jobStatus.LastError != nil {
					log.Errorf("Failed to create deployment: %s", serviceKey)
					break
				}

				time.Sleep(500 * time.Millisecond)
			}
		}(job.JobID, key)

	}
	wg.Wait()

	log.Info("All jobs have been completed.")
	return nil
}
