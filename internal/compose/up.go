package compose

import (
	"fmt"
	"go-deploy/dto/v2/body"
	"strings"
	"sync"
	"time"

	"github.com/Phillezi/kthcloud-cli/internal/model"
	"github.com/Phillezi/kthcloud-cli/pkg/util"

	"github.com/briandowns/spinner"
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

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Color("blue")
	s.Start()
	defer s.Stop()

	for key, service := range services {
		resp, err := session.Client.Req("/v2/deployments", "POST", serviceToDepl(service, key, projectDir))
		if err != nil {
			log.Errorln("error: ", err, " response: ", resp)
			return err
		}
		if strings.HasPrefix(resp.String(), "{\"errors\":") {
			errors, err := util.ProcessResponse[model.ErrorResponse](resp.String())
			if err != nil {
				return err
			}
			s.Color("red")
			// TODO: handle race condition here later
			s.Stop()
			log.Errorf("Error when trying to create deployment %s: %v", key, *errors)
			s.Start()

			return fmt.Errorf("error when trying to create deployment %s: %v", key, *errors)
		}
		job, err := util.ProcessResponse[body.DeploymentCreated](resp.String())
		if err != nil {
			return err
		}

		wg.Add(1)
		go func(jobId string, serviceKey string) {
			defer wg.Done()

			for {
				jobResp, err := session.Client.Req("/v2/jobs/"+jobId, "GET", nil)
				if err != nil {
					s.Color("red")
					log.Errorf("Failed to get job status for deployment %s: %v", serviceKey, err)
					return
				}

				jobStatus, err := util.ProcessResponse[body.JobRead](jobResp.String())
				if err != nil {
					s.Color("red")
					log.Errorf("Error processing job status for deployment %s: %v", serviceKey, err)
					return
				}

				// Break out of loop when job is complete
				if jobStatus.Status == "finished" {
					s.Color("green")
					// TODO: handle race condition here later
					s.Stop()
					log.Infof("Deployment %s created successfully", serviceKey)
					s.Start()
					break
				}
				if jobStatus.Status == "terminated" {
					s.Color("yellow")
					// TODO: handle race condition here later
					s.Stop()
					log.Infof("Job for deployment %s was terminated", serviceKey)
					s.Start()
					break
				}
				if jobStatus.LastError != nil {
					s.Color("red")
					// TODO: handle race condition here later
					s.Stop()
					log.Errorf("Failed to create deployment: %s", serviceKey)
					s.Start()
					break
				}

				time.Sleep(500 * time.Millisecond)
			}
		}(job.JobID, key)

	}
	wg.Wait()
	s.Color("green")
	s.Stop()

	log.Info("All jobs have been completed.")
	return nil
}
