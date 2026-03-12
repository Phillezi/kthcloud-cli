package compose

import (
	composev2 "github.com/compose-spec/compose-go/v2/types"
	"github.com/kthcloud/cli/pkg/kthcloud/project"
)

type Converter interface {
	Convert(project *composev2.Project) (*project.Project, error)
}
