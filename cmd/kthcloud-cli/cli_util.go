package main

import (
	"go-deploy/dto/v2/body"
	"kthcloud-cli/pkg/auth"
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
		client, err := auth.GetClient()
		if err != nil {
			log.Fatalln(err)
		}

		resp, err := client.Req("/v2/users", "GET", nil)
		if err != nil {
			log.Fatalln(err)
		}
		if resp.IsError() {
			log.Fatalln("non ok responsecode")
		}

		users, err := util.ProcessResponseArr[body.UserRead](resp.String())
		if err != nil {
			log.Fatalln(err)
		}

		if len(users) != 1 {
			log.Fatalln("recieved more than one user")
		}

		user := users[0]

		name := "cli-access"

		// check if it already is in the user.ApiKeys array
		if util.Contains(util.GetNames(user.ApiKeys), name) {
			// TODO: check how we want to handle this, should a new name be genereated for the key
			log.Fatalf("A token with the name '%s' already exists in the user's API keys.\n", name)
		}

		// create a new api token that expires in one month
		resp, err = client.Req("/v2/users/"+user.ID+"/apiKeys", "POST", &body.ApiKeyCreate{Name: name, ExpiresAt: time.Now().AddDate(0, 1, 0)})
		if err != nil {
			log.Fatalln(err)
		}

		util.HandleResponse(resp)
		key, err := util.ProcessResponse[body.ApiKeyCreated](resp.String())
		if err != nil {
			log.Fatalln(err)
		}
		// TODO: handle better
		viper.Set("api-token", key.Key)
		viper.WriteConfig()
		log.Infoln("Successfully generated and added an api token")
	},
}

func init() {
	utilCmd.AddCommand(tokenCmd)

	rootCmd.AddCommand(utilCmd)
}
