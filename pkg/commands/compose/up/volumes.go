package up

import (
	"github.com/Phillezi/kthcloud-cli/pkg/storage"
	"github.com/sirupsen/logrus"
)

func (c *Command) volumes() error {

	_, err := storage.CreateVolumes(c.client, c.compose)
	if err != nil {
		logrus.Error("could not create volumes: ", err)
		return err
	}
	return nil
}
