package compose

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Phillezi/common/utils/or"
	composev2 "github.com/compose-spec/compose-go/v2/types"
	"github.com/kthcloud/cli/internal/defaults"
	"github.com/kthcloud/cli/pkg/deploy"
	"github.com/kthcloud/cli/pkg/kthcloud/project"
)

type ConverterImpl struct {
	Zone string
}

func (ci ConverterImpl) Convert(in *composev2.Project) (*project.Project, error) {
	out := &project.Project{
		Sevices: make(map[string]project.Service, len(in.Services)),
	}

	if err := ci.ConvertServices(in.Services, out.Sevices, in.Name, in.WorkingDir); err != nil {
		return nil, err
	}

	return out, nil
}

func (ci ConverterImpl) ConvertServices(in composev2.Services, out project.Services, projectName, cwd string) error {
	var errs error
	for k, v := range in {
		service, err := ci.ConvertService(k, v, projectName, cwd)
		if err != nil {
			errs = errors.Join(errs, err)
		}
		out[k] = service
	}

	return errs
}

func (ci ConverterImpl) ConvertService(name string, in composev2.ServiceConfig, projectName, cwd string) (out project.Service, err error) {
	out.Name = or.Or(in.ContainerName, name)
	if out.Name == "" {
		err = errors.Join(err, ErrNoNameOnService)
	}

	if in.Build != nil {
		if in.Image != "" {
			err = errors.Join(err, ErrBuildAndImageProvided)
		}
		if out.Build == nil {
			out.Build = &project.Build{}
		}
		out.Build.ContainerFile = in.Build.Dockerfile
		out.Build.Context = in.Build.Context
	}

	if in.Image != "" {
		out.Image = &in.Image
	}

	envs, errEnv := convertEnvs(in.Environment, &out.BodyDeploymentCreate)
	if errEnv != nil {
		err = errors.Join(err, errEnv)
	}
	envPorts, errPort := convertPorts(in.Ports)
	if errPort != nil {
		err = errors.Join(err, errPort)
	}

	if envPorts != nil && envs != nil {
		out.Envs = new(append(envs, envPorts...))
	} else if envs != nil {
		out.Envs = &envs
	} else {
		out.Envs = &envPorts
	}

	volumes, errVol := convertVolumes(in.Volumes, projectName, cwd)
	if errVol != nil {
		err = errors.Join(err, errVol)
	}
	out.Volumes = &volumes

	if in.Deploy != nil {
		if in.Deploy.Resources.Limits != nil {
			// TODO: implement
		}
	}

	if in.Scale != nil {
		out.Replicas = in.Scale
	}

	if len(in.DependsOn) > 0 {
		out.Dependencies = make([]string, 0, len(in.DependsOn))
		for k := range in.DependsOn {
			// todo support serviceStarted vs serviceHealhty?
			out.Dependencies = append(out.Dependencies, k)
		}
	}

	if in.Command != nil {
		out.Args = new(make([]string, 0, len(in.Command)))
		for _, cmd := range in.Command {
			*out.Args = append(*out.Args, cmd)
		}
	}

	if out.Visibility != nil && string(*out.Visibility) == "" {
		if len(in.Ports) == 0 {
			out.Visibility = new(deploy.BodyDeploymentCreateVisibilityPrivate)
		} else {
			out.Visibility = new(deploy.BodyDeploymentCreateVisibilityPublic)
		}
	}

	gpus, errGpu := convertGpus(in.Gpus)
	if errGpu != nil {
		err = errors.Join(err, errGpu)
	}
	if gpus != nil {
		*out.Gpus = append(*out.Gpus, gpus...)
	}

	applyDefaults(&out, ci.Zone)

	return
}

func applyDefaults(out *project.Service, zone string) {
	if out.Zone == nil && zone != "" {
		out.Zone = new(zone)
	}

	if out.HealthCheckPath == nil {
		out.HealthCheckPath = new(defaults.DefaultDeploymentHealthPath)
	}

	if out.CpuCores == nil {
		out.CpuCores = new(defaults.DefaultDeploymentSpecsCores)
	}

	if out.Ram == nil {
		out.Ram = new(defaults.DefaultDeploymentSpecsRam)
	}

	if out.Replicas == nil {
		out.Replicas = new(defaults.DefaultDeploymentSpecsReplicas)
	}
}

func convertEnvs(in composev2.MappingWithEquals, out *deploy.BodyDeploymentCreate) ([]deploy.BodyEnv, error) {
	if in != nil {
		envs := make([]deploy.BodyEnv, 0, len(in))
		for k, v := range in {
			if v == nil {
				continue
			}
			if h, isSpecial := specialEnvsHandlers[k]; isSpecial {
				h(*v, out)
			} else {
				envs = append(envs, deploy.BodyEnv{Name: k, Value: *v})
			}
		}
		return envs, nil
	}
	return nil, nil
}

func convertPorts(in []composev2.ServicePortConfig) ([]deploy.BodyEnv, error) {
	if in != nil {
		envs := make([]deploy.BodyEnv, 0, len(in))
		internalPorts := make([]string, 0)
		for _, env := range envs {
			if env.Name == "INTERNAL_PORTS" {
				ports := strings.Split(env.Value, ",")
				internalPorts = append(internalPorts, ports...)
				break // only read the first occurance of this env
			}
		}
		for i, port := range in {
			if i == 0 {
				envs = append(envs, deploy.BodyEnv{Name: "PORT", Value: fmt.Sprintf("%d", port.Target)})
			} else {
				internalPorts = append(internalPorts, fmt.Sprintf("%d", port.Target))
			}
		}
		if len(internalPorts) > 0 {
			envs = append(envs, deploy.BodyEnv{Name: "INTERNAL_PORTS", Value: strings.Join(internalPorts, ",")})
		}
		return envs, nil
	}

	return nil, nil
}

func convertVolumes(in []composev2.ServiceVolumeConfig, projectName string, cwd string) ([]deploy.BodyVolume, error) {
	if in != nil {
		volumes := make([]deploy.BodyVolume, 0, len(in))

		for i, vol := range in {
			if projectName != "" {
				resolvedServerPath := filepath.Join(projectName, strings.TrimPrefix(vol.Source, cwd))

				volumes = append(volumes, deploy.BodyVolume{Name: fmt.Sprintf("cli-%d", i), ServerPath: filepath.ToSlash(resolvedServerPath), AppPath: filepath.ToSlash(vol.Target)})
			} else {
				volumes = append(volumes, deploy.BodyVolume{Name: fmt.Sprintf("cli-%d", i), ServerPath: filepath.ToSlash(vol.Source), AppPath: filepath.ToSlash(vol.Target)})
			}
		}
		return volumes, nil
	}

	return nil, nil
}

func convertGpus(in []composev2.DeviceRequest) ([]deploy.BodyDeploymentGPU, error) {
	if in != nil {
		// NOT IMPLEMENTED

		/*gpus := make([]body.DeploymentGPU, 0, len(in))
		for _, gpu := range in {
			gpus = append(gpus, body.DeploymentGPU{
				Name: "hard-to-convert",
				ClaimName: "idk-how-to-convert",
			})
		}*/
	}

	return nil, nil
}

func convertDNS(in composev2.ServiceConfig) {
	if in.DNS != nil {
		// NOT IMPLEMENTED
	}
}
