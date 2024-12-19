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
			allCancelled := true
			allCancelledOrErrored := true

			var states string

			for name, id := range jobIDs {
				state, err := sched.GetJobState(id)
				if err != nil {
					logrus.Errorf("Failed to get state for job %s: %v", name, err)
					return fmt.Errorf("failed to get job %s state", name)
				}

				states += fmt.Sprintf("\t%s has state: Job%s\n", name, state.String())
				if state != scheduler.Cancelled {
					allCancelled = false
				}
				if state != scheduler.Cancelled && state != scheduler.Errored {
					allCancelledOrErrored = false
				}
				if state != scheduler.Done {
					allDone = false
				}
			}

			if allDone {
				spinner.FinalMSG = "All jobs have completed successfully\n"
				spinner.Color("green")
				return nil
			} else if allCancelled {
				spinner.FinalMSG = "Cancelled\n"
				spinner.Color("yellow")
				return fmt.Errorf("cancelled")
			} else if allCancelledOrErrored {
				spinner.FinalMSG = "Error\n"
				spinner.Color("red")
				return fmt.Errorf("error occurred")
			} else {
				spinner.Suffix = states
			}

		}
	}
}
