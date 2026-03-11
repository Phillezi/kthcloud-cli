package auth

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/kthcloud/cli/internal/defaults"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type Server struct {
	addr       string
	oauth2Conf *oauth2.Config
	server     *http.Server

	closedMu  sync.RWMutex
	closeOnce sync.Once
	closed    bool
	tokenCh   chan *oauth2.Token

	l *zap.Logger
}

func NewServer(opts ...Option) *Server {
	s := &Server{
		addr:    defaults.DefaultLoginServerAddress,
		tokenCh: make(chan *oauth2.Token, 1),
		l:       zap.NewNop(),
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

func (s *Server) Token() <-chan *oauth2.Token {
	return s.tokenCh
}

func (s *Server) Url() string {
	return fmt.Sprintf("http://%s/login", s.addr)
}
