package auth

import (
	"context"
	"net/http"
	"time"
)

func (s *Server) Start() error {
	s.l.Debug("server.Start")
	defer s.l.Debug("server.Start exited")

	mux := http.NewServeMux()
	mux.HandleFunc("/login", s.LoginHandler)
	mux.HandleFunc("/callback", s.CallbackHandler)

	s.server = &http.Server{
		Addr:         ":" + s.port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	s.l.Sugar().Infof("Server started on %s", s.Url())

	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.l.Debug("server.Shutdown")
	defer s.l.Debug("server.Shutdown exited")

	if s.server != nil {
		return s.server.Shutdown(ctx)
	}

	s.closeOnce.Do(func() {
		s.closedMu.Lock()
		defer s.closedMu.Unlock()
		close(s.tokenCh)
		s.closed = true
	})

	return nil
}
