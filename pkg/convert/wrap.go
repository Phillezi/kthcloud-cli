package convert

import (
	"github.com/compose-spec/compose-go/v2/types"
	"github.com/kthcloud/go-deploy/dto/v2/body"
)

type Wrap struct {
	Deployments  []body.DeploymentCreate
	Dependencies map[string][]string
	Source       *types.Project
}
