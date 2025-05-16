package up

import (
	"github.com/Phillezi/kthcloud-cli/pkg/storage"
	"github.com/sirupsen/logrus"
)

func (c *Command) volumes() error {
	if c.tryVolumes {
		_, err := storage.CreateVolumes(c.client, c.compose)
		if err != nil {
			logrus.Fatal(err)
		}
	} else {
		logrus.Infoln("Skipping volume creation from local structure")
		logrus.Infoln("If enabled it will \"steal\" cookies from your browser to authenticate")
		logrus.Infoln("use --try-volumes to try")
	}
	return nil
}
