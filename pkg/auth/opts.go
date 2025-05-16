package auth

import (
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/session"
	"github.com/go-resty/resty/v2"
)

type ClientOpts struct {
	KeycloakBaseURL      *string
	KeycloakClientID     *string
	KeycloakClientSecret *string
	KeycloakRealm        *string

	RedirectHost *string
	RedirectPath *string

	SessionPath *string

	RequestTimeout *time.Duration

	Session *session.Session

	Client *resty.Client
}
