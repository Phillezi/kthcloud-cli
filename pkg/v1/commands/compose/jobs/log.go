package jobs

import (
	"fmt"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/scheduler"
	"github.com/briandowns/spinner"
	"github.com/sirupsen/logrus"
)

func MonitorJobStates(jobIDs map[string]string, sched *scheduler.Sched, spinner *spinner.Spinner) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			allDone := true

			var states string

			for name, id := range jobIDs {
				state, err := sched.GetJobState(id)
				if err != nil {
					logrus.Errorf("Failed to get state for job %s: %v", name, err)
					return fmt.Errorf("failed to get job %s state", name)
				}

				//logrus.Infof("%s has state %v", name, state)
				states += fmt.Sprintf("%s has state %v\n", name, state)

				if state == scheduler.Errored {
					logrus.Debugf("job %s is in an ERRORED state\n", name)
					return fmt.Errorf("job %s is in an ERRORED state", name)
				}

				if state != scheduler.Done {
					allDone = false
				}
			}

			if allDone {
				logrus.Infoln("All jobs have completed successfully")
				return nil
			} else {
				spinner.Suffix = states
			}

		}
	}
}
