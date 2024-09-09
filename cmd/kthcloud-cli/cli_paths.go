package main

import (
	"fmt"
	"kthcloud-cli/internal/api"
	"kthcloud-cli/pkg/util"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	_, err := os.Stat("swagger.json")
	if os.IsNotExist(err) {
		log.Println("swagger.json not found, downloading...")

		// Download the swagger.json from the URL
		err = util.DownloadFile("swagger.json", "https://raw.githubusercontent.com/kthcloud/go-deploy/main/docs/api/v2/V2_swagger.json")
		if err != nil {
			log.Printf("Failed to download swagger.json: %v\n", err)
			return
		}
		fmt.Println("Downloaded swagger.json successfully.")
	} else if err != nil {
		log.Printf("Error checking file: %v\n", err)
		return
	} else {
		log.Println("swagger.json found.")
	}

	swagger, err := api.LoadSwaggerDoc("swagger.json")
	if err != nil {
		log.Fatal("Error loading swagger file: %v\n", err)
	}

	// Dynamically add commands based on the Swagger doc

	// wait until login init has ben run
	commands := CreateCommandsFromSwagger(swagger)
	for _, cmd := range commands {
		rootCmd.AddCommand(cmd)
	}
}

func CreateCommandsFromSwagger(swagger *api.SwaggerDoc) []*cobra.Command {
	var commands []*cobra.Command

	token := viper.GetString("auth-token")

	if token == "" {
		log.Warn("not logged in")
		return commands
	}

	for path, operations := range swagger.Paths {
		for method, operation := range operations {
			cmd := &cobra.Command{
				Use:   fmt.Sprintf("%s-%s", method, path),
				Short: operation.Summary,
				Long:  operation.Description,
				Run: func(cmd *cobra.Command, args []string) {
					// Prepare the resource request based on the method
					resource := path // You can modify this to handle different parameters

					client := api.NewClient(viper.GetString("api-url"), token)
					// Call the client to fetch the resource
					resp, err := client.FetchResource(resource, strings.ToUpper(method))
					if err != nil {
						fmt.Printf("Error: %v\n", err)
					} else {
						fmt.Println(resp)
					}
				},
			}

			// Add parameters to the command
			for _, param := range operation.Parameters {
				cmd.Flags().String(param.Name, "", param.Description)
				viper.BindPFlag(param.Name, cmd.Flags().Lookup(param.Name))
			}

			commands = append(commands, cmd)
		}
	}

	return commands
}
