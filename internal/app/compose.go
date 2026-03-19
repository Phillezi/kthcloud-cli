package app

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/compose-spec/compose-go/v2/cli"
	"github.com/compose-spec/compose-go/v2/consts"
	"github.com/kthcloud/cli/internal/defaults"
	"github.com/kthcloud/cli/pkg/compose"
	"github.com/kthcloud/cli/pkg/kthcloud/project"
	"go.uber.org/zap"
)

func (app *App) Compose(composeFilePaths []string, projectName string) (*project.Project, error) {
	options, err := func() (*cli.ProjectOptions, error) {
		if len(composeFilePaths) > 0 {
			return cli.NewProjectOptions(
				composeFilePaths,
				WithDefaultConfigPath,
				WithEnvFiles(defaults.DefaultComposeDotEnvFileNames...),
				cli.WithOsEnv,
				cli.WithDotEnv,
				// cli.WithName(projectName),
			)
		}
		return cli.NewProjectOptions(
			nil,
			WithDefaultConfigPath,
			WithEnvFiles(defaults.DefaultComposeDotEnvFileNames...),
			cli.WithOsEnv,
			cli.WithDotEnv,
			// cli.WithName(projectName),
		)
	}()
	if err != nil {
		return nil, err
	}

	project, err := options.LoadProject(app.ctx)
	if err != nil {
		return nil, err
	}

	return compose.ConverterImpl{}.Convert(project)
}

func WithDefaultConfigPath(o *cli.ProjectOptions) error {
	if len(o.ConfigPaths) > 0 {
		return nil
	}
	pwd, err := o.GetWorkingDir()
	if err != nil {
		return err
	}
	for {
		candidates := findFiles(defaults.DefaultComposeFileNames, pwd)
		if len(candidates) > 0 {
			winner := candidates[0]
			if len(candidates) > 1 {
				zap.L().Sugar().Warnf("Found multiple config files with supported names: %s", strings.Join(candidates, ", "))
				zap.L().Sugar().Warnf("Using %s", winner)
			}
			o.ConfigPaths = append(o.ConfigPaths, winner)
		}
		parent := filepath.Dir(pwd)
		if parent == pwd {
			// no config file found, but that's not a blocker if caller only needs project name
			return nil
		}
		pwd = parent
	}
}

func WithEnvFiles(file ...string) cli.ProjectOptionsFn {
	return func(o *cli.ProjectOptions) error {
		if len(file) > 0 {
			wd, err := o.GetWorkingDir()
			if err != nil {
				return err
			}
			o.EnvFiles = findFiles(file, wd)
			return nil
		}
		if v, ok := os.LookupEnv(consts.ComposeDisableDefaultEnvFile); ok {
			b, err := strconv.ParseBool(v)
			if err != nil {
				return err
			}
			if b {
				return nil
			}
		}

		wd, err := o.GetWorkingDir()
		if err != nil {
			return err
		}
		defaultDotEnv := filepath.Join(wd, ".env")

		s, err := os.Stat(defaultDotEnv)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		if !s.IsDir() {
			o.EnvFiles = []string{defaultDotEnv}
		}
		return nil
	}
}

func findFiles(names []string, pwd string) []string {
	candidates := []string{}
	for _, n := range names {
		f := filepath.Join(pwd, n)
		if _, err := os.Stat(f); err == nil {
			candidates = append(candidates, f)
		}
	}
	return candidates
}
