package main

import (
	"go-deploy/dto/v2/body"
	"kthcloud-cli/internal/model"
	"kthcloud-cli/pkg/util"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var utilCmd = &cobra.Command{
	Use:   "util",
	Short: "Utility functionality",
}

var tokenCmd = &cobra.Command{
	Use:   "api-token",
	Short: "Generate and save api token",
	Run: func(cmd *cobra.Command, args []string) {
		session, err := model.Load(viper.GetString("session-path"))
		if err != nil {
			log.Fatalln("No active session. Please log in")
		}
		if session.AuthSession.IsExpired() {
			log.Fatalln("Session is expired. Please log in again")
		}
		session.SetupClient()
		session.FetchUser()

		name := "cli-access"

		// Check if it already is in the user.ApiKeys array
		if util.Contains(util.GetNames(session.User.ApiKeys), name) {
			// TODO: check how we want to handle this, should a new name be genereated for the key
			log.Fatalf("A token with the name '%s' already exists in the user's API keys.\n", name)
		}

		// create a new api token that expires in one month
		resp, err := session.Client.Req("/v2/users/"+session.User.ID+"/apiKeys", "POST", &body.ApiKeyCreate{Name: name, ExpiresAt: time.Now().AddDate(0, 1, 0)})
		if err != nil {
			log.Fatalln(err)
		}

		util.HandleResponse(resp)
		key, err := util.ProcessResponse[body.ApiKeyCreated](resp.String())
		if err != nil {
			log.Fatalln(err)
		}

		session.ApiKey = &model.ApiKey{Key: key.Key, Expiry: key.ExpiresAt}
		err = session.Save(viper.GetString("session-path"))
		if err != nil {
			log.Fatalln(err)
		}
		log.Infoln("Successfully generated and added an api token")
	},
}

func init() {
	utilCmd.AddCommand(tokenCmd)

	rootCmd.AddCommand(utilCmd)
}
