package check

import (
	"context"

	"github.com/Phillezi/kthcloud-cli/pkg/deploy"
)

type Command struct {
	ctx    context.Context
	client *deploy.Client
}

func New(opts ...CommandOpts) *Command {
	var opt CommandOpts
	if len(opts) > 0 {
		opt = opts[0]
	}
	return &Command{
		ctx:    context.Background(),
		client: opt.Client,
	}
}

func (c *Command) WithContext(ctx context.Context) *Command {
	c.ctx = ctx
	return c
}
