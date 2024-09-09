package compose

import (
	"encoding/json"
	"fmt"
	"kthcloud-cli/internal/api"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Up(filename string) error {
	services, err := ParseComposeFile(filename)
	if err != nil {
		log.Errorln(err)
	}

	client := api.NewClient(viper.GetString("api-url"), viper.GetString("auth-token"))

	for key, service := range services {
		resp, err := client.Req("/v2/deployments", "POST", serviceToDepl(service, key))
		if err != nil {
			log.Errorln("error: ", err, " response: ", resp)
			return err
		}
		var responseMap map[string]interface{}
		if err := json.Unmarshal([]byte(resp), &responseMap); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}

		if errors, ok := responseMap["errors"]; ok {
			log.Errorf("response contains errors: %v", errors)
			return fmt.Errorf("response contains errors: %v", errors)
		}
		log.Info("response: ", resp)
		log.Info("Created deployment: ", key)
	}
	return nil
}
