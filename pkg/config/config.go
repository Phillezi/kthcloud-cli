package config

import (
	"errors"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// InitConfig initializes configuration
func InitConfig() {
	viper.SetConfigName("config") // Name of the config file (without extension)
	viper.SetConfigType("yaml")   // File format (yaml)
	viper.AddConfigPath(GetConfigPath())
	viper.AutomaticEnv() // Read environment variables

	// Load config file
	if err := viper.ReadInConfig(); err != nil {
		logrus.Debugf("Config file not found, using defaults or environment variables.\n")
	} else {
		logrus.Debugf("Using config file: %s", viper.ConfigFileUsed())
	}

}

func getConfigPath() (string, error) {
	basePath, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	configPath := path.Join(basePath, ".kthcloud")
	fileDescr, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(configPath, os.ModePerm)
		if err != nil {
			return "", err
		}
		fileDescr, err = os.Stat(configPath)
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}
	if !fileDescr.IsDir() {
		return "", errors.New("default config dir is file")
	}
	return configPath, nil
}

func GetConfigPath() string {
	configPath, err := getConfigPath()
	if err != nil {
		logrus.Errorln(err)
		configPath = "."
	}
	return configPath
}
