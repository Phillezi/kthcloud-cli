package app

import (
	"context"

	"github.com/kthcloud/cli/internal/constants"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func FromViper(ctx context.Context, logger *zap.Logger, v viper.Viper) *App {
	config := Config{
		DeployAPIBaseURL: viper.GetString(constants.ViperDeployAPIBaseURL),

		KeycloakClientID:     viper.GetString(constants.ViperKeycloakClientId),
		KeycloakClientSecret: viper.GetString(constants.ViperKeycloakClientSecret),
		KeycloakBaseURL:      viper.GetString(constants.ViperKeycloakBaseURL),
		KeycloakRealm:        viper.GetString(constants.ViperKeycloakRealm),

		LoginServerAddress: viper.GetString(constants.ViperLoginServerAddress),

		SessionKey:         viper.GetString(constants.ViperSessionKey),
		SessionService:     viper.GetString(constants.ViperSessionService),
		SessionFallbackDir: viper.GetString(constants.ViperSessionFallbackDir),
	}

	if logger != nil {
		config.Logger = logger
	}

	return New(ctx, WithConfig(config))
}
