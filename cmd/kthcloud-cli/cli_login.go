package main

import (
	"kthcloud-cli/internal/model"
	"kthcloud-cli/pkg/auth"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to kthcloud using Keycloak and retrieve the authentication token",
	Run: func(cmd *cobra.Command, args []string) {
		//clientID := viper.GetString("client-id")
		//clientSecret := viper.GetString("client-secret")
		//authURL := viper.GetString("auth-url")
		//tokenURL := viper.GetString("token-url")
		//redirectURI := viper.GetString("redirect-uri")

		//url := authURL + "?client_id=" + clientID + "&redirect_uri=" + url.QueryEscape(redirectURI) + "&response_type=code&scope=openid"

		err := auth.OpenBrowser("http://localhost:3000")
		if err != nil {
			log.Fatalf("Failed to open browser: %v", err)
		}

		authSession, err := auth.StartLocalServer()
		if err != nil {
			log.Fatalf("Failed to start local server: %v", err)
		}

		session := model.NewSession(authSession)

		err = session.FetchUser()
		if err != nil {
			log.Fatalln(err)
		}
		err = session.Save(viper.GetString("session-path"))
		if err != nil {
			log.Fatalln("eeee", err)
		}

		log.Infoln("Login successful. Access token stored in config file.")
	},
}

func init() {
	// Add the login command
	rootCmd.AddCommand(loginCmd)

	// Add flags for Keycloak credentials
	loginCmd.Flags().StringP("client-id", "c", "landing", "Keycloak client ID")
	loginCmd.Flags().StringP("keycloak-realm", "r", "cloud", "Keycloak realm")
	loginCmd.Flags().StringP("client-secret", "s", "", "Keycloak client secret")
	loginCmd.Flags().StringP("keycloak-host", "k", "https://iam.cloud.cbh.kth.se", "Keycloak server endpoint")
	loginCmd.Flags().StringP("redirect-uri", "f", "http://localhost:3000", "Keycloak redirect endpoint URI")
	viper.BindPFlag("client-id", loginCmd.Flags().Lookup("client-id"))
	viper.BindPFlag("keycloak-realm", loginCmd.Flags().Lookup("keycloak-realm"))
	viper.BindPFlag("client-secret", loginCmd.Flags().Lookup("client-secret"))
	viper.BindPFlag("keycloak-host", loginCmd.Flags().Lookup("keycloak-host"))
	loginCmd.Flags().StringP("auth-url", "a", viper.GetString("keycloak-host")+"/realms/"+viper.GetString("keycloak-realm")+"/protocol/openid-connect/auth", "Keycloak auth endpoint URL")
	loginCmd.Flags().StringP("token-url", "t", viper.GetString("keycloak-host")+"/realms/"+viper.GetString("keycloak-realm")+"/protocol/openid-connect/token", "Keycloak token endpoint URL")
	viper.BindPFlag("auth-url", loginCmd.Flags().Lookup("auth-url"))
	viper.BindPFlag("token-url", loginCmd.Flags().Lookup("token-url"))
	viper.BindPFlag("redirect-uri", loginCmd.Flags().Lookup("redirect-uri"))

}
