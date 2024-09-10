package compose

import (
	"kthcloud-cli/internal/model"
	"kthcloud-cli/pkg/util"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Up(filename string) error {
	services, err := ParseComposeFile(filename)
	if err != nil {
		log.Errorln(err)
	}

	// load the session from the session.json file
	session, err := model.Load(viper.GetString("session-path"))
	if err != nil {
		log.Fatalln("No active session. Please log in")
	}
	if session.AuthSession.IsExpired() {
		log.Fatalln("Session is expired. Please log in again")
	}
	session.SetupClient()

	for key, service := range services {
		resp, err := session.Client.Req("/v2/deployments", "POST", serviceToDepl(service, key))
		if err != nil {
			log.Errorln("error: ", err, " response: ", resp)
			return err
		}
		if err := util.HandleResponse(resp); err != nil {
			return err
		}
		log.Info("response: ", resp)
		log.Info("Created deployment: ", key)
	}
	return nil
}
