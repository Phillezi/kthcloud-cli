package create

import (
	"github.com/Phillezi/kthcloud-cli/pkg/deploy"
)

type CommandOpts struct {
	Client *deploy.Client

	RootDir        *string
	CreateWorkFlow *bool
	DeploymentName *string
}
