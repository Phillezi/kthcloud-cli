package compose

import (
	"go-deploy/dto/v2/body"
	"kthcloud-cli/pkg/util"
	"strings"
)

func serviceToDepl(service Service, name string) *body.DeploymentCreate {
	envsMap := make(map[string]bool)
	var envs []body.Env
	for envName, value := range service.Environment {
		if _, exists := envsMap[envName]; !exists {
			envsMap[envName] = true
			envs = append(envs, body.Env{
				Name:  envName,
				Value: value,
			})
		}
	}

	// Default to public
	visibility := "public"
	if len(service.Ports) == 0 {
		// if no ports are forwarded make the deployment private
		visibility = "private"
	} else {
		ports := strings.Split(service.Ports[0], ":")
		if len(ports) != 0 {
			port := ports[len(ports)-1]
			// Add the PORT environment variable and set it to the first exposed port
			if _, exists := envsMap["PORT"]; !exists {
				envsMap["PORT"] = true
				envs = append(envs, body.Env{
					Name:  "PORT",
					Value: port,
				})
			}
		}
	}

	return &body.DeploymentCreate{
		Name:       name,
		CpuCores:   util.Float64Pointer(0.2),
		RAM:        util.Float64Pointer(0.5),
		Replicas:   util.IntPointer(1),
		Envs:       envs,
		Image:      &service.Image,
		Visibility: visibility,
		Args:       service.Command,
	}
}
