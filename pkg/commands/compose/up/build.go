package up

import (
	"github.com/Phillezi/kthcloud-cli/pkg/builder"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/sirupsen/logrus"
)

func (c *Command) build() error {
	if c.buildAll {
		logrus.Debugln("buildAll is true")
		for n, s := range c.compose.Source.Services {
			if s.Build != nil {
				if err := builder.Build(c.client, c.ctx, n, s, c.nonInteractive); err != nil {
					logrus.Fatalln("Could not build service:", n, "Error:", err)
				}
				logrus.Debugln("build of", n, "is done!")
			}
		}
	} else if len(c.servicesToBuild) > 0 {
		logrus.Debugln("services to build are specified")
		for n, s := range c.compose.Source.Services {
			if s.Build != nil {
				if util.Contains(c.servicesToBuild, n) {
					if err := builder.Build(c.client, c.ctx, n, s, c.nonInteractive); err != nil {
						logrus.Fatalln("Could not build service:", n, "Error:", err)
					}
					logrus.Debugln("build of", n, "is done!")
				}
			}
		}
	}

	buildsReq, err := builder.GetBuildsRequired(c.client, *c.compose.Source)
	if err != nil {
		logrus.Fatalln("Error getting builds required:", err)
	}
	for n, needsBuild := range buildsReq {
		if needsBuild {
			if err := builder.Build(c.client, c.ctx, n, c.compose.Source.Services[n], c.nonInteractive); err != nil {
				logrus.Fatalln("Could not build service:", n, "Error:", err)
			}
			logrus.Debugln("build of", n, "is done!")
		}
	}
	return nil
}
