package project

import (
	"github.com/kthcloud/cli/pkg/deploy"
)

type Project struct {
	Sevices Services `yaml:"services,omitempty" json:"services,omitempty"`
}

type Build struct {
	ContainerFile string
	Context       string
}

type Services = map[string]Service

type Service struct {
	deploy.BodyDeploymentCreate `yaml:",inline" json:",inline"`
	Dependencies                []string `yaml:"depends_on,omitempty" json:"depends_on,omitempty"`
	Build                       *Build   `yaml:"build,omitempty" json:"build,omitempty"`
}

func (s Service) GetDependencies() []string {
	return s.Dependencies
}

func (s Service) Body() deploy.BodyDeploymentCreate {
	return s.BodyDeploymentCreate
}
