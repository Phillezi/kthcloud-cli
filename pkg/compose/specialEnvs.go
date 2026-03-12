package compose

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/kthcloud/cli/pkg/deploy"
)

// "Special" environment variables
const (
	// Configure how many cores that should be requested
	KTHCLOUD_CORES = "KTHCLOUD_CORES"
	// Configure how much RAM should be requested
	KTHCLOUD_RAM = "KTHCLOUD_RAM"
	// Configure how many replicas should be requested
	KTHCLOUD_REPLICAS = "KTHCLOUD_REPLICAS"
	// Configure gpus for deployment, syntax <claim-name>/<gpu-name>,<claim-name-1>/<gpu-name-1> (comma separated)
	KTHCLOUD_GPUS = "KTHCLOUD_GPUS"
	// Configure health path that should be polled
	KTHCLOUD_HEALTH_PATH = "KTHCLOUD_HEALTH_PATH"
	// Configure the visibility of the deployment
	KTHCLOUD_VISIBILITY = "KTHCLOUD_VISIBILITY"

	KTHCLOUD_ZONE = "KTHCLOUD_ZONE"

	KTHCLOUD_CUSTOMDOMAIN = "KTHCLOUD_CUSTOMDOMAIN"

	KTHCLOUD_ADMIN_NEVER_STALE = "KTHCLOUD_ADMIN_NEVER_STALE"
)

var (
	specialEnvs = []string{
		KTHCLOUD_CORES,
		KTHCLOUD_RAM,
		KTHCLOUD_REPLICAS,
		KTHCLOUD_GPUS,
		KTHCLOUD_HEALTH_PATH,
		KTHCLOUD_VISIBILITY,
		KTHCLOUD_ZONE,
		KTHCLOUD_CUSTOMDOMAIN,
		KTHCLOUD_ADMIN_NEVER_STALE,
	}

	specialEnvsHandlers = map[string]func(value string, deployment *deploy.BodyDeploymentCreate) error{
		KTHCLOUD_CORES: func(value string, deployment *deploy.BodyDeploymentCreate) error {
			cores, err := strconv.ParseFloat(value, 32)
			if err != nil {
				return err
			}
			deployment.CpuCores = new(float32(cores))
			return nil
		},
		KTHCLOUD_RAM: func(value string, deployment *deploy.BodyDeploymentCreate) error {
			ram, err := strconv.ParseFloat(value, 32)
			if err != nil {
				return err
			}
			deployment.Ram = new(float32(ram))
			return nil
		},
		KTHCLOUD_REPLICAS: func(value string, deployment *deploy.BodyDeploymentCreate) error {
			replicas, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			deployment.Replicas = &replicas
			return nil
		},
		KTHCLOUD_GPUS: func(value string, deployment *deploy.BodyDeploymentCreate) error {
			if value == "" {
				return nil
			}
			pairs := strings.Split(value, ",")
			if (*deployment).Gpus == nil {
				(*deployment).Gpus = new(make([]deploy.BodyDeploymentGPU, 0, len(pairs)))
			}
			for _, pair := range pairs {
				pair = strings.TrimSpace(pair)
				parts := strings.Split(pair, "/")
				partsLen := len(parts)
				if partsLen != 2 {
					return fmt.Errorf("gpu config env needs exactly two parts split by /, got: %v", parts)
				}
				*(*deployment).Gpus = append(*(*deployment).Gpus, deploy.BodyDeploymentGPU{ClaimName: parts[0], Name: parts[1]})
			}
			return nil
		},
		KTHCLOUD_HEALTH_PATH: func(value string, deployment *deploy.BodyDeploymentCreate) error {
			if value == "" {
				value = "/"
			}
			deployment.HealthCheckPath = &value
			return nil
		},
		KTHCLOUD_VISIBILITY: func(value string, deployment *deploy.BodyDeploymentCreate) error {
			visibility := strings.ToLower(value)
			switch visibility {
			case "private", "public", "auth":
				deployment.Visibility = new(deploy.BodyDeploymentCreateVisibility(visibility))
			default:
				return ErrInvalidDeploymentVisibility
			}
			return nil
		},
		KTHCLOUD_ZONE: func(value string, deployment *deploy.BodyDeploymentCreate) error {
			if value := strings.TrimSpace(value); value != "" {
				deployment.Zone = new(value)
			}
			return nil
		},
		KTHCLOUD_CUSTOMDOMAIN: func(value string, deployment *deploy.BodyDeploymentCreate) error {
			if len(value) > 243 {
				return ErrCustomDomainTooLong
			}
			if value := strings.TrimSpace(value); value != "" {
				deployment.CustomDomain = new(value)
			}
			return nil
		},
		KTHCLOUD_ADMIN_NEVER_STALE: func(value string, deployment *deploy.BodyDeploymentCreate) error {
			v, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}
			deployment.NeverStale = new(v)
			return nil
		},
	}
)

func GetSpecialEnvs() []string {
	return slices.Clone(specialEnvs)
}

func GetSpecialEnvsHandlers() map[string]func(value string, deployment *deploy.BodyDeploymentCreate) error {
	return maps.Clone(specialEnvsHandlers)
}
