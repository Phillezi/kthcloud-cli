package scheduler

import (
	"context"
	"fmt"
	"time"
)

type Job interface {
	Run(ctx context.Context) error
	Revert(ctx context.Context) error
}

type JobImpl struct {
	Value string
	Fail  bool
}

func (v *JobImpl) Run(ctx context.Context) error {
	if v.Fail {
		return fmt.Errorf("i wanted to fail")
	}
	select {
	case <-time.After(1 * time.Second):
	case <-ctx.Done():
		return ctx.Err()
	}
	fmt.Println(v.Value)
	return nil
}

func (v *JobImpl) Revert(ctx context.Context) error {
	select {
	case <-time.After(1 * time.Second):
	case <-ctx.Done():
		return ctx.Err()
	}
	fmt.Println("revert: " + v.Value)
	return nil
}

type Node struct {
	Task Job
	Next *Node
}

type Scheduler struct {
}
