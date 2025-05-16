package web

import (
	"context"
	"net/http"

	"github.com/Phillezi/kthcloud-cli/pkg/defaults"
	"github.com/Phillezi/kthcloud-cli/pkg/session"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
)

type Server struct {
	ctx context.Context

	address      string
	keycloakURL  string
	redirectHost string
	redirectPath string
	redirectURI  string

	sessionChannel chan *session.Session

	fetchOAuthToken func(redirectURI, code string) (*http.Response, error)

	cancelServer context.CancelFunc

	// state
	sent bool
}

func New(opts ...ServerOpts) *Server {
	var opt ServerOpts
	if len(opts) > 0 {
		opt = opts[0]
	}

	s := &Server{
		ctx: context.Background(),

		address:     util.PtrOr(opt.Address, ":3000"),
		keycloakURL: util.PtrOr(opt.KeycloakURL),

		redirectHost: util.PtrOr(opt.RedirectHost, defaults.DefaultRedirectSchemeHostPort),
		redirectPath: util.PtrOr(opt.RedirectPath, defaults.DefaultRedirectBasePath),

		sessionChannel:  util.PtrOr(opt.SessionChannel),
		fetchOAuthToken: opt.FetchOAuthToken,
	}

	s.redirectURI, _ = s.constructRedirectURI()

	return s
}

func (s *Server) WithContext(ctx context.Context) *Server {
	s.ctx = ctx
	return s
}
