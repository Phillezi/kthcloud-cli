package load

import (
	"sync"

	"github.com/Phillezi/kthcloud-cli/internal/interrupt"
	"github.com/Phillezi/kthcloud-cli/pkg/convert"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/compose-spec/compose-go/v2/cli"
	"github.com/spf13/viper"
)

var (
	once     sync.Once
	instance *convert.Wrap
	lastErr  error
)

type LoadOpts struct {
	File string
}

func GetCompose(opts ...LoadOpts) (*convert.Wrap, error) {
	once.Do(func() {
		instance, lastErr = InternalGetCompose(opts...)
	})
	return instance, lastErr
}

// Needs to be public for testing
func InternalGetCompose(opts ...LoadOpts) (*convert.Wrap, error) {
	var opt LoadOpts
	if len(opts) > 0 {
		opt = opts[0]
	}
	composeFilePath := util.Or(opt.File, viper.GetString("file"))
	projectName := "kthcloud-cli-compose-project"

	options, err :=
		func() (*cli.ProjectOptions, error) {
			if composeFilePath != "" {
				return cli.NewProjectOptions(
					[]string{composeFilePath},
					WithDefaultConfigPath,
					cli.WithOsEnv,
					cli.WithDotEnv,
					cli.WithName(projectName),
				)
			}
			return cli.NewProjectOptions(
				nil,
				WithDefaultConfigPath,
				cli.WithOsEnv,
				cli.WithDotEnv,
				cli.WithName(projectName),
			)

		}()
	if err != nil {
		return nil, err
	}

	project, err := options.LoadProject(interrupt.GetInstance().Context())
	if err != nil {
		return nil, err
	}

	var wrap convert.Wrap
	err = convert.ToCloud(project, &wrap)
	if err != nil {
		return nil, err
	}
	return &wrap, nil
}
