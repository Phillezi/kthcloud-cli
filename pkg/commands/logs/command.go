package logs

import (
	"context"

	"github.com/Phillezi/kthcloud-cli/pkg/deploy"
)

type Command struct {
	ctx    context.Context
	client *deploy.Client

	services []string
}

func New(opts ...CommandOpts) *Command {
	var opt CommandOpts
	if len(opts) > 0 {
		opt = opts[0]
	}
	return &Command{
		ctx:    context.Background(),
		client: opt.Client,

		services: opt.Services,
	}
}

func (c *Command) WithContext(ctx context.Context) *Command {
	c.ctx = ctx
	return c
}
