package scheduler

import "context"

type JobState string

const (
	JobPending  JobState = "pending"
	JobRunning  JobState = "running"
	JobDone     JobState = "done"
	JobFailed   JobState = "failed"
	JobReverted JobState = "reverted"
)

type JobStatus struct {
	ID    string
	State JobState
	Error error
}

type Job interface {
	Run(ctx context.Context) error
	Revert(ctx context.Context) error
}

type Scheduler interface {
	Start() error
	Add(j Job, deps ...string) (id string, err error)
	Status() <-chan []JobStatus
}
