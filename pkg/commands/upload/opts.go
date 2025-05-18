package upload

import (
	"github.com/Phillezi/kthcloud-cli/pkg/deploy"
)

type CommandOpts struct {
	Client *deploy.Client

	SrcPath  *string
	DestPath *string

	StorageURL      *string
	KeycloakBaseURL *string
}
