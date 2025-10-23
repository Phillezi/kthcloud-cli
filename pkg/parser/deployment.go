// pkg/parser/parser.go
package parser

import (
	"fmt"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/kthcloud/cli/pkg/deploy"
	"github.com/kthcloud/cli/pkg/parser/validation"
	"github.com/kthcloud/cli/pkg/utils"
	"github.com/kthcloud/go-deploy/dto/v2/body"
)

// DeploymentFlags holds user-provided flags
type DeploymentFlags struct {
	Name       string
	Env        []string
	Volume     []string
	Port       []string
	CPU        float32
	RAM        float32
	Replicas   int
	Zone       string
	NeverStale bool
	Args       []string
	Domain     string
	Health     string
	Visibility string
}

// ParseDeployment takes positional args and flags, and constructs a deploy.BodyDeploymentCreate
func ParseDeployment(image string, flags DeploymentFlags) (*deploy.BodyDeploymentCreate, error) {
	if strings.TrimSpace(flags.Name) == "" {
		return nil, ErrNameReqired
	}
	body := &deploy.BodyDeploymentCreate{
		Name:  flags.Name,
		Image: &image,
	}

	if len(flags.Args) > 0 {
		body.Args = &flags.Args
	}
	if flags.CPU > 0 {
		body.CpuCores = &flags.CPU
	}
	if flags.RAM > 0 {
		body.Ram = &flags.RAM
	}
	if flags.Replicas != 1 {
		body.Replicas = &flags.Replicas
	}
	if flags.Zone != "" {
		body.Zone = &flags.Zone
	}
	if flags.Domain != "" {
		body.CustomDomain = &flags.Domain
	}
	if flags.Health != "" {
		body.HealthCheckPath = &flags.Health
	}
	if flags.Visibility != "" {
		v := deploy.BodyDeploymentCreateVisibility(flags.Visibility)
		body.Visibility = &v
	}
	if flags.NeverStale {
		body.NeverStale = &flags.NeverStale
	}

	if len(flags.Env) > 0 {
		envs := parseKeyValuePairs(flags.Env)
		body.Envs = &envs
	}

	if len(flags.Volume) > 0 {
		vols := parseVolumePairs(flags.Volume)
		body.Volumes = &vols
	}

	if len(flags.Port) > 0 {
		if body.Envs == nil {
			body.Envs = utils.PtrOf(make([]deploy.BodyEnv, 0, len(flags.Port)+1))
		}

		(*body.Envs) = append((*body.Envs), parsePorts(flags.Port)...)
	}

	backendModel, err := toBackendDeployment(body)
	if err != nil {
		return nil, err
	}
	if err := validation.Validate(backendModel); err != nil {
		return nil, err
	}

	return body, nil
}

// parseKeyValuePairs parses KEY=VALUE strings into []BodyEnv
func parseKeyValuePairs(pairs []string) []deploy.BodyEnv {
	var result []deploy.BodyEnv
	for _, p := range pairs {
		if kv := strings.SplitN(p, "=", 2); len(kv) == 2 {
			result = append(result, deploy.BodyEnv{
				Name:  kv[0],
				Value: kv[1],
			})
		}
	}
	return result
}

// parseVolumePairs parses local:remote pairs into []BodyVolume
func parseVolumePairs(pairs []string) []deploy.BodyVolume {
	var result []deploy.BodyVolume
	for i, v := range pairs {
		if parts := strings.SplitN(v, ":", 2); len(parts) == 2 {
			result = append(result, deploy.BodyVolume{
				Name:       fmt.Sprintf("cli-%d", i),
				AppPath:    parts[0],
				ServerPath: parts[1],
			})
		}
	}
	return result
}

func parsePorts(ports []string) []deploy.BodyEnv {
	envs := make([]deploy.BodyEnv, 0, len(ports)+1)

	var internalPorts []string

	for i, p := range ports {
		parts := strings.Split(p, ":")
		var containerPort string

		switch len(parts) {
		case 1:
			containerPort = parts[0]
		case 2:
			containerPort = parts[1]
		default:
			continue // skip invalid ones like "a:b:c"
		}

		internalPorts = append(internalPorts, containerPort)

		// First port gets the PORT env var
		// this will be used for the ingress
		if i == 0 {
			envs = append(envs, deploy.BodyEnv{
				Name:  "PORT",
				Value: containerPort,
			})
		}
	}

	// Add INTERNAL_PORTS env var with comma+space separation
	if len(internalPorts) > 1 {
		envs = append(envs, deploy.BodyEnv{
			Name:  "INTERNAL_PORTS",
			Value: strings.Join(internalPorts, ", "),
		})
	}

	return envs
}

func toBackendDeployment(b *deploy.BodyDeploymentCreate) (*body.DeploymentCreate, error) {
	backendModel := &body.DeploymentCreate{}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json", // respect `json` tags for field matching
		Result:  backendModel,
	})
	if err != nil {
		return nil, err
	}

	if err := decoder.Decode(b); err != nil {
		return nil, err
	}

	return backendModel, nil
}
