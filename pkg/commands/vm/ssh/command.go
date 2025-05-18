package ssh

import (
	"context"

	"github.com/Phillezi/kthcloud-cli/pkg/deploy"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
)

type Command struct {
	ctx    context.Context
	client *deploy.Client

	id   string
	name string
}

func New(opts ...CommandOpts) *Command {
	var opt CommandOpts
	if len(opts) > 0 {
		opt = opts[0]
	}
	return &Command{
		ctx:    context.Background(),
		client: opt.Client,

		id:   util.PtrOr(opt.ID),
		name: util.PtrOr(opt.Name),
	}
}

func (c *Command) WithContext(ctx context.Context) *Command {
	c.ctx = ctx
	return c
}
