package options

import (
	"fmt"
	"net/url"

	"github.com/Phillezi/kthcloud-cli/internal/interrupt"
	"github.com/Phillezi/kthcloud-cli/pkg/auth"
	"github.com/Phillezi/kthcloud-cli/pkg/deploy"
	"github.com/Phillezi/kthcloud-cli/pkg/filebrowser"
	"github.com/spf13/viper"
)

func DeployOpts() deploy.ClientOpts {
	baseURL := viper.GetString("api-url")

	return deploy.ClientOpts{
		BaseURL: &baseURL,
	}
}

func AuthOpts() auth.ClientOpts {
	keycloakBaseURL := viper.GetString("keycloak-host")
	keycloakClientID := viper.GetString("client-id")
	keycloakClientSecret := viper.GetString("client-secret")
	keycloakRealm := viper.GetString("keycloak-realm")

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
	}
}

func FilebrowserOpts() filebrowser.ClientOpts {
	keycloakBaseURL := viper.GetString("keycloak-host")

	requestTimeout := viper.GetDuration("request-timeout")

	return filebrowser.ClientOpts{
		KeycloakBaseURL: &keycloakBaseURL,
		RequestTimeout:  &requestTimeout,
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
