package ps

import (
	"github.com/Phillezi/kthcloud-cli/pkg/deploy"
)

type CommandOpts struct {
	Client *deploy.Client

	All *bool
}
