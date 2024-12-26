package progress

import (
	"fmt"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/progress/spinner"
	"github.com/Phillezi/kthcloud-cli/pkg/scheduler"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
)

type Tracker struct {
	Scheduler  *scheduler.Sched
	Rows       []*spinner.Spinner
	multi      pterm.MultiPrinter
	idIndexMap map[string]int
}

func New(scheduler *scheduler.Sched) *Tracker {
	return &Tracker{
		Scheduler:  scheduler,
		multi:      pterm.DefaultMultiPrinter,
		idIndexMap: make(map[string]int),
	}
}

func (t *Tracker) createRows() {
	t.idIndexMap = make(map[string]int)
	jobs := t.Scheduler.GetJobs()
	t.Rows = make([]*spinner.Spinner, len(jobs))
	i := 0
	for id, job := range jobs {
		t.createRow(t.Rows[i], job)
		t.idIndexMap[id] = i
		i++
	}
}

func (t *Tracker) createRow(row *spinner.Spinner, job *scheduler.Job) {
	switch job.State {
	case scheduler.Errored:
		row = spinner.New(t.multi.NewWriter(), scheduler.Errored, pterm.FgWhite, pterm.BgRed)
	case scheduler.Cancelled:
		row = spinner.New(t.multi.NewWriter(), scheduler.Cancelled, pterm.FgGray)
	case scheduler.Cancelling:
		row = spinner.New(t.multi.NewWriter(), scheduler.Cancelling, pterm.FgYellow)
	case scheduler.Done:
		row = spinner.New(t.multi.NewWriter(), scheduler.Done, pterm.FgGreen)
	case scheduler.Created:
		row = spinner.New(t.multi.NewWriter(), scheduler.Created, pterm.FgBlue)
	default:
		row = spinner.New(t.multi.NewWriter(), scheduler.Started, pterm.FgBlue)
	}

	row.Printer.Start(fmt.Sprintf("%s - Job%s", job.DisplayName, job.State.String()))
	if job.Start != nil {
		row.Printer.SetStartedAt(*job.Start)
	}
}

func (t *Tracker) renderRow(row *spinner.Spinner, job *scheduler.Job) {
	if row == nil || job == nil {
		return
	}
	if row.PrevState != job.State {
		t.createRow(row, job)
		return
	}

	switch job.State {
	case scheduler.Errored:
		row.Printer.UpdateText(fmt.Sprintf("%s - Job%s", job.DisplayName, job.State.String()))
	case scheduler.Cancelled:
		row.Printer.UpdateText(fmt.Sprintf("%s - Job%s", job.DisplayName, job.State.String()))
	case scheduler.Cancelling:
		row.Printer.UpdateText(fmt.Sprintf("%s - Job%s", job.DisplayName, job.State.String()))
	case scheduler.Done:
		row.Printer.UpdateText(fmt.Sprintf("%s - Job%s", job.DisplayName, job.State.String()))
	case scheduler.Created:
		row.Printer.UpdateText(fmt.Sprintf("%s - Job%s", job.DisplayName, job.State.String()))
	default:
		row.Printer.UpdateText(fmt.Sprintf("%s - Job%s", job.DisplayName, job.State.String()))
	}
}

func (t *Tracker) Render() {
	if t.Scheduler.NumJobs() != len(t.idIndexMap) {
		t.createRows()
		return
	}

	for id, job := range t.Scheduler.GetJobs() {
		t.renderRow(t.Rows[t.idIndexMap[id]], job)
	}

}

func (t *Tracker) TrackJobs() error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	defer t.multi.Stop()

	t.multi.Start()
	//logrus.AddHook(NewTrackerHook(t))
	logrus.SetOutput(t.multi.NewWriter())
	defer logrus.SetOutput(logrus.StandardLogger().Out)

	for {
		select {
		case <-ticker.C:
			allDone := true
			allCancelled := true
			allCancelledOrErrored := true

			if t.Scheduler.NumJobs() != len(t.idIndexMap) {
				t.createRows()
				continue
			}

			for id, job := range t.Scheduler.GetJobs() {
				t.renderRow(t.Rows[t.idIndexMap[id]], job)
				state := job.State
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
				logrus.Println("All jobs have completed successfully")
				return nil
			} else if allCancelled {
				logrus.Println("cancelled")
				return fmt.Errorf("cancelled")
			} else if allCancelledOrErrored {
				logrus.Println("Error")
				return fmt.Errorf("error occurred")
			}

		}
	}

}
