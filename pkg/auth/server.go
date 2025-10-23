package auth

import (
	"fmt"
	"net/http"
	"sync"

	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type Server struct {
	port       string
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
		port:    "3000",
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
	return fmt.Sprintf("http://localhost:%s/login", s.port)
}
