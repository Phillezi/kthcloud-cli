package defaults

import "time"

const (
	DefaultRequestTimeout time.Duration = 10 * time.Second

	DefaultDeployBaseURL string = "https://api.cloud.cbh.kth.se/deploy"

	DefaultRedirectSchemeHostPort string = "http://localhost:3000"
	DefaultRedirectBasePath       string = "/auth/callback"

	DefaultKeycloakBaseURL      string = "https://iam.cloud.cbh.kth.se"
	DefaultKeycloakRealm        string = "cloud"
	DefaultKeycloakClientID     string = "landing"
	DefaultKeycloakClientSecret string = ""

	DefaultStorageManagerProxy string = "https://sm-proxy.app.cloud.cbh.kth.se"

	DefaultCreateWorkflow bool = true

	DefaultSaveSecrets bool = false

	DefaultComposeUpBuildAll   bool = false
	DefaultComposeUpDetached   bool = false
	DefaultComposeUpTryVolumes bool = false

	DefaultNonInteractive bool = false

	DefaultPsAll bool = false

	DefaultParseOutputOnlyJSON bool = false

	DefaultRemoveCustomDeployments bool = false

	DefaultDeploymentCores    float64 = 0.2
	DefaultDeploymentMemory   float64 = 0.5
	DefaultDeploymentReplicas int     = 1

	DefaultDeploymentVisibility string = "public"

	DefaultDeploymentRunInteractive bool = false
	DefaultDeploymentRunTTY         bool = false
	DefaultDeploymentRunRemove      bool = false
	DefaultDeploymentRunDetatch     bool = false

	DefaultDeploymentZone string = "se-flem-2"
)
