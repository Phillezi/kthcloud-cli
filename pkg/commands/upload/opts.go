package upload

import (
	"context"

	"github.com/Phillezi/kthcloud-cli/pkg/deploy"
)

type CommandOpts struct {
	Context *context.Context
	Client  *deploy.Client

	SrcPath  *string
	DestPath *string

	StorageURL      *string
	KeycloakBaseURL *string
}
