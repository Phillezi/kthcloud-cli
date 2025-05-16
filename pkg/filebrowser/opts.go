package filebrowser

import (
	"context"
	"net/http"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/session"
)

type ClientOpts struct {
	Context *context.Context

	KeycloakBaseURL *string
	FilebrowserURL  *string

	RequestTimeout *time.Duration

	Session *session.Session

	Client *http.Client
}
