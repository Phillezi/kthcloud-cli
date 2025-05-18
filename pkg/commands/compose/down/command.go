package down

import (
	"context"

	"github.com/Phillezi/kthcloud-cli/pkg/convert"
	"github.com/Phillezi/kthcloud-cli/pkg/defaults"
	"github.com/Phillezi/kthcloud-cli/pkg/deploy"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
)

type Command struct {
	ctx    context.Context
	client *deploy.Client

	compose *convert.Wrap

	all bool
}

func New(opts ...CommandOpts) *Command {
	var opt CommandOpts
	if len(opts) > 0 {
		opt = opts[0]
	}
	return &Command{
		ctx:    context.Background(),
		client: opt.Client,

		compose: opt.Compose,

		all: util.PtrOr(opt.All, defaults.DefaultRemoveCustomDeployments),
	}
}

func (c *Command) WithContext(ctx context.Context) *Command {
	c.ctx = ctx
	return c
}
