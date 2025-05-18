package run

import (
	"context"

	"github.com/Phillezi/kthcloud-cli/pkg/defaults"
	"github.com/Phillezi/kthcloud-cli/pkg/deploy"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
)

type Command struct {
	ctx    context.Context
	client *deploy.Client

	interactive bool
	tty         bool
	remove      bool
	detatch     bool

	envs map[string]string
	port []int

	visibility string

	image string

	memory   float64
	cores    float64
	replicas int

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

		interactive: util.PtrOr(opt.Interactive, defaults.DefaultDeploymentRunInteractive),
		tty:         util.PtrOr(opt.TTY, defaults.DefaultDeploymentRunTTY),
		remove:      util.PtrOr(opt.Remove, defaults.DefaultDeploymentRunRemove),
		detatch:     util.PtrOr(opt.Detatch, defaults.DefaultDeploymentRunDetatch),

		envs: opt.Envs,
		port: opt.Port,

		visibility: util.PtrOr(opt.Visibility, defaults.DefaultDeploymentVisibility),

		image: util.PtrOr(opt.Image),

		memory:   util.PtrOr(opt.Memory, defaults.DefaultDeploymentMemory),
		cores:    util.PtrOr(opt.Cores, defaults.DefaultDeploymentCores),
		replicas: util.PtrOr(opt.Replicas, defaults.DefaultDeploymentReplicas),

		name: util.PtrOr(opt.Name),
	}
}

func (c *Command) WithContext(ctx context.Context) *Command {
	c.ctx = ctx
	return c
}
