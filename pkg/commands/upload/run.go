package upload

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func (c *Command) Run() error {
	if c.client == nil {
		return fmt.Errorf("client is nil")
	}

	isAuth, err := c.client.Storage().WithFilebrowserURL(c.storageURL).Auth()
	if err != nil {
		logrus.Fatal(err)
	}
	if !isAuth {
		logrus.Fatal("user is not authenticated on storage manager")
	}
	content, err := os.ReadFile(c.srcPath)
	if err != nil {
		logrus.Fatal(err)
	}
	uploaded, err := c.client.Storage().UploadFile(c.destPath, content)
	if err != nil {
		logrus.Fatal(err)
	}
	if uploaded {
		logrus.Info("uploaded file!")
	}
	return nil
}
