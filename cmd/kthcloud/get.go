package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"text/tabwriter"

	"github.com/kthcloud/cli/internal/app"
	"github.com/kthcloud/cli/pkg/deploy"
	"github.com/kthcloud/cli/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get deployments",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
		defer cancel()

		a := app.New(ctx, app.WithKeycloakOptions(
			viper.GetString("keycloak-client-id"),
			viper.GetString("keycloak-base-url"),
			viper.GetString("keycloak-realm"),
		),
			app.WithSessionKey(viper.GetString("session-key")),
			app.WithLogger(zap.L()),
		)

		params := &deploy.GetV2DeploymentsParams{
			All:    utils.PtrOf(viper.GetBool("all")),
			Shared: utils.PtrOf(viper.GetBool("shared")),
		}
		if userIDFilter := viper.GetString("by-user-id"); userIDFilter != "" {
			params.UserId = &userIDFilter
		}

		r, err := a.Deploy().GetV2Deployments(ctx, params)
		if err != nil {
			zap.L().Fatal("Error on request", zap.Error(err))
		}

		if r == nil || r.Body == nil {
			zap.L().Fatal("empty response from GetV2Deployments")
		}
		defer r.Body.Close()

		if r.StatusCode >= 400 {
			zap.L().Fatal("bad response from GetV2Deployments", zap.Int("statusCode", r.StatusCode))
		}

		var deployments []deploy.BodyDeploymentRead
		if err := json.NewDecoder(r.Body).Decode(&deployments); err != nil {
			zap.L().Fatal("failed to decode response body", zap.Error(err))
		}

		output := viper.GetString("output")
		switch output {
		case "json":
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			if err := enc.Encode(deployments); err != nil {
				zap.L().Fatal("failed to encode json", zap.Error(err))
			}
		case "yaml":
			data, err := yaml.Marshal(deployments)
			if err != nil {
				zap.L().Fatal("failed to marshal yaml", zap.Error(err))
			}
			fmt.Print(string(data))
		default:
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

			showOwner := true
			var firstOwner string
			for i, d := range deployments {
				if d.OwnerId == nil {
					showOwner = true
					break
				}
				if i == 0 {
					firstOwner = *d.OwnerId
					continue
				}
				if *d.OwnerId != firstOwner {
					showOwner = true
					break
				}
			}

			header := "ID\tName"
			if showOwner {
				header += "\tOwner"
			}
			header += "\tStatus"
			fmt.Fprintln(w, header)

			for _, d := range deployments {
				id := utils.DerefOrZero(d.Id)
				name := utils.DerefOrZero(d.Name)
				status := utils.DerefOrZero(d.Status)

				line := fmt.Sprintf("%s\t%s", id, name)
				if showOwner {
					line += "\t" + utils.DerefOrZero(d.OwnerId)
				}
				line += "\t" + status
				fmt.Fprintln(w, line)
			}

			w.Flush()
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.PersistentFlags().BoolP("all", "a", false, "Get all deployments")
	getCmd.PersistentFlags().BoolP("shared", "s", false, "Get shared deployments")
	getCmd.PersistentFlags().String("by-user-id", "", "Get all deployments by a userID")
	getCmd.PersistentFlags().StringP("output", "o", "table", "Output format: table, json, yaml")

	viper.BindPFlags(getCmd.PersistentFlags())
}
