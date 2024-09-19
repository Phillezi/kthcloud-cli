package storage

import (
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	storageclient "github.com/Phillezi/kthcloud-cli/pkg/v1/auth/storage-client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Check() {
	c := client.Get()
	if !c.HasValidSession() {
		logrus.Fatal("not logged in")
	}
	user, err := c.User()
	if err != nil {
		logrus.Fatal(err)
	}
	if user.StorageURL == nil {
		logrus.Fatal("user doesnt have storageurl")
	}
	if c.StorageClient == nil {
		c.StorageClient = storageclient.GetInstance(*user.StorageURL, viper.GetString("keycloak-host"))
	}
	isAuth, err := c.StorageClient.Auth()
	if err != nil {
		logrus.Fatal(err)
	}
	if !isAuth {
		logrus.Fatal("not authenticated on storage url" + *user.StorageURL)
	}
	logrus.Infoln("Passed :)")
}
