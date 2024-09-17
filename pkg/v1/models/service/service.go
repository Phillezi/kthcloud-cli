package service

import (
	"fmt"
	"go-deploy/dto/v2/body"
	"strings"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

type Service struct {
	Image       string   `yaml:"image,omitempty"`
	Environment EnvVars  `yaml:"environment,omitempty"`
	Ports       []string `yaml:"ports,omitempty"`
	Volumes     []string `yaml:"volumes,omitempty"`
	Command     []string `yaml:"command,omitempty"`
}

func (s *Service) ToDeployment(name string, projectDir string) *body.DeploymentCreate {
	specialEnvs := []string{
		"KTHCLOUD_CORES",
		"KTHCLOUD_RAM",
		"KTHCLOUD_REPLICAS",
		"KTHCLOUD_HEALTH_PATH",
		"KTHCLOUD_VISIBILITY",
	}
	envsMap := make(map[string]bool)
	var envs []body.Env
	for envName, value := range s.Environment {
		if _, exists := envsMap[envName]; !exists {
			if !util.Contains(specialEnvs, envName) {
				envsMap[envName] = true
				envs = append(envs, body.Env{
					Name:  envName,
					Value: value,
				})
			}
		}
	}

	// Default to public
	visibility := "public"
	if len(s.Ports) == 0 {
		// if no ports are forwarded make the deployment private
		visibility = "private"
	} else {
		ports := strings.Split(s.Ports[0], ":")
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

	// Get the preferred zone and check if it is valid
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

	volumes := ToVolumes(s.Volumes, projectDir)

	// Get configuration from "special" set environment variables
	cores, ram, replicas, healthPath, visibilityConf := s.Environment.ResolveConfigEnvs()
	if cores == nil {
		cores = util.Float64Pointer(0.2)
	}
	if ram == nil {
		ram = util.Float64Pointer(0.5)
	}
	if replicas == nil {
		replicas = util.IntPointer(1)
	}
	if visibilityConf != "" {
		visibilityConf := strings.ToLower(visibilityConf)
		if util.Contains([]string{"private", "public", "auth"}, visibilityConf) {
			visibility = visibilityConf
		} else {
			log.Warnln("KTHCLOUD_VISIBILITY is set to: ", visibilityConf, " which is not valid, must be one of: private public auth.")
		}
	}

	return &body.DeploymentCreate{
		Name:            name,
		CpuCores:        cores,
		RAM:             ram,
		Replicas:        replicas,
		Envs:            envs,
		Image:           &s.Image,
		Visibility:      visibility,
		Args:            s.Command,
		Zone:            zone,
		Volumes:         volumes,
		HealthCheckPath: healthPath,
	}
}

func (s *Service) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("  Image: %s\n", s.Image))

	if len(s.Environment) > 0 {
		sb.WriteString("  Environment:\n")
		for key, value := range s.Environment {
			sb.WriteString(fmt.Sprintf("    %s: %s\n", key, value))
		}
	} else {
		sb.WriteString("  Environment: None\n")
	}

	if len(s.Ports) > 0 {
		sb.WriteString("  Ports:\n")
		for _, port := range s.Ports {
			sb.WriteString(fmt.Sprintf("    %s\n", port))
		}
	} else {
		sb.WriteString("  Ports: None\n")
	}

	if len(s.Volumes) > 0 {
		sb.WriteString("  Volumes:\n")
		for _, volume := range s.Volumes {
			sb.WriteString(fmt.Sprintf("    %s\n", volume))
		}
	} else {
		sb.WriteString("  Volumes: None\n")
	}

	if len(s.Command) > 0 {
		sb.WriteString("  Command:\n")
		for _, cmd := range s.Command {
			sb.WriteString(fmt.Sprintf("    %s\n", cmd))
		}
	} else {
		sb.WriteString("  Command: None\n")
	}

	return sb.String()
}
