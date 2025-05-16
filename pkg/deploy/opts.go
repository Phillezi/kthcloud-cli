package deploy

import (
	"github.com/Phillezi/kthcloud-cli/pkg/session"
	"github.com/go-resty/resty/v2"
)

type ClientOpts struct {
	BaseURL *string

	Session *session.Session

	Client *resty.Client
}
