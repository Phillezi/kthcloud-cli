package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kthcloud-cli/internal/api"
	"kthcloud-cli/pkg/util"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var apiCmd = &cobra.Command{
	Use:   "api [method] [resource]",
	Short: "Fetch the kthcloud api",
}

var getCmd = &cobra.Command{
	Use:   "get [resource]",
	Short: "Fetch the kthcloud api",
}

var postCmd = &cobra.Command{
	Use:   "post [resource]",
	Short: "Fetch the kthcloud api",
}

var putCmd = &cobra.Command{
	Use:   "put [resource]",
	Short: "Fetch the kthcloud api",
}

var deleteCmd = &cobra.Command{
	Use:   "delete [resource]",
	Short: "Fetch the kthcloud api",
}

var pathCmd = &cobra.Command{
	Use:   "path [resource]",
	Short: "Fetch data from the API with authorization",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logLevel, err := log.ParseLevel(viper.GetString("loglevel"))
		if err != nil {
			log.Fatal(err)
		}
		log.SetLevel(logLevel)

		resource := args[0]
		apiURL := viper.GetString("api-url")

		token := viper.GetString("auth-token")
		if token == "" {
			log.Fatal("No authentication token found. Please log in first.")
		}

		client := api.NewClient(apiURL, token)

		resp, err := client.FetchResource(resource, "GET")
		if err != nil {
			log.Fatalf("Failed to fetch resource: %v", err)
		}

		log.Infof("Response: %s", resp)
	},
}

func init() {

	rootCmd.AddCommand(apiCmd)
	getCmd.AddCommand(pathCmd)

	_, err := os.Stat("swagger.json")
	if os.IsNotExist(err) {
		log.Println("swagger.json not found, downloading...")

		// Download the swagger.json from the URL
		err = util.DownloadFile("swagger.json", "https://raw.githubusercontent.com/kthcloud/go-deploy/main/docs/api/v2/V2_swagger.json")
		if err != nil {
			log.Errorln("Failed to download swagger.json: %v\n", err)
			return
		}
		fmt.Println("Downloaded swagger.json successfully.")
	} else if err != nil {
		log.Errorln("Error checking file: %v\n", err)
		return
	} else {
		log.Debugln("swagger.json found.")
	}

	swagger, err := api.LoadSwaggerDoc("swagger.json")
	if err != nil {
		log.Fatal("Error loading swagger file: %v\n", err)
	}

	// Dynamically add commands based on the Swagger doc

	commands := CreateCommandsFromSwagger(swagger)
	for key, cmds := range commands {
		for _, cmd := range cmds {
			switch key {
			case "GET":
				getCmd.AddCommand(cmd)
			case "POST":
				postCmd.AddCommand(cmd)
			case "PUT":
				putCmd.AddCommand(cmd)
			case "DELETE":
				deleteCmd.AddCommand(cmd)
			default:
				log.Warnln("not supported method", key)
				apiCmd.AddCommand(cmd)
			}
		}
	}
	apiCmd.AddCommand(getCmd)
	apiCmd.AddCommand(postCmd)
	apiCmd.AddCommand(putCmd)
	apiCmd.AddCommand(deleteCmd)
}

func CreateCommandsFromSwagger(swagger *api.SwaggerDoc) map[string][]*cobra.Command {
	commandsByMethod := make(map[string][]*cobra.Command)

	for path, operations := range swagger.Paths {
		for method, operation := range operations {
			cmd := &cobra.Command{
				Use:   fmt.Sprintf("%s", path),
				Short: operation.Summary,
				Long:  operation.Description,
				Run: func(cmd *cobra.Command, args []string) {
					resource := path
					client := api.NewClient(viper.GetString("api-url"), viper.GetString("auth-token"))

					resp, err := client.FetchResource(resource, strings.ToUpper(method))
					if err != nil {
						log.Errorf("Error: %v\n", err)
					} else {
						var prettyJSON bytes.Buffer
						err = json.Indent(&prettyJSON, []byte(resp), "", "  ")
						if err != nil {
							log.Errorf("Error formatting JSON: %v\n", err)
							fmt.Println(resp)
						} else {
							fmt.Println(prettyJSON.String())
						}
					}
				},
			}

			for _, param := range operation.Parameters {
				cmd.Flags().String(param.Name, "", param.Description)
				viper.BindPFlag(param.Name, cmd.Flags().Lookup(param.Name))
			}

			commandsByMethod[strings.ToUpper(method)] = append(commandsByMethod[strings.ToUpper(method)], cmd)
		}
	}

	return commandsByMethod
}
