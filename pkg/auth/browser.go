package auth

import (
	_ "embed"
	"fmt"
	"kthcloud-cli/internal/model"
	"os/exec"
	"runtime"
	"time"

	"github.com/briandowns/spinner"
)

//go:embed static/authenticated.html
var authenticatedHTML string

//go:embed static/authenticate.html
var authenticateHTML string

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
	fmt.Printf("Trying to open: %s in web browser\n\n", url)
	return cmd.Start()
}

func StartLocalServer() (*model.AuthSession, error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Color("blue")
	s.Prefix = "Waiting for login\n"
	s.Start()
	defer s.Stop()
	server := model.NewServer(":3000", authenticateHTML, authenticatedHTML)
	return server.Start()
}
