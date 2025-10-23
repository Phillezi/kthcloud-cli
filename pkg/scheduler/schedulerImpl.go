package scheduler

import (
	"context"

	"github.com/heimdalr/dag"
)

type SchedulerImpl struct {
	ctx      context.Context
	dag      *dag.DAG
	executor *Executor
}

func New(ctx context.Context, opts ...Option) *SchedulerImpl {
	si := &SchedulerImpl{
		ctx: ctx,
		dag: dag.NewDAG(),
	}

	for _, opt := range opts {
		opt(si)
	}

	return si
}
