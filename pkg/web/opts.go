package web

import (
	"net/http"

	"github.com/Phillezi/kthcloud-cli/pkg/session"
)

type ServerOpts struct {
	Address      *string
	KeycloakURL  *string
	RedirectHost *string
	RedirectPath *string

	SessionChannel *chan *session.Session

	FetchOAuthToken func(redirectURI, code string) (*http.Response, error)
}
