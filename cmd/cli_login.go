package cmd

import (
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to kthcloud using Keycloak and retrieve the authentication token",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.Get()

		session, err := c.Login()
		if err != nil {
			log.Fatal(err)
		}
		if session == nil {
			log.Fatal("Could not login")
		}
		if c.HasValidSession() {
			log.Info("Logged in")
		}

		err = c.Session.Save(viper.GetString("session-path"))
		if err != nil {
			log.Errorln(err)
		}
		log.Info("Saved session to file")
	},
}

func init() {
	// Add the login command
	rootCmd.AddCommand(loginCmd)

	// Add flags for Keycloak credentials
	loginCmd.Flags().StringP("redirect-uri", "f", "http://localhost:3000", "Keycloak redirect endpoint URI")
	loginCmd.Flags().StringP("auth-url", "a", viper.GetString("keycloak-host")+"/realms/"+viper.GetString("keycloak-realm")+"/protocol/openid-connect/auth", "Keycloak auth endpoint URL")
	loginCmd.Flags().StringP("token-url", "t", viper.GetString("keycloak-host")+"/realms/"+viper.GetString("keycloak-realm")+"/protocol/openid-connect/token", "Keycloak token endpoint URL")
	viper.BindPFlag("auth-url", loginCmd.Flags().Lookup("auth-url"))
	viper.BindPFlag("token-url", loginCmd.Flags().Lookup("token-url"))
	viper.BindPFlag("redirect-uri", loginCmd.Flags().Lookup("redirect-uri"))

}
