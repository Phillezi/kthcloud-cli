package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/kthcloud/go-deploy/dto/v2/body"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/commands/upload"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var utilCmd = &cobra.Command{
	Use:   "util",
	Short: "Utility functionality",
}

var tokenCmd = &cobra.Command{
	Use:   "api-key",
	Short: "Generate and save a new api key",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.Get()
		name := "cli-access"
		c.DropUserCache()
		user, err := c.User()
		if err != nil {
			logrus.Fatal(err)
		}

		// Check if it already is in the user.ApiKeys array
		for i := 0; util.Contains(util.GetNames(user.ApiKeys), name); i++ {
			// if it already exists, try to add/increment a number to it and check again
			name = fmt.Sprintf("cli-access-%d", i)
		}

		// create a new api token that expires in one month
		resp, err := c.Create(&body.ApiKeyCreate{Name: name, ExpiresAt: time.Now().AddDate(0, 1, 0)})
		if err != nil {
			logrus.Fatalln(err)
		}

		util.HandleResponse(resp)
		key, err := util.ProcessResponse[body.ApiKeyCreated](resp.String())
		if err != nil {
			logrus.Fatalln(err)
		}

		c.Session.ApiKey = key
		err = c.Session.Save(viper.GetString("session-path"))
		if err != nil {
			logrus.Fatalln(err)
		}
		logrus.Infoln("Successfully generated and added an api token")
	},
}
var uploadCmd = &cobra.Command{
	Use:   "upload <local-file-path> <server-file-path>",
	Short: "Upload a file",
	Long: `Upload a file to the server.

Arguments:
  <local-file-path>   The local path to the file that you want to upload.
  <server-file-path>  The destination path on the server where the file will be uploaded, including the filename.`,
	Example: "upload ./myfile.txt existingpath/myfile.txt",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Fatal(cmd.Usage())
		}
		upload.Upload(args[0], args[1])
	},
}

func init() {
	utilCmd.AddCommand(tokenCmd)
	utilCmd.AddCommand(uploadCmd)

	rootCmd.AddCommand(utilCmd)
}
