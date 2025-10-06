package convert

import (
	"strconv"
	"strings"

	"github.com/kthcloud/go-deploy/dto/v2/body"
)

const (
	KTHCLOUD_CORES       = "KTHCLOUD_CORES"
	KTHCLOUD_RAM         = "KTHCLOUD_RAM"
	KTHCLOUD_REPLICAS    = "KTHCLOUD_REPLICAS"
	KTHCLOUD_HEALTH_PATH = "KTHCLOUD_HEALTH_PATH"
	KTHCLOUD_VISIBILITY  = "KTHCLOUD_VISIBILITY"
)

var (
	SpecialEnvs = []string{
		KTHCLOUD_CORES,
		KTHCLOUD_RAM,
		KTHCLOUD_REPLICAS,
		KTHCLOUD_HEALTH_PATH,
		KTHCLOUD_VISIBILITY,
	}
)

func HandleSpecialEnvs(key, value string, deployment *body.DeploymentCreate) bool {
	switch key {
	case KTHCLOUD_CORES:
		cores, err := strconv.ParseFloat(value, 64)
		if err != nil {
			break
		}
		deployment.CpuCores = &cores
	case KTHCLOUD_RAM:
		ram, err := strconv.ParseFloat(value, 64)
		if err != nil {
			break
		}
		deployment.RAM = &ram
	case KTHCLOUD_REPLICAS:
		replicas, err := strconv.Atoi(value)
		if err != nil {
			break
		}
		deployment.Replicas = &replicas
	case KTHCLOUD_HEALTH_PATH:
		if value == "" {
			value = "/"
		}
		deployment.HealthCheckPath = &value
	case KTHCLOUD_VISIBILITY:
		visibility := strings.ToLower(value)
		switch visibility {
		case "private", "public", "auth":
			deployment.Visibility = visibility
		}
	default:
		return false
	}
	return true
}
