package ps

import (
	"fmt"
	"go-deploy/dto/v2/body"
	"os"
	"strings"

	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/sirupsen/logrus"
)

func List(all bool) {
	c := client.Get()

	vms, err := c.Vms()
	if err != nil {
		logrus.Fatal("could not get vms:", err)
	}

	renderVmsTable(vms, all)
}

func getGPUName(gpuNames *map[string]string, groupID string) (string, error) {
	if name, ok := (*gpuNames)[groupID]; ok {
		return name, nil
	}
	c := client.Get()

	group, err := c.GpuGroupByID(groupID)
	if err != nil {
		return "", err
	}

	return group.DisplayName, nil
}

func renderVmsTable(vms []body.VmRead, all bool) {

	gpuNames := make(map[string]string)

	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Name", "Status", "GPU", "Visit"})

	for _, vm := range vms {
		if vm.Status == "resourceStopped" && !all {
			continue
		}

		vm.Status = strings.Replace(vm.Status, "resource", "", 1)

		gpu := ""
		if vm.GPU != nil {
			logrus.Infoln(*vm.GPU)
			gpu, _ = getGPUName(&gpuNames, vm.GPU.GpuGroupID)
		}

		var visitPorts []string
		for _, port := range vm.Ports {
			if port.HttpProxy != nil && port.HttpProxy.URL != nil {
				if port.ExternalPort != nil {
					visitPorts = append(visitPorts, fmt.Sprintf("\u001b]8;;%s\u0007%s\u001b]8;;\u0007 (%d:%d)", *port.HttpProxy.URL, port.Name, *port.ExternalPort, port.Port))
				} else {
					visitPorts = append(visitPorts, fmt.Sprintf("\u001b]8;;%s\u0007%s\u001b]8;;\u0007", *port.HttpProxy.URL, port.Name))
				}
			} else if port.ExternalPort != nil {
				visitPorts = append(visitPorts, fmt.Sprintf("\u001b]8;;%s:%d\u0007%s\u001b]8;;\u0007 (%d:%d)", "deploy.cloud.cbh.kth.se", *port.ExternalPort, port.Name, *port.ExternalPort, port.Port))
			} else {
				visitPorts = append(visitPorts, fmt.Sprintf("%s (NONE:%d)", port.Name, port.Port))
			}
		}

		var rows []table.Row
		rows = append(rows, table.Row{vm.ID, vm.Name, vm.Status, gpu, (func() string {
			if len(visitPorts) > 0 {
				return visitPorts[0]
			}
			return ""
		})()})

		for i := 1; i < len(visitPorts); i++ {
			rows = append(rows, table.Row{"", "", "", "", visitPorts[i]})
		}

		t.AppendRows(rows)
		t.AppendSeparator()
	}

	t.Render()
}
