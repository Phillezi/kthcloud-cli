package check

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func (c *Command) Run() error {
	if c.client == nil {
		return fmt.Errorf("client is nil")
	}

	user, err := c.client.User()
	if err != nil {
		logrus.Fatal(err)
	}
	if user.StorageURL == nil {
		logrus.Fatal("user doesnt have storageurl")
	}

	isAuth, err := c.client.Storage().Auth()
	if err != nil {
		logrus.Fatal(err)
	}
	if !isAuth {
		logrus.Fatal("not authenticated on storage url" + *user.StorageURL)
	}
	logrus.Infoln("Passed :)")
	return nil
}
