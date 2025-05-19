package up

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

	compose         *convert.Wrap
	services        []string
	servicesToBuild []string

	buildAll       bool
	detached       bool
	nonInteractive bool

	// state
	creationDone bool
	cancelled    bool
}

func New(opts ...CommandOpts) *Command {
	var opt CommandOpts
	if len(opts) > 0 {
		opt = opts[0]
	}
	return &Command{
		ctx:    context.Background(),
		client: opt.Client,

		compose:         opt.Compose,
		services:        opt.Services,
		servicesToBuild: opt.ServicesToBuild,

		buildAll:       util.PtrOr(opt.BuildAll, defaults.DefaultComposeUpBuildAll),
		detached:       util.PtrOr(opt.Detatched, defaults.DefaultComposeUpDetached),
		nonInteractive: util.PtrOr(opt.NonInteractive, defaults.DefaultNonInteractive),
	}
}

func (c *Command) WithContext(ctx context.Context) *Command {
	c.ctx = ctx
	return c
}
