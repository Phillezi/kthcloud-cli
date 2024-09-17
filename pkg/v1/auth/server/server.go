package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/session"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/token"
)

type Server struct {
	addr           string
	sessionChannel chan *session.Session
	//cookies         []*http.Cookie
	kcURL           string
	fetchOAuthToken func(redirectURI, code string) (*http.Response, error)
}

func New(
	addr string,
	kcURL string,
	sessionChannel chan *session.Session,
	//cookies []*http.Cookie,
	fetchOAuthToken func(redirectURI, code string) (*http.Response, error),
) *Server {
	return &Server{
		addr:           addr,
		kcURL:          kcURL,
		sessionChannel: sessionChannel,
		//cookies:         cookies,
		fetchOAuthToken: fetchOAuthToken,
	}
}

func (s *Server) Start() {
	go func() {
		server := &http.Server{Addr: s.addr}
		var sess *session.Session

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, s.kcURL, http.StatusFound)
		})

		http.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			if code == "" {
				//http.Redirect(w, r, newURL, http.StatusFound)
				fmt.Fprintln(w, "no code provided")
				http.Redirect(w, r, s.kcURL, http.StatusFound)
				return
			}

			for _, cookie := range r.Cookies() {
				fmt.Println(cookie)
				http.SetCookie(w, cookie)
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
			for _, cookie := range resp.Cookies() {
				fmt.Println("got cookie: ", cookie)
			}

			//fmt.Fprintln(w, "code: ", code)
			http.Redirect(w, r, "http://localhost:3000/auth/done", http.StatusFound)
		})

		// Handle /auth/callback
		http.HandleFunc("/auth/done", func(w http.ResponseWriter, r *http.Request) {

			s.sessionChannel <- sess
			fmt.Fprintln(w, "Callback received. Server will now shut down.")

			// Close the server
			go func() {
				time.Sleep(500 * time.Millisecond)
				if err := server.Shutdown(context.Background()); err != nil {
					log.Fatalf("Server Shutdown Failed:%+v", err)
				}
				fmt.Println("Server stopped after serving the callback request")
			}()
		})

		// Start the server on localhost:3000
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s", err)
		}
	}()
}
