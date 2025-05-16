package create

import (
	"context"
	"os"

	"github.com/Phillezi/kthcloud-cli/pkg/defaults"
	"github.com/Phillezi/kthcloud-cli/pkg/deploy"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
)

type Command struct {
	ctx    context.Context
	client *deploy.Client

	rootDir        string
	createWorkFlow bool
	deploymentName string
}

func New(opts ...CommandOpts) *Command {
	var opt CommandOpts
	if len(opts) > 0 {
		opt = opts[0]
	}
	return &Command{
		ctx:    util.PtrOr(opt.Context, context.Background()),
		client: opt.Client,

		rootDir:        util.PtrOr(opt.RootDir, func() string { wd, _ := os.Getwd(); return wd }()),
		createWorkFlow: util.PtrOr(opt.CreateWorkFlow, defaults.DefaultCreateWorkflow),
		deploymentName: util.PtrOr(opt.DeploymentName),
	}
}

func (c *Command) WithContext(ctx context.Context) *Command {
	c.ctx = ctx
	return c
}
