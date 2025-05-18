package defaults

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
