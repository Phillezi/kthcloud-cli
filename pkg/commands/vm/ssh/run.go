package ssh

import (
	"fmt"
	"strings"

	"github.com/kthcloud/go-deploy/dto/v2/body"
	"github.com/sirupsen/logrus"
)

func (c *Command) Run() error {
	if c.client == nil {
		return fmt.Errorf("client is nil")
	}

	vms, err := c.client.Vms()
	if err != nil {
		logrus.Error("could not get vms:", err)
		return err
	}

	var vm *body.VmRead
	if c.id != "" || c.name != "" {
		vm, err = selectNonInteractive(vms, c.id, c.name)
		if err != nil {
			return err
		}
	} else {
		vm, err = selectVm(vms)
		if err != nil {
			return err
		}
	}

	if vm == nil {
		return fmt.Errorf("no vm selected")
	}

	if vm.SshConnectionString == nil {
		return fmt.Errorf("vm does not have an ssh connectionstring")
	}

	connstr := strings.Split((*vm.SshConnectionString), " ")[1:]
	if len(connstr) != 3 {
		return fmt.Errorf("unexpected connectionstring format")
	}

	cmd, err := ssh(connstr[0], connstr[2])
	if err != nil {
		return fmt.Errorf("failed to create ssh command: %v", err)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute ssh command: %v", err)
	}

	return nil
}
