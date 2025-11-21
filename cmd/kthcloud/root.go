package main

import (
	viperconf "github.com/Phillezi/common/config/viper"
	zetup "github.com/Phillezi/common/logging/zap"
	"github.com/kthcloud/cli/internal/constants"
	"github.com/kthcloud/cli/internal/defaults"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:     "kthcloud",
	Short:   "A CLI tool for kthcloud.",
	Long:    banner,
	Version: version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		zetup.Setup()
	},
}

func init() {
	cobra.OnInitialize(func() { viperconf.InitConfig("kthcloud") })

	rootCmd.PersistentFlags().String(constants.ViperKeycloakClientId, defaults.DefaultKeycloakClientID, "Keycloak client ID")
	rootCmd.PersistentFlags().String(constants.ViperKeycloakBaseURL, defaults.DefaultKeycloakBaseURL, "Keycloak base URL")
	rootCmd.PersistentFlags().String(constants.ViperKeycloakRealm, defaults.DefaultKeycloakRealm, "Keycloak realm")

	rootCmd.PersistentFlags().String(constants.ViperDeployAPIToken, "", "Use an API token for authentication")
	rootCmd.PersistentFlags().String(constants.ViperSessionKey, defaults.DefaultKeystoreSessionKey, "The session key to store the session as, can be used with different users on the same computer user at the same time with this option")

	viper.BindPFlags(rootCmd.PersistentFlags())
}
