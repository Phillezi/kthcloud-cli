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
	Cancelled
)

// Job represents a unit of work in the scheduler.
type Job struct {
	ID           string
	Dependencies []*Job
	State        JobState
	Action       func(ctx context.Context, callback func(cArg interface{})) error
	Callback     func(cArg interface{})
	ctx          context.Context
	cancel       context.CancelFunc
	mu           sync.Mutex
}

func NewJob(
	action func(
		ctx context.Context,
		callback func(cArg interface{}),
	) error,
	callback func(cArg interface{}),
	dependencies ...*Job,
) *Job {
	if dependencies == nil {
		dependencies = []*Job{}
	}
	return &Job{
		ID:           "",
		Dependencies: dependencies,
		State:        Created,
		Action:       action,
		Callback:     callback,
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
