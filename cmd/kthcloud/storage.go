package main

import (
	"os"

	"github.com/kthcloud/cli/internal/app"
	"github.com/kthcloud/cli/internal/defaults"
	"github.com/kthcloud/cli/pkg/storage/filebrowser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var storageCmd = &cobra.Command{
	Use: "storage",
	Run: func(cmd *cobra.Command, args []string) {
		a := app.New(cmd.Context(), app.WithKeycloakOptions(
			viper.GetString("keycloak-client-id"),
			viper.GetString("keycloak-base-url"),
			viper.GetString("keycloak-realm"),
		),
			app.WithSessionKey(viper.GetString("session-key")),
			app.WithLogger(zap.L()),
		)

		c := filebrowser.NewFileBrowserClient(
			defaults.DefaultSMProxyAPIBaseURL,
			"",
			filebrowser.WithLogger(zap.L().Named("storage")),
			filebrowser.WithMiddleware(a.SessionMiddleware()),
		)

		f, _ := os.Open("test.txt")
		defer f.Close()

		progress := make(chan filebrowser.UploadProgress)
		go func() {
			for p := range progress {
				zap.L().Info("progress", zap.Int64("uploaded", p.Uploaded))
			}
		}()

		if err := c.UploadFileWithProgress(cmd.Context(), "v2/test.txt", f, progress); err != nil {
			zap.L().Error("upload failed", zap.Error(err))
		}
	},
}

func init() {
	rootCmd.AddCommand(storageCmd)
}
