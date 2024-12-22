package connect

import (
	"fmt"
	"go-deploy/dto/v2/body"

	"github.com/manifoldco/promptui"
)

func selectVm(vms []body.VmRead) (*body.VmRead, error) {
	if len(vms) < 1 {
		return nil, fmt.Errorf("no VMs available")
	}

	startIndex := 0

	availableVMs := []*body.VmRead{}
	vmItems := make([]string, len(vms))
	for i, vm := range vms {
		if vm.SshConnectionString != nil &&
			vm.Status == "resourceRunning" {
			vmItems[i] = fmt.Sprintf("%s", vm.Name)
			availableVMs = append(availableVMs, &vm)
		}
	}

	if len(availableVMs) < 1 {
		return nil, fmt.Errorf("no running VMs")
	}

	if len(availableVMs) == 1 {
		return availableVMs[0], nil
	}

	prompt := promptui.Select{
		Label:     "Select a VM\n",
		Items:     vmItems,
		CursorPos: startIndex,
	}
	index, _, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to select a VM: %v", err)
	}

	return availableVMs[index], nil
}

func selectNonInteractive(vms []body.VmRead, id string, name string) (*body.VmRead, error) {
	if len(vms) == 0 {
		return nil, fmt.Errorf("no VMs available")
	}

	if id != "" {
		for i, vm := range vms {
			if vm.ID == id {
				return &vms[i], nil
			}
		}
	} else if name != "" {
		for i, vm := range vms {
			if vm.Name == name {
				return &vms[i], nil
			}
		}
	} else {
		return nil, fmt.Errorf("id or name needs to be provided to select non-interactively")
	}

	return nil, fmt.Errorf("VM not found")
}
