package server

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/session"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/token"
)

//go:embed static/authenticated.html
var authenticatedHTML string

type Server struct {
	addr            string
	sessionChannel  chan *session.Session
	kcURL           string
	fetchOAuthToken func(redirectURI, code string) (*http.Response, error)
	ctx             context.Context
}

func New(
	addr string,
	kcURL string,
	sessionChannel chan *session.Session,
	fetchOAuthToken func(redirectURI, code string) (*http.Response, error),
	ctx context.Context,
) *Server {
	return &Server{
		addr:            addr,
		kcURL:           kcURL,
		sessionChannel:  sessionChannel,
		fetchOAuthToken: fetchOAuthToken,
		ctx:             ctx,
	}
}

func (s *Server) Start() {
	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Minute)

	go func() {
		defer cancel()

		server := &http.Server{Addr: s.addr}
		var sess *session.Session
		doneChan := make(chan bool)

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, s.kcURL, http.StatusFound)
		})

		http.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			if code == "" {
				fmt.Fprintln(w, "no code provided")
				http.Redirect(w, r, s.kcURL, http.StatusFound)
				return
			}

			resp, err := s.fetchOAuthToken("http://localhost:3000/auth/callback", code)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error reading response body: %v\n", err)
				return
			}
			jwt, err := util.ProcessResponse[token.JWTToken](string(body))
			if err != nil {
				fmt.Println(err)
				return
			}
			sess = session.New(*jwt)

			http.Redirect(w, r, "http://localhost:3000/auth/done", http.StatusFound)
		})

		http.HandleFunc("/auth/done", func(w http.ResponseWriter, r *http.Request) {
			doneChan <- true
			s.sessionChannel <- sess
			fmt.Fprintln(w, authenticatedHTML)

			go func() {
				time.Sleep(500 * time.Millisecond)
				if err := server.Shutdown(context.Background()); err != nil {
					log.Fatalf("Server Shutdown Failed:%+v", err)
				}
				fmt.Println("Server stopped after serving the callback request")
			}()
		})

		go func() {
			select {
			case <-ctx.Done():
				s.sessionChannel <- nil
				if err := server.Shutdown(context.Background()); err != nil {
					log.Fatalf("Server Shutdown Failed:%+v", err)
				}
				fmt.Println("Server stopped after timeout of 3 minutes")
			case <-doneChan:
				return
			}
		}()

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s", err)
		}
	}()
}
