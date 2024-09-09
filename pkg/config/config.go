package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// InitConfig initializes configuration
func InitConfig() {
	viper.SetConfigName("config") // Name of the config file (without extension)
	viper.SetConfigType("yaml")   // File format (yaml)
	viper.AddConfigPath(".")      // Search in the current directory
	viper.AutomaticEnv()          // Read environment variables

	// Load config file
	if err := viper.ReadInConfig(); err != nil {
		log.Warn("Config file not found, using defaults or environment variables.")
	} else {
		log.Debugf("Using config file: %s", viper.ConfigFileUsed())
	}
}
