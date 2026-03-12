package scheduler

import "context"

type FuncJob struct {
	RunFunc    func(ctx context.Context) error
	RevertFunc func(ctx context.Context) error
}

func (j *FuncJob) Run(ctx context.Context) error {
	if j.RunFunc == nil {
		return nil
	}
	return j.RunFunc(ctx)
}

func (j *FuncJob) Revert(ctx context.Context) error {
	if j.RevertFunc == nil {
		return nil
	}
	return j.RevertFunc(ctx)
}
