package compose

import "errors"

var (
	ErrNoNameOnService             = errors.New("service is required to have a name")
	ErrBuildAndImageProvided       = errors.New("service cant provide both build and image")
	ErrInvalidDeploymentVisibility = errors.New("invalid deployment visibility, expected public, private or auth")
	ErrCustomDomainTooLong         = errors.New("the provided custom domain is too long, please make sure that it doesnt exceed the length of 243 characters")
)

type Warning error

var (
	WarnNotImplServiceDeployResourcesLimits = Warning(errors.New("service.Deploy.Resources.Limits is not implemented"))
)
