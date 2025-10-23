package browser

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"runtime"
)

var commonCandidates = map[string][]string{
	"code": {"--open-url"},
}

type ErrCommandNotFound struct {
	Cmds []string
}

func (e *ErrCommandNotFound) Error() string {
	return fmt.Sprintf("none of the browser open commands were found: %v", e.Cmds)
}

// Open opens the specified URL in the default browser.
func Open(url string, opts ...Option) error {
	var cmdCandidates []string
	var args []string

	switch runtime.GOOS {
	case "linux", "freebsd", "openbsd", "netbsd", "dragonfly":
		cmdCandidates = []string{"xdg-open", "gnome-open", "kde-open"}
		args = []string{url}
	case "darwin":
		cmdCandidates = []string{"open"}
		args = []string{url}
	case "windows":
		cmdCandidates = []string{"rundll32"}
		args = []string{"url.dll,FileProtocolHandler", url}
	default:
		return ErrUnsupportedPlatform
	}

	// Apply options
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	out := o.out
	if out == nil {
		out = io.Discard
	}

	var lastErr error
	tryCommand := func(cmd string, args ...string) bool {
		fmt.Fprintf(out, "[browser] would execute: %s %v\n", cmd, args)
		if o.dryRun {
			return true
		}

		if path, err := exec.LookPath(cmd); err == nil {
			ecmd := exec.Command(path, args...)
			ecmd.Stdout = out
			ecmd.Stderr = out
			if err := ecmd.Start(); err != nil {
				lastErr = fmt.Errorf("failed to start %s: %w", cmd, err)
				return false
			}
			return true
		}
		lastErr = fmt.Errorf("%s not found", cmd)
		return false
	}

	for _, cmd := range cmdCandidates {
		if tryCommand(cmd, args...) {
			return nil
		}
	}

	for cmd, extraArgs := range commonCandidates {
		if tryCommand(cmd, append(extraArgs, url)...) {
			return nil
		}
	}

	return errors.Join(&ErrCommandNotFound{Cmds: cmdCandidates}, lastErr)
}
