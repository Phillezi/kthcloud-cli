package compose

import (
	"strconv"
	"strings"

	"github.com/kthcloud/cli/pkg/utils"
	"github.com/kthcloud/go-deploy/dto/v2/body"
)

// "Special" environment variables
const (
	// Configure how many cores that should be requested
	KTHCLOUD_CORES = "KTHCLOUD_CORES"
	// Configure how much RAM should be requested
	KTHCLOUD_RAM = "KTHCLOUD_RAM"
	// Configure how many replicas should be requested
	KTHCLOUD_REPLICAS = "KTHCLOUD_REPLICAS"
	// Configure health path that should be polled
	KTHCLOUD_HEALTH_PATH = "KTHCLOUD_HEALTH_PATH"
	// Configure the visibility of the deployment
	KTHCLOUD_VISIBILITY = "KTHCLOUD_VISIBILITY"

	KTHCLOUD_ZONE = "KTHCLOUD_ZONE"

	KTHCLOUD_CUSTOMDOMAIN = "KTHCLOUD_CUSTOMDOMAIN"

	KTHCLOUD_ADMIN_NEVER_STALE = "KTHCLOUD_ADMIN_NEVER_STALE"
)

var (
	_ = []string{
		KTHCLOUD_CORES,
		KTHCLOUD_RAM,
		KTHCLOUD_REPLICAS,
		KTHCLOUD_HEALTH_PATH,
		KTHCLOUD_VISIBILITY,
		KTHCLOUD_ZONE,
		KTHCLOUD_CUSTOMDOMAIN,
		KTHCLOUD_ADMIN_NEVER_STALE,
	}

	_ = map[string]func(value string, deployment *body.DeploymentCreate) error{
		KTHCLOUD_CORES: func(value string, deployment *body.DeploymentCreate) error {
			cores, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			deployment.CpuCores = &cores
			return nil
		},
		KTHCLOUD_RAM: func(value string, deployment *body.DeploymentCreate) error {
			ram, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			deployment.RAM = &ram
			return nil
		},
		KTHCLOUD_REPLICAS: func(value string, deployment *body.DeploymentCreate) error {
			replicas, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			deployment.Replicas = &replicas
			return nil
		},
		KTHCLOUD_HEALTH_PATH: func(value string, deployment *body.DeploymentCreate) error {
			if value == "" {
				value = "/"
			}
			deployment.HealthCheckPath = &value
			return nil
		},
		KTHCLOUD_VISIBILITY: func(value string, deployment *body.DeploymentCreate) error {
			visibility := strings.ToLower(value)
			switch visibility {
			case "private", "public", "auth":
				deployment.Visibility = visibility
			default:
				return ErrInvalidDeploymentVisibility
			}
			return nil
		},
		KTHCLOUD_ZONE: func(value string, deployment *body.DeploymentCreate) error {
			if value := strings.TrimSpace(value); value != "" {
				deployment.Zone = utils.PtrOf(value)
			}
			return nil
		},
		KTHCLOUD_CUSTOMDOMAIN: func(value string, deployment *body.DeploymentCreate) error {
			if len(value) > 243 {
				return ErrCustomDomainTooLong
			}
			if value := strings.TrimSpace(value); value != "" {
				deployment.CustomDomain = utils.PtrOf(value)
			}
			return nil
		},
		KTHCLOUD_ADMIN_NEVER_STALE: func(value string, deployment *body.DeploymentCreate) error {
			v, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}
			deployment.NeverStale = v
			return nil
		},
	}
)
