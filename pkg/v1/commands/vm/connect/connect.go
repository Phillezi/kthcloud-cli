package connect

import (
	"go-deploy/dto/v2/body"
	"strings"

	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	"github.com/sirupsen/logrus"
)

func SSH(id string, name string) {
	c := client.Get()

	vms, err := c.Vms()
	if err != nil {
		logrus.Fatal("could not get vms:", err)
	}

	var vm *body.VmRead

	if id != "" || name != "" {
		vm, err = selectNonInteractive(vms, id, name)
		if err != nil {
			logrus.Fatalln(err)
		}
	} else {
		vm, err = selectVm(vms)
		if err != nil {
			logrus.Fatalln(err)
		}
	}

	if vm == nil {
		logrus.Fatalln("No vm selected")
	}

	if vm.SshConnectionString == nil {
		logrus.Fatalln("VM does not have a SSH connectionstring")
	}

	connstr := strings.Split((*vm.SshConnectionString), " ")[1:]
	if len(connstr) != 3 {
		logrus.Fatal("unexpected connectionstring format")
	}

	cmd, err := ssh(connstr[0], connstr[2])
	if err != nil {
		logrus.Fatalf("failed to create SSH command: %v", err)
	}

	if err := cmd.Run(); err != nil {
		logrus.Fatalf("failed to execute SSH command: %v", err)
	}
}
