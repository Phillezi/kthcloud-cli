package auth

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func (s *Server) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		s.l.Warn("missing code in callback")
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}

	tok, err := s.oauth2Conf.Exchange(r.Context(), code)
	if err != nil {
		s.l.Error("token exchange failed", zap.Error(err))
		http.Error(w, "token exchange failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	s.closedMu.RLock()

	if s.closed {
		s.l.Warn("callback received but server already closed")
		writeHTML(w, "An error occurred. You can close this window.")
		s.closedMu.RUnlock()
		return
	}

	select {
	case s.tokenCh <- tok:
		s.closedMu.RUnlock()
		s.l.Info("token received from callback")
	default:
		s.closedMu.RUnlock()
		s.l.Warn("token channel full or closed")
		writeHTML(w, "An error occurred. You can close this window.")
		return
	}

	writeHTMLWithScript(w, "Login complete. You can close this window.")

	go func(s *Server) { s.Shutdown(context.Background()) }(s)
}

// writeHTML writes a simple HTML page with a message
func writeHTML(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Login Complete</title>
<style>
  body {
    display:flex;
    justify-content:center;
    align-items:center;
    height:100vh;
    margin:0;
    font-family:sans-serif;
    background-color:#fff;
    color:#000;
  }
  @media(prefers-color-scheme:dark){
    body{background-color:#121212;color:#e0e0e0}
  }
</style>
</head>
<body>
  <h1>%s</h1>
</body>
</html>`, message)
}

// writeHTMLWithScript writes HTML with a message and window.close() script
func writeHTMLWithScript(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Login Complete</title>
<style>
  body {
    display:flex;
    justify-content:center;
    align-items:center;
    height:100vh;
    margin:0;
    font-family:sans-serif;
    background-color:#fff;
    color:#000;
  }
  @media(prefers-color-scheme:dark){
    body{background-color:#121212;color:#e0e0e0}
  }
</style>
</head>
<body>
  <h1>%s</h1>
<script>window.close();</script>
</body>
</html>`, message)
}
