package deploy

import (
	"github.com/kthcloud/cli/pkg/deploy/deployment"
	"github.com/kthcloud/cli/pkg/deploy/user"
)

type Client interface {
	deployment.Client
	user.Client
}
