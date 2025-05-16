package down

import (
	"context"

	"github.com/Phillezi/kthcloud-cli/pkg/deploy"
	"github.com/Phillezi/kthcloud-cli/pkg/models/compose"
)

type CommandOpts struct {
	Context *context.Context
	Client  *deploy.Client

	Compose *compose.Compose

	All *bool
}
