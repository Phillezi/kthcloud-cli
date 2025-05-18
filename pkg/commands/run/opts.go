package run

import (
	"github.com/Phillezi/kthcloud-cli/pkg/deploy"
)

type CommandOpts struct {
	Client *deploy.Client

	Interactive *bool
	TTY         *bool
	Remove      *bool
	Detatch     *bool

	Envs map[string]string
	Port []int

	Visibility *string

	Image *string

	Memory   *float64
	Cores    *float64
	Replicas *int

	Name *string
}
