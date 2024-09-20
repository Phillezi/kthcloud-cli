package upload

import (
	"os"

	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	storageclient "github.com/Phillezi/kthcloud-cli/pkg/v1/auth/storage-client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Upload(localPath, serverPath string) {
	c := client.Get()
	if !c.HasValidSession() {
		logrus.Fatal("no valid session, log in and try again")
	}
	user, err := c.User()
	if err != nil {
		logrus.Fatal(err)
	}
	sc := storageclient.GetInstance(*user.StorageURL, viper.GetString("keycloak-host"))
	isAuth, err := sc.Auth()
	if err != nil {
		logrus.Fatal(err)
	}
	if !isAuth {
		logrus.Fatal("user is not authenticated on storage manager")
	}
	content, err := os.ReadFile(localPath)
	if err != nil {
		logrus.Fatal(err)
	}
	uploaded, err := sc.UploadFile(serverPath, content)
	if err != nil {
		logrus.Fatal(err)
	}
	if uploaded {
		logrus.Info("uploaded file!")
	}
}
