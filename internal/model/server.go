package model

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

type Server struct {
	Addr              string
	codeChannel       chan *AuthSession
	authenticateHTML  string
	authenticatedHTML string
}

func NewServer(addr string, authenticateHTML string, authenticatedHTML string) *Server {
	authHTML, err := NewTemplate(viper.GetString("keycloak-host"), viper.GetString("keycloak-realm"), viper.GetString("client-id")).Replace(authenticateHTML)
	if err != nil {
		log.Fatal("could not replace html template variables")
		return nil
	}
	return &Server{
		Addr:              addr,
		codeChannel:       make(chan *AuthSession),
		authenticateHTML:  authHTML,
		authenticatedHTML: authenticatedHTML,
	}
}

func (s *Server) Start() (*AuthSession, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleRoot)
	mux.HandleFunc("/auth", s.handleAuth)

	server := &http.Server{
		Addr:    s.Addr,
		Handler: mux,
	}

	go func() {
		log.Infof("Starting server on %s\n", s.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v\n", err)
		}
	}()

	select {
	case code := <-s.codeChannel:
		return code, nil
	case <-time.After(5 * time.Minute):
		return nil, fmt.Errorf("timeout waiting for authorization code")
	}
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, s.authenticateHTML)
}

func (s *Server) handleAuth(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	token := r.URL.Query().Get("token")
	expires_in, err := strconv.Atoi(r.URL.Query().Get("expires_in"))
	if err != nil {
		expires_in = 0
	}
	refresh_token := r.URL.Query().Get("refresh_token")
	now := time.Now()
	if token != "" && id != "" {
		s.codeChannel <- &AuthSession{
			SessionStart: &now,
			Token:        token,
			RefreshToken: refresh_token,
			ExpiresIn:    expires_in,
			SessionId:    &id,
		}
		fmt.Fprintln(w, s.authenticatedHTML)
		return
	}
	fmt.Fprintln(w, "Failed to get authorization code.")
}
