package storage

import (
	"github.com/Phillezi/kthcloud-cli/pkg/deploy"
	"github.com/sirupsen/logrus"
)

func Check(client *deploy.Client) {
	user, err := client.User()
	if err != nil {
		logrus.Fatal(err)
	}
	if user.StorageURL == nil {
		logrus.Fatal("user doesnt have storageurl")
	}

	isAuth, err := client.Storage().WithFilebrowserURL(*user.StorageURL).Auth()
	if err != nil {
		logrus.Fatal(err)
	}
	if !isAuth {
		logrus.Fatal("not authenticated on storage url" + *user.StorageURL)
	}
	logrus.Infoln("Passed :)")
}
