package compose

import (
	"go-deploy/dto/v2/body"
	"kthcloud-cli/internal/model"
	"kthcloud-cli/pkg/util"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

func serviceToDepl(service model.Service, name string, projectDir string) *body.DeploymentCreate {
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

	// Get the prefered zone and check if it is valid
	// Note: zoneName is the zoneName name
	zoneName := viper.GetString("zone")
	// TODO: Might be better to fetch /v2/zones and check if the zone is valid that way
	// the only downside to that is that we have to rewrite the converter to either take in
	// a Session or a list of valid zones, since it would be highly inefficient to fetch
	// /v2/zones for every conversion
	if zoneName != "" && zoneName != "se-flem2" && zoneName != "se-kista" {
		log.Warnln("Specified zone is (most likely) not valid")
	}

	var zone *string
	if zoneName == "" {
		zone = nil
	} else {
		zone = &zoneName
	}

	volumes := ToVolumes(service.Volumes, projectDir)

	return &body.DeploymentCreate{
		Name:       name,
		CpuCores:   util.Float64Pointer(0.2),
		RAM:        util.Float64Pointer(0.5),
		Replicas:   util.IntPointer(1),
		Envs:       envs,
		Image:      &service.Image,
		Visibility: visibility,
		Args:       service.Command,
		Zone:       zone,
		Volumes:    volumes,
	}
}
