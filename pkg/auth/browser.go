package auth

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

func OpenBrowser(url string) error {
	var cmd *exec.Cmd
	switch {
	case runtime.GOOS == "linux":
		cmd = exec.Command("xdg-open", url)
	case runtime.GOOS == "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case runtime.GOOS == "darwin":
		cmd = exec.Command("open", url)
	default:
		return fmt.Errorf("unsupported platform")
	}
	fmt.Printf("Trying to open: %s in web browser", url)
	return cmd.Start()
}

func StartLocalServer() (string, error) {
	fmt.Println("Starting server")
	codeChannel := make(chan string)
	server := &http.Server{Addr: ":3000", Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code != "" {
			codeChannel <- code
			http.ServeFile(w, r, "static/authenticated.html")
			return
		}
		fmt.Fprintln(w, "Failed to get authorization code.")
	})}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Errorf("Server error: %v\n", err)
		}
	}()

	select {
	case code := <-codeChannel:
		return code, nil
	case <-time.After(5 * time.Minute):
		return "", fmt.Errorf("timeout waiting for authorization code")
	}
}
