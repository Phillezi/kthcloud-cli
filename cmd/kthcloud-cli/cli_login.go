package main

import (
	"kthcloud-cli/pkg/auth"
	"kthcloud-cli/pkg/util"
	"net/url"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Keycloak and retrieve the authentication token",
	Run: func(cmd *cobra.Command, args []string) {
		// Get login credentials
		clientID := viper.GetString("client-id")
		clientSecret := viper.GetString("client-secret")
		//username := viper.GetString("username")
		//password := viper.GetString("password")
		authURL := viper.GetString("auth-url")
		tokenURL := viper.GetString("token-url")
		redirectURI := viper.GetString("redirect-uri")

		/*// Call the Keycloak login function to get the token
		token, err := auth.KeycloakLogin(authURL, clientID, clientSecret, username, password)
		if err != nil {
			logrus.Fatalf("Login failed: %v", err)
		}

		// Store the token using viper (or securely as needed)
		viper.Set("auth-token", token)

		fmt.Println("Login successful. Token stored for future requests.")*/

		if clientID == "" || redirectURI == "" {
			log.Fatal("Client ID, and Redirect URI must be set in config.")
		}

		// Open browser for user login
		err := auth.OpenBrowser(authURL + "?client_id=" + clientID + "&redirect_uri=" + url.QueryEscape(redirectURI) + "&response_type=code&scope=openid")
		if err != nil {
			log.Fatalf("Failed to open browser: %v", err)
		}

		// Start local server to capture authorization code
		code, err := auth.StartLocalServer()
		if err != nil {
			log.Fatalf("Failed to start local server: %v", err)
		}

		// Exchange authorization code for access token
		token, err := auth.GetAccessToken(code, clientID, clientSecret, tokenURL, redirectURI)
		if err != nil {
			log.Fatalf("Failed to get access token: %v", err)
		}

		viper.Set("auth-token", token)

		configPath := "config.yaml"

		// Ensure the config file exists
		if err := util.EnsureFileExists(configPath); err != nil {
			log.Fatalf("Error: %v", err)
		}

		// Initialize Viper
		viper.SetConfigFile(configPath)

		// Write the token to the config file
		err = viper.WriteConfig()
		if err != nil {
			log.Fatalf("Failed to write token to config file: %v", err)
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
	loginCmd.Flags().StringP("auth-token", "x", "", "Authentication token")
	loginCmd.Flags().StringP("username", "u", "", "Username for login")
	loginCmd.Flags().StringP("password", "p", "", "Password for login")
	loginCmd.Flags().StringP("keycloak-host", "k", "https://iam.cloud.cbh.kth.se", "Keycloak server endpoint")
	loginCmd.Flags().StringP("redirect-uri", "f", "http://localhost:3000", "Keycloak redirect endpoint URI")
	viper.BindPFlag("client-id", loginCmd.Flags().Lookup("client-id"))
	viper.BindPFlag("keycloak-realm", loginCmd.Flags().Lookup("keycloak-realm"))
	viper.BindPFlag("client-secret", loginCmd.Flags().Lookup("client-secret"))
	viper.BindPFlag("auth-token", loginCmd.Flags().Lookup("auth-token"))
	viper.BindPFlag("username", loginCmd.Flags().Lookup("username"))
	viper.BindPFlag("password", loginCmd.Flags().Lookup("password"))
	viper.BindPFlag("keycloak-host", loginCmd.Flags().Lookup("keycloak-host"))
	loginCmd.Flags().StringP("auth-url", "a", viper.GetString("keycloak-host")+"/realms/"+viper.GetString("keycloak-realm")+"/protocol/openid-connect/auth", "Keycloak auth endpoint URL")
	loginCmd.Flags().StringP("token-url", "t", viper.GetString("keycloak-host")+"/realms/"+viper.GetString("keycloak-realm")+"/protocol/openid-connect/token", "Keycloak token endpoint URL")
	viper.BindPFlag("auth-url", loginCmd.Flags().Lookup("auth-url"))
	viper.BindPFlag("token-url", loginCmd.Flags().Lookup("token-url"))
	viper.BindPFlag("redirect-uri", loginCmd.Flags().Lookup("redirect-uri"))

}
