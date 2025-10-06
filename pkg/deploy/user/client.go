package user

import (
	"github.com/kthcloud/cli/internal/types"
	"github.com/kthcloud/cli/pkg/deploy/common"
)

type GetUsersOption common.RequestOption

type Client interface {
	GetUser(id string) (types.TODO, error)
	GetUsers(id string, opts ...GetUsersOption) ([]types.TODO, error)

	CreateUser(user types.TODO) (string, error)

	UpdateUser(id string, user types.TODO) error

	DeleteUser(id string) error
}
