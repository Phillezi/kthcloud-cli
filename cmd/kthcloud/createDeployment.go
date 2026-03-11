package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/kthcloud/cli/internal/app"
	"github.com/kthcloud/cli/internal/body"
	"github.com/kthcloud/cli/internal/constants"
	"github.com/kthcloud/cli/pkg/deploy"
	"github.com/kthcloud/cli/pkg/parser"
	"github.com/kthcloud/cli/pkg/session"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var flags parser.DeploymentFlags

var createDeploymentCmd = &cobra.Command{
	Use: "deployment",
	Aliases: []string{
		"deployments",
	},
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			return
		}

		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
		defer cancel()

		a := app.New(ctx, app.WithKeycloakOptions(
			viper.GetString(constants.ViperKeycloakClientId),
			viper.GetString(constants.ViperKeycloakBaseURL),
			viper.GetString(constants.ViperKeycloakRealm),
		),
			app.WithSessionKey(viper.GetString(constants.ViperSessionKey)),
			app.WithAPITokenSession(viper.GetString(constants.ViperDeployAPIToken)),
			app.WithLogger(zap.L()),
		)

		flags.Args = args[1:]
		reqBody, err := parser.ParseDeployment(args[0], flags)
		if err != nil {
			zap.L().Fatal("Error parsing args", zap.Error(err))
		}

		reader, err := body.Reader(reqBody)
		if err != nil {
			zap.L().Fatal("Error marshal", zap.Error(err))
		}

		r, err := a.Deploy().PostV2DeploymentsWithBodyWithResponse(ctx, "application/json", reader)
		if err != nil {
			if errors.Is(err, session.ErrLoginRequired) {
				zap.L().Fatal("Login is required, please run the login command")
			}
			zap.L().Fatal("Error on request", zap.Error(err))
		}

		obj, err := deploy.HandleAndAssert[*deploy.BodyDeploymentRead](r, "create")
		if err != nil {
			zap.L().Fatal("Error on handle", zap.Error(err))
		}

		if obj != nil {
			if obj.Id != nil {
				_, _ = fmt.Fprintln(os.Stdout, *obj.Id)
			}
		}
	},
}

func init() {
	createCmd.AddCommand(createDeploymentCmd)

	createDeploymentCmd.Flags().StringVar(&flags.Name, "name", "", "Deployment name")
	_ = createDeploymentCmd.MarkFlagRequired("name")
	createDeploymentCmd.Flags().StringSliceVarP(&flags.Env, "env", "e", nil, "Set environment variables (KEY=value)")
	createDeploymentCmd.Flags().StringSliceVarP(&flags.Volume, "volume", "v", nil, "Bind mount a volume (local:remote)")
	createDeploymentCmd.Flags().StringSliceVarP(&flags.Port, "publish", "p", nil, "Publish a port (local:remote)")
	createDeploymentCmd.Flags().Float32Var(&flags.CPU, "cpu", 0, "CPU cores")
	createDeploymentCmd.Flags().Float32Var(&flags.RAM, "ram", 0, "RAM (in GB)")
	createDeploymentCmd.Flags().StringSliceVar(&flags.GPU, "gpu", nil, "GPU (<claim-name>/<gpu-name>)")
	createDeploymentCmd.Flags().IntVar(&flags.Replicas, "replicas", 1, "Number of replicas")
	createDeploymentCmd.Flags().StringVar(&flags.Zone, "zone", "", "Deployment zone")
	createDeploymentCmd.Flags().StringVar(&flags.Domain, "domain", "", "Custom domain")
	createDeploymentCmd.Flags().StringVar(&flags.Health, "health", "", "Health check path")
	createDeploymentCmd.Flags().StringVar(&flags.Visibility, "visibility", "", "Visibility (public, private, auth)")
	createDeploymentCmd.Flags().BoolVar(&flags.NeverStale, "never-stale", false, "Prevent auto-disable")
}
