package auth

import (
	"errors"
	"kthcloud-cli/internal/api"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

func GetToken() (string, error) {
	token := viper.GetString("api-token")
	if token == "" {
		token = viper.GetString("auth-token")
		if token == "" {
			log.Errorln("No authentication token found. Please log in first.")
			return "", errors.New("no authentication token found. please log in first")
		}
	}
	return token, nil
}

func GetClient() (*api.Client, error) {
	token := viper.GetString("api-token")
	if token == "" {
		token = viper.GetString("auth-token")
		if token == "" {
			log.Errorln("No authentication token found. Please log in first.")
			return nil, errors.New("no authentication token found. please log in first")
		} else {
			return api.NewClient(viper.GetString("api-url"), token), nil
		}
	}
	return api.NewClient(viper.GetString("api-url"), token), nil
}
