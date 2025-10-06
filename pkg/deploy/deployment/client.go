package deployment

import (
	"github.com/kthcloud/cli/internal/types"
	"github.com/kthcloud/cli/pkg/deploy/common"
)

type GetDeploymentsOption common.RequestOption

type Client interface {
	GetDeployment(id string) (types.TODO, error)
	GetDeployments(opts ...GetDeploymentsOption) ([]types.TODO, error)

	CreateDeployment(depl types.TODO) (string, error)

	UpdateDeployment(id string, depl types.TODO) (string, error)

	DeleteDeployment(id string) error
}
