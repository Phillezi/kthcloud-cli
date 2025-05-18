package converting

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Phillezi/kthcloud-cli/pkg/convert"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/compose-spec/compose-go/v2/types"
	"github.com/kthcloud/go-deploy/utils"
)

func TestComposeConversion(t *testing.T) {
	comp := &types.Project{
		Services: types.Services{
			"app": {
				Image: "nginx:latest",
				Environment: types.MappingWithEquals{
					"foo": utils.PtrOf("bar"),
				},
				Ports: []types.ServicePortConfig{
					{Target: 8080},
					{Target: 9090},
				},
				Volumes: []types.ServiceVolumeConfig{
					{Source: "./data", Target: "/app/data"},
				},
				Command: types.ShellCommand{
					"bash -c",
					"\"echo 'hello world'\"",
				},
				DependsOn: types.DependsOnConfig{
					"otherapp": types.ServiceDependency{},
				},
			},
			"otherapp": {
				Image: "nginx:latest",
				Environment: types.MappingWithEquals{
					"foo":            utils.PtrOf("bar"),
					"KTHCLOUD_CORES": utils.PtrOf("99"),
					"KTHCLOUD_RAM":   utils.PtrOf("99"),
				},
				Ports: []types.ServicePortConfig{
					{Target: 8080},
					{Target: 9090},
				},
				Volumes: []types.ServiceVolumeConfig{
					{Source: "./data", Target: "/app/data"},
				},
				Command: types.ShellCommand{
					"bash -c",
					"\"echo 'hello world'\"",
				},
			},
		},
	}

	var out convert.Wrap
	if err := convert.ToCloud(comp, &out); err != nil {
		t.Error(err)
		return
	}

	t.Run("dependencies match", func(t *testing.T) {
		if !util.Contains(out.Dependencies["app"], "otherapp") || len(out.Dependencies["app"]) != 1 {
			t.Error("dependencies does not match")
		}
	})

	t.Run("service conversion", func(t *testing.T) {
		for _, depl := range out.Deployments {
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
				if fmt.Sprintf("%v", *depl.CpuCores) != *val {
					t.Errorf("CPU cores mismatch for deployment %s: expected %s, got %v", depl.Name, *val, *depl.CpuCores)
				}
			}

			if val, ok := associatedService.Environment["KTHCLOUD_RAM"]; ok && depl.RAM != nil {
				if fmt.Sprintf("%v", *depl.RAM) != *val {
					t.Errorf("RAM mismatch for deployment %s: expected %s, got %v", depl.Name, *val, *depl.RAM)
				}
			}
		}
	})

}
