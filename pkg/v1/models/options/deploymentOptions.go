package options

import (
	"strings"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/kthcloud/go-deploy/dto/v2/body"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type DeploymentOptions struct {
	Image           string
	Entrypoint      *string
	ContainerName   string
	InteractiveLogs bool
	RemoveOnExit    bool
	Visibility      string
	Cores           float64
	Ram             float64
	Replicas        int
	HealthCheck     string
	EnvVars         []string
	Ports           []string
	Volumes         []string
}

func (do *DeploymentOptions) ToDeploymentCreate() (*body.DeploymentCreate, error) {
	specialEnvs := []string{
		"KTHCLOUD_CORES",
		"KTHCLOUD_RAM",
		"KTHCLOUD_REPLICAS",
		"KTHCLOUD_HEALTH_PATH",
		"KTHCLOUD_VISIBILITY",
	}
	envsMap := make(map[string]bool)
	var envs []body.Env
	for _, value := range do.EnvVars {
		env := strings.SplitN(value, "=", 2)
		if len(env) < 1 {
			continue
		} else if len(env) == 1 {
			env = append(env, "")
		}
		if _, exists := envsMap[env[0]]; !exists {
			if !util.Contains(specialEnvs, env[0]) {
				envsMap[env[0]] = true
				envs = append(envs, body.Env{
					Name:  env[0],
					Value: env[1],
				})
			}
		}
	}

	visibility := do.Visibility
	if visibility == "" {
		visibility = "public"
	}
	if len(do.Ports) == 0 {
		// if no ports are forwarded make the deployment private
		visibility = "private"
	} else {
		ports := strings.Split(do.Ports[0], ":")
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
	zoneName := viper.GetString("zone")
	// TODO: Might be better to fetch /v2/zones and check if the zone is valid that way
	// the only downside to that is that we have to rewrite the converter to either take in
	// a Session or a list of valid zones, since it would be highly inefficient to fetch
	// /v2/zones for every conversion
	if zoneName != "" && zoneName != "se-flem2" && zoneName != "se-kista" {
		logrus.Warnln("Specified zone is (most likely) not valid")
	}

	var zone *string
	if zoneName == "" {
		zone = nil
	} else {
		zone = &zoneName
	}
	var commands []string
	if do.Entrypoint != nil {
		commands = append(commands, strings.Split(*do.Entrypoint, " ")...)
	}
	var volumes []body.Volume
	for _, value := range do.Volumes {
		vol := strings.SplitN(value, "=", 2)
		if len(vol) < 2 {
			continue
		}
		volumes = append(volumes, body.Volume{
			Name:       "cli-generated",
			ServerPath: vol[0],
			AppPath:    vol[1],
		})
	}
	return &body.DeploymentCreate{
		Name:            do.ContainerName,
		CpuCores:        util.OrNil(do.Cores),
		RAM:             util.OrNil(do.Ram),
		Replicas:        util.OrNil(do.Replicas),
		Envs:            envs,
		Image:           util.OrNil(do.Image),
		Visibility:      visibility,
		Args:            commands,
		Zone:            zone,
		Volumes:         volumes,
		HealthCheckPath: util.OrNil(do.HealthCheck),
	}, nil
}
