package down

import (
	"github.com/Phillezi/kthcloud-cli/pkg/convert"
	"github.com/Phillezi/kthcloud-cli/pkg/deploy"
)

type CommandOpts struct {
	Client *deploy.Client

	Compose *convert.Wrap

	All *bool
}
