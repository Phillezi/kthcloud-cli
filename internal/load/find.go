package load

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Phillezi/kthcloud-cli/pkg/defaults"
	"github.com/compose-spec/compose-go/v2/cli"
	"github.com/sirupsen/logrus"
)

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
				logrus.Warnf("Found multiple config files with supported names: %s", strings.Join(candidates, ", "))
				logrus.Warnf("Using %s", winner)
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
