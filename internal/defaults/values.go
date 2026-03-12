package defaults

import "time"

const (
	DefaultRequestTimeout = 30 * time.Second

	DefaultDeployAPIBaseURL  = "https://api.cloud.cbh.kth.se/deploy"
	DefaultSMProxyAPIBaseURL = "https://sm-proxy.app.cloud.cbh.kth.se"

	DefaultKeycloakBaseURL      = "https://iam.cloud.cbh.kth.se"
	DefaultKeycloakRealm        = "cloud"
	DefaultKeycloakClientID     = "landing"
	DefaultKeycloakClientSecret = ""
	DefaultLoginServerAddress   = "localhost:3000"

	DefaultDeploymentVisibility    = "public"
	DefaultDeploymentHealthPath    = "/healthz"
	DefaultDeploymentSpecsCores    = float32(0.2)
	DefaultDeploymentSpecsRam      = float32(0.5)
	DefaultDeploymentSpecsReplicas = 1

	DefaultVMSpecsCores = float64(4)
	DefaultVMSpecsRam   = float64(8)
	DefaultVMDiskSize   = float64(20)

	DefaultZone = "se-flem-2"

	DefaultKeystoreSessionKey  = "default"
	DefaultKeystoreServiceName = "kthcloud-cli"
	DefaultKeystoreFallbackDir = "/tmp/kthcloud-cli" // TODO: this is temporary

	DefaultSSHBinary = "ssh"
	DefaultSSHHost   = "deploy.cloud.cbh.kth.se"
)

var (
	DefaultComposeFileNames []string = []string{
		"kthcloud.docker-compose.yaml",
		"kthcloud.docker-compose.yml",
		"kthcloud.compose.yaml",
		"kthcloud.compose.yml",
		"docker-compose.yaml",
		"docker-compose.yml",
		"compose.yaml",
		"compose.yml",
	}

	DefaultComposeDotEnvFileNames []string = []string{
		".env",
		".env.kthcloud",
	}
)
