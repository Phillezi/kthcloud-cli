package auth

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

func GetToken() (string, error) {
	token := viper.GetString("api-token")
	if token == "" {
		token = viper.GetString("auth-token")
		if token == "" {
			log.Errorln("No authentication token found. Please log in first.")
			return "", errors.New("No authentication token found. Please log in first.")
		}
	}
	return token, nil
}
