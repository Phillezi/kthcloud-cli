package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/kthcloud/cli/internal/app"
	"github.com/kthcloud/cli/internal/constants"
	"github.com/kthcloud/cli/pkg/deploy"
	"github.com/kthcloud/cli/pkg/session"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var whoamiCmd = &cobra.Command{
	Use: "whoami",
	Run: func(cmd *cobra.Command, args []string) {
		a := app.FromViper(cmd.Context(), zap.L(), *viper.GetViper())

		mgr, ok := a.Session().(*session.DefaultManager)
		if !ok {
			zap.L().Fatal("Unsupported session type")
		}

		sess, err := mgr.GetSession(viper.GetString(constants.ViperSessionKey))
		if err != nil {
			zap.L().Fatal("Could not get session", zap.Error(err))
		}

		sub, err := getSubFromJWT(sess.Token.AccessToken)
		if err != nil {
			zap.L().Fatal("Could not extract sub claim from token", zap.Error(err))
		}

		resp, err := a.Deploy().GetV2UsersUserIdWithResponse(cmd.Context(), sub, &deploy.GetV2UsersUserIdParams{})
		if err != nil {
			zap.L().Fatal("Error getting user info", zap.Error(err))
		}

		obj, err := deploy.HandleAndAssert[*deploy.BodyUserRead](resp, "get")
		if err != nil {
			zap.L().Fatal("Error on handle", zap.Error(err))
		}

		dat, err := yaml.Marshal(obj)
		if err != nil {
			zap.L().Fatal("Error on marshal", zap.Error(err))
		}

		fmt.Fprintln(os.Stderr, string(dat))
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}

func getSubFromJWT(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid JWT")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	var claims struct {
		Sub string `json:"sub"`
	}

	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", err
	}

	return claims.Sub, nil
}
