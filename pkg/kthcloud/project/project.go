package project

import (
	"github.com/kthcloud/cli/pkg/deploy"
)

type Project struct {
	Sevices Services
}

type Services = map[string]Service

type Service struct {
	deploy.BodyDeploymentCreate `yaml:",inline" json:",inline"`
	Dependencies                []string
}

func (s Service) GetDependencies() []string {
	return s.Dependencies
}
