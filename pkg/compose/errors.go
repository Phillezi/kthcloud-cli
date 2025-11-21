package compose

import (
	"errors"

	"github.com/kthcloud/cli/pkg/warnings"
)

var (
	ErrNoNameOnService             = errors.New("service is required to have a name")
	ErrBuildAndImageProvided       = errors.New("service cant provide both build and image")
	ErrInvalidDeploymentVisibility = errors.New("invalid deployment visibility, expected public, private or auth")
	ErrCustomDomainTooLong         = errors.New("the provided custom domain is too long, please make sure that it doesnt exceed the length of 243 characters")
)

var (
	WarnNotImplServiceDeployResourcesLimits = warnings.New("service.Deploy.Resources.Limits is not implemented")
)
