package compose

import (
	"kthcloud-cli/pkg/auth"
	"kthcloud-cli/pkg/util"

	log "github.com/sirupsen/logrus"
)

func Up(filename string) error {
	services, err := ParseComposeFile(filename)
	if err != nil {
		log.Errorln(err)
	}

	client, err := auth.GetClient()
	if err != nil {
		return err
	}

	for key, service := range services {
		resp, err := client.Req("/v2/deployments", "POST", serviceToDepl(service, key))
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
