package init

import (
	"github.com/Phillezi/kthcloud-cli/pkg/deploy"
)

type CommandOpts struct {
	Client      *deploy.Client
	SaveSecrets *bool
}
