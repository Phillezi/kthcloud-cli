package defaults

import "time"

const (
	DefaultRequestTimeout = 30 * time.Second

	DefaultDeployAPIBaseURL  = "https://api.cloud.cbh.kth.se/deploy/v2"
	DefaultSMProxyAPIBaseURL = "https://sm-proxy.app.cloud.cbh.kth.se"

	DefaultKeycloakBaseURL      = "https://iam.cloud.cbh.kth.se"
	DefaultKeycloakRealm        = "cloud"
	DefaultKeycloakClientID     = "landing"
	DefaultKeycloakClientSecret = ""

	DefaultDeploymentVisibility    = "public"
	DefaultDeploymentHealthPath    = "/healthz"
	DefaultDeploymentSpecsCores    = float64(0.2)
	DefaultDeploymentSpecsRam      = float64(0.5)
	DefaultDeploymentSpecsReplicas = 1

	DefaultVMSpecsCores = float64(4)
	DefaultVMSpecsRam   = float64(8)
	DefaultVMDiskSize   = float64(20)

	DefaultZone = "se-flem-2"
)
