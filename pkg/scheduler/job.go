package scheduler

import (
	"context"
	"sync"
)

type JobState int64

const (
	Created JobState = iota
	Started
	Done
	Errored
	Cancelling
	Cancelled
)

func (s JobState) String() string {
	switch s {
	case Created:
		return "Created"
	case Started:
		return "Started"
	case Done:
		return "Done"
	case Errored:
		return "Errored"
	case Cancelling:
		return "Cancelling"
	case Cancelled:
		return "Cancelled"
	default:
		return "unknown"
	}
}

// Job represents a unit of work in the scheduler.
type Job struct {
	ID             string
	Dependencies   []*Job
	State          JobState
	Action         func(ctx context.Context, callback func()) error
	CancelCallback func()
	ctx            context.Context
	cancel         context.CancelFunc
	mu             sync.Mutex
}

func NewJob(
	action func(
		ctx context.Context,
		callback func(),
	) error,
	cancelCallback func(),
	dependencies ...*Job,
) *Job {
	if dependencies == nil {
		dependencies = []*Job{}
	}
	return &Job{
		ID:             "",
		Dependencies:   dependencies,
		State:          Created,
		Action:         action,
		CancelCallback: cancelCallback,
	}
}

// After adds dependencies that the job depends on
func (j *Job) After(dependencies ...*Job) {
	if j.Dependencies == nil {
		j.Dependencies = dependencies
	} else {
		j.Dependencies = append(j.Dependencies, dependencies...)
	}
}
