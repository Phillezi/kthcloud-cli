package auth

import (
	"fmt"
	"strings"

	"github.com/kthcloud/cli/pkg/keycloak"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type Option func(*Server)

func WithPort(port string) Option {
	return func(s *Server) {
		oldPort := s.port
		s.port = port
		if s.oauth2Conf != nil {
			s.oauth2Conf.RedirectURL = strings.ReplaceAll(s.oauth2Conf.RedirectURL, oldPort, port)
		}
	}
}

func WithOAuth2Config(conf *oauth2.Config) Option {
	return func(s *Server) {
		s.oauth2Conf = conf
	}
}

func WithKeycloakOAuth2Config(clientID, baseURL, realm string) Option {
	return func(s *Server) {
		redirectURL := fmt.Sprintf("http://localhost:%s/callback", s.port)
		s.oauth2Conf = keycloak.Config(clientID, baseURL, redirectURL, realm)
	}
}

func WithLogger(l *zap.Logger) Option {
	return func(s *Server) {
		s.l = l
	}
}
