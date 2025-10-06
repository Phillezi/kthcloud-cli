package vm

import (
	"github.com/kthcloud/cli/internal/types"
	"github.com/kthcloud/cli/pkg/deploy/common"
)

type GetVMsOption common.RequestOption

type Client interface {
	GetVM(id string) (types.TODO, error)
	GetVMs(opts ...GetVMsOption) ([]types.TODO, error)

	CreateVM(depl types.TODO) (string, error)

	UpdateVM(id string, depl types.TODO) (string, error)

	DeleteVM(id string) error
}
