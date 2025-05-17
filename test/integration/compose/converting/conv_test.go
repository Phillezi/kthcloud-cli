package converting

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Phillezi/kthcloud-cli/pkg/models/compose"
	"github.com/Phillezi/kthcloud-cli/pkg/models/service"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
)

func TestComposeConversion(t *testing.T) {
	comp := &compose.Compose{
		Services: map[string]*service.Service{
			"app": {
				Image: "nginx:latest",
				Environment: service.EnvVars{
					"foo": "bar",
				},
				Ports: []string{
					"8080:8080",
					"9090:9090",
				},
				Volumes: []string{
					"./data:/app/data",
				},
				Command: []string{
					"bash -c",
					"\"echo 'hello world'\"",
				},
				Dependencies: []string{
					"otherapp",
				},
			},
			"otherapp": {
				Image: "nginx:latest",
				Environment: service.EnvVars{
					"foo":            "bar",
					"KTHCLOUD_CORES": "99",
					"KTHCLOUD_RAM":   "99",
				},
				Ports: []string{
					"8080:8080",
					"9090:9090",
				},
				Volumes: []string{
					"./data:/app/data",
				},
				Command: []string{
					"bash -c",
					"\"echo 'hello world'\"",
				},
			},
		},
	}

	depls, deps := comp.ToDeploymentsWDeps()

	t.Run("dependencies match", func(t *testing.T) {
		if !util.Contains(deps["app"], "otherapp") || len(deps["app"]) != 1 {
			t.Error("dependencies does not match")
		}
	})

	t.Run("service conversion", func(t *testing.T) {
		for _, depl := range depls {
			associatedService, ok := comp.Services[depl.Name]
			if !ok {
				t.Error("could not find associated service")
				continue
			}

			if depl.Image == nil && strings.TrimSpace(associatedService.Image) != "" || !strings.EqualFold(*depl.Image, associatedService.Image) {
				t.Error("images dont match")
			}

			if len(depl.Args) != len(associatedService.Command) {
				t.Errorf("init command count mismatch for deployment %s: expected %d, got %d", depl.Name, len(associatedService.Command), len(depl.Args))
			} else {
				for i, cmd := range depl.Args {
					if cmd != associatedService.Command[i] {
						t.Errorf("init command mismatch for deployment %s at index %d: expected %s, got %s", depl.Name, i, associatedService.Command[i], cmd)
					}
				}
			}

			if val, ok := associatedService.Environment["KTHCLOUD_CORES"]; ok && depl.CpuCores != nil {
				if fmt.Sprintf("%v", *depl.CpuCores) != val {
					t.Errorf("CPU cores mismatch for deployment %s: expected %s, got %v", depl.Name, val, *depl.CpuCores)
				}
			}

			if val, ok := associatedService.Environment["KTHCLOUD_RAM"]; ok && depl.RAM != nil {
				if fmt.Sprintf("%v", *depl.RAM) != val {
					t.Errorf("RAM mismatch for deployment %s: expected %s, got %v", depl.Name, val, *depl.RAM)
				}
			}
		}
	})

}
