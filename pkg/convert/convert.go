package convert

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/compose-spec/compose-go/v2/types"
	"github.com/kthcloud/go-deploy/dto/v2/body"
	"github.com/sirupsen/logrus"
)

func ToCloud(in *types.Project, out *Wrap) error {
	if in == nil {
		return fmt.Errorf("input is nil")
	}
	if out == nil {
		return fmt.Errorf("out is nil")
	}
	if out.Dependencies == nil {
		out.Dependencies = make(map[string][]string)
	}

	if out.Deployments == nil {
		out.Deployments = make([]body.DeploymentCreate, 0)
	}

	for name, service := range in.Services {
		depl, deps, err := ServiceToDeployment(&service, ServiceToDeploymentAdditionalInfo{Name: name, Hash: HashServices(in.Services), CWD: in.WorkingDir})
		if err != nil {
			return err
		}
		out.Deployments = append(out.Deployments, *depl)
		if len(deps) > 0 {
			out.Dependencies[depl.Name] = deps
		}
	}

	out.Source = in

	return nil
}

type ServiceToDeploymentAdditionalInfo struct {
	Name string
	Hash string
	CWD  string
}

func ServiceToDeployment(in *types.ServiceConfig, additionalInfo ...ServiceToDeploymentAdditionalInfo) (out *body.DeploymentCreate, dependencies []string, err error) {
	var extra ServiceToDeploymentAdditionalInfo
	if len(additionalInfo) > 0 {
		extra = additionalInfo[0]
	}

	out = &body.DeploymentCreate{}

	out.Name = util.Or(in.ContainerName, extra.Name)
	if out.Name == "" {
		return nil, nil, fmt.Errorf("service needs a name")
	}

	if in.Build != nil {
		if in.Image != "" {
			return nil, nil, fmt.Errorf("service cant provide both image and build")
		}
	}

	if in.Image != "" {
		out.Image = &in.Image
	}

	if in.Environment != nil {
		if len(in.Environment) > 0 {
			out.Envs = make([]body.Env, 0)
		}
		for k, v := range in.Environment {
			vv := util.PtrOr(v)
			if !HandleSpecialEnvs(k, vv, out) {
				out.Envs = append(out.Envs, body.Env{Name: k, Value: vv})
			}
		}
	}

	if len(in.Ports) > 0 {
		if out.Envs == nil {
			out.Envs = make([]body.Env, 0)
		}
		internalPorts := make([]string, 0)
		for _, env := range out.Envs {
			if env.Name == "INTERNAL_PORTS" {
				ports := strings.Split(env.Value, ",")
				internalPorts = append(internalPorts, ports...)
				break // only read the first occurance of this env
			}
		}
		for i, port := range in.Ports {
			if i == 0 {
				out.Envs = append(out.Envs, body.Env{Name: "PORT", Value: fmt.Sprintf("%d", port.Target)})
			} else {
				internalPorts = append(internalPorts, fmt.Sprintf("%d", port.Target))
			}

		}
		if len(internalPorts) > 0 {
			out.Envs = append(out.Envs, body.Env{Name: "INTERNAL_PORTS", Value: strings.Join(internalPorts, ",")})
		}
	}

	if len(in.DependsOn) > 0 {
		dependencies = make([]string, 0)
		for k := range in.DependsOn {
			dependencies = append(dependencies, k)
		}
	}

	if in.Deploy != nil {
		if in.Deploy.Resources.Limits != nil {
			// TODO: implement
		}
	}

	if in.Scale != nil {
		out.Replicas = in.Scale
	}

	if len(in.Volumes) > 0 {
		if out.Volumes == nil {
			out.Volumes = make([]body.Volume, 0)
		}
		for i, vol := range in.Volumes {
			if extra.Hash != "" {
				resolvedServerPath := filepath.Join(extra.Hash, strings.TrimPrefix(vol.Source, extra.CWD))

				out.Volumes = append(out.Volumes, body.Volume{Name: fmt.Sprintf("cli-%d", i), ServerPath: resolvedServerPath, AppPath: vol.Target})
			} else {
				out.Volumes = append(out.Volumes, body.Volume{Name: fmt.Sprintf("cli-%d", i), ServerPath: vol.Source, AppPath: vol.Target})
			}

		}
	}

	if in.Command != nil {
		if len(in.Command) > 0 {
			out.Args = make([]string, 0)
		}
		for _, cmd := range in.Command {
			out.Args = append(out.Args, cmd)
		}
	}

	if out.Visibility == "" {
		if len(in.Ports) == 0 {
			out.Visibility = "private"
		} else {
			out.Visibility = "public"
		}
	}

	if in.Gpus != nil {
		logrus.Warnln("service.Gpus is not implemented, yet ;)")
	}

	return out, dependencies, nil
}
