package web

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/session"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/sirupsen/logrus"
)

//go:embed static/authenticated.html
var authenticatedHTML string

func (s *Server) Serve() error {
	if s.sessionChannel == nil {
		// log err here
		logrus.Errorln("sessionChannel is nil")
		return fmt.Errorf("sessionchannel is nil")
	}

	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Minute)
	defer cancel()

	s.cancelServer = cancel

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    s.address,
		Handler: mux,
	}

	s.setupRoutes(mux)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			logrus.Debug("recv cancel signal, shutting down server...")
			if err := server.Shutdown(s.ctx); err != nil {
				// log err here
				logrus.Errorln(err)
			}
			// log ctx cancellation here
			logrus.Debug("server cancelled")
		case <-s.ctx.Done():
			logrus.Debug("recv cancel signal, shutting down server...")
			if err := server.Shutdown(s.ctx); err != nil {
				// log err here
				logrus.Errorln(err)
			}
			// log ctx cancellation here
			logrus.Debug("server cancelled")
		}
	}()

	defer logrus.Infoln("server closed")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Errorf("Server failed: %s", err)
		if s.cancelServer != nil {
			s.once.Do(s.cancelServer)
		}
		go func() {
			select {
			case s.sessionChannel <- nil:
			default:
				logrus.Warn("session channel was full")
			}
		}()
		wg.Wait()
		return err
	}

	wg.Wait()
	return nil
}

func (s *Server) setupRoutes(mux *http.ServeMux) {
	donePath := "/auth/done"
	doneURL := s.redirectHost + donePath

	mux.HandleFunc("/", s.handleRoot)
	mux.HandleFunc(s.redirectPath, s.handleAuthRedirect(doneURL))
	mux.HandleFunc(donePath, s.handleAuthDone())
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, s.keycloakURL, http.StatusFound)
}

func (s *Server) handleAuthRedirect(doneURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Redirect(w, r, s.keycloakURL, http.StatusFound)
			return
		}

		resp, err := s.fetchOAuthToken(s.redirectURI, code)
		if err != nil {
			http.Error(w, "Failed to fetch OAuth token", http.StatusInternalServerError)
			//log.Println(err)
			// log err here

			logrus.Errorln(err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read response body", http.StatusInternalServerError)
			//log.Println(err)
			// log err here

			logrus.Errorln(err)
			return
		}

		jwt, err := util.ProcessResponse[session.JWTToken](string(body))
		if err != nil {
			http.Error(w, "Failed to process JWT token", http.StatusInternalServerError)
			//log.Println(err)
			// log err here

			logrus.Errorln(err)
			return
		}

		if jwt != nil {
			go func() {
				select {
				case s.sessionChannel <- session.New(*jwt):
					select {
					case <-s.authDoneVisited:
						logrus.Debugln("auth page was visited closing gracefully.")
					case <-time.After(10 * time.Second):
						logrus.Warnln("timeout reached before auth page was visited, forcing exit.")
					}
					if s.cancelServer != nil {
						s.once.Do(s.cancelServer)
					} else {
						logrus.Warnln("cancel server func is nil!")
					}
				default:
					logrus.Error("failed to send jwt, channel full")
				}
			}()
			http.Redirect(w, r, doneURL, http.StatusFound)
			return
		}

		// log err here
		logrus.Errorln("unexpected, jwt token was nil")
		http.Error(w, "jwt token was nil", http.StatusInternalServerError)
	}
}

func (s *Server) handleAuthDone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, authenticatedHTML)

		go func() {
			select {
			case s.authDoneVisited <- struct{}{}:
			default:
				logrus.Warn("authDoneVisited channel is full")
			}
		}()
	}
}
