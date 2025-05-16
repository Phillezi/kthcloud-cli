package create

import (
	"context"

	"github.com/Phillezi/kthcloud-cli/pkg/deploy"
)

type CommandOpts struct {
	Context *context.Context
	Client  *deploy.Client

	RootDir        *string
	CreateWorkFlow *bool
	DeploymentName *string
}
