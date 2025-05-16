package ssh

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func ssh(host string, port string) (*exec.Cmd, error) {
	if host == "" {
		return nil, fmt.Errorf("host cannot be empty")
	}

	if port == "" {
		port = "22"
	}

	logrus.Debugln("ssh command called with addr:", host, "port:", port)
	cmd := exec.Command("ssh", host, "-p", port)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd, nil
}
