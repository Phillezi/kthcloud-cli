package options

import (
	"fmt"
	"net/url"
	"path"
	"sync"

	"github.com/Phillezi/kthcloud-cli/internal/interrupt"
	"github.com/Phillezi/kthcloud-cli/pkg/auth"
	"github.com/Phillezi/kthcloud-cli/pkg/config"
	"github.com/Phillezi/kthcloud-cli/pkg/defaults"
	"github.com/Phillezi/kthcloud-cli/pkg/deploy"
	"github.com/Phillezi/kthcloud-cli/pkg/filebrowser"
	"github.com/Phillezi/kthcloud-cli/pkg/session"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/spf13/viper"
)

var (
	once       sync.Once
	deployOpts *deploy.ClientOpts
)

func DeployOpts() deploy.ClientOpts {
	once.Do(func() {
		baseURL := util.Or(viper.GetString("api-url"), defaults.DefaultDeployBaseURL)
		sess, _ := session.Load(util.Or(viper.GetString("session-path"), path.Join(config.GetConfigPath(), "session.json")))

		deployOpts = &deploy.ClientOpts{
			BaseURL: &baseURL,
			Session: sess,
		}
	})
	return *deployOpts
}

func AuthOpts() auth.ClientOpts {
	keycloakBaseURL := util.Or(viper.GetString("keycloak-host"), defaults.DefaultKeycloakBaseURL)
	keycloakClientID := util.Or(viper.GetString("client-id"), defaults.DefaultKeycloakClientID)
	keycloakClientSecret := util.Or(viper.GetString("client-secret"), defaults.DefaultKeycloakClientSecret)
	keycloakRealm := util.Or(viper.GetString("keycloak-realm"), defaults.DefaultKeycloakRealm)

	redirectURI := viper.GetString("redirect-uri")
	var (
		redirectHost string
		redirectPath string
	)

	if url, _ := url.Parse(redirectURI); url != nil {
		redirectHost = fmt.Sprintf("%s://%s", url.Scheme, url.Host)
		redirectPath = url.Path
	}

	sessionPath := viper.GetString("session-path")

	requestTimeout := viper.GetDuration("request-timeout")

	return auth.ClientOpts{
		KeycloakBaseURL:      &keycloakBaseURL,
		KeycloakClientID:     &keycloakClientID,
		KeycloakClientSecret: &keycloakClientSecret,
		KeycloakRealm:        &keycloakRealm,

		RedirectHost: &redirectHost,
		RedirectPath: &redirectPath,

		SessionPath: &sessionPath,

		RequestTimeout: &requestTimeout,
		Session:        DeployOpts().Session,
	}
}

func FilebrowserOpts() filebrowser.ClientOpts {
	keycloakBaseURL := viper.GetString("keycloak-host")
	filebrowserURL := util.Or(viper.GetString("filebrowser-url"), defaults.DefaultStorageManagerProxy)

	requestTimeout := viper.GetDuration("request-timeout")

	return filebrowser.ClientOpts{
		KeycloakBaseURL: &keycloakBaseURL,
		RequestTimeout:  &requestTimeout,
		Session:         DeployOpts().Session,
		FilebrowserURL:  &filebrowserURL,
	}
}

func DefaultClient() *deploy.Client {
	return deploy.GetInstance(
		DeployOpts(),
	).WithContext(
		interrupt.GetInstance().Context(),
	).WithAuthClient(
		AuthOpts(),
	).WithStorageClient(
		FilebrowserOpts(),
	)
}
