package compose

import (
	"fmt"
	"kthcloud-cli/internal/model"

	log "github.com/sirupsen/logrus"

	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// ParseComposeFile reads and parses a docker-compose file
func ParseComposeFile(filename string) (map[string]model.Service, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var composeFile model.ComposeFile
	if err := yaml.Unmarshal(data, &composeFile); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	services := make(map[string]model.Service)
	for name, serviceData := range composeFile.Services {
		service := model.Service{}
		if err := processService(serviceData, &service); err != nil {
			return nil, err
		}
		services[name] = service
	}

	return services, nil
}

func processService(serviceData map[string]interface{}, service *model.Service) error {
	if env, ok := serviceData["environment"]; ok {
		service.Environment = parseEnvironment(env)
	}
	if ports, ok := serviceData["ports"]; ok {
		service.Ports = parseStringList(ports)
	}
	if volumes, ok := serviceData["volumes"]; ok {
		service.Volumes = parseStringList(volumes)
	}
	if command, ok := serviceData["command"]; ok {
		service.Command = parseStringList(command)
	}
	if image, ok := serviceData["image"].(string); ok {
		service.Image = image
	}

	// Add more fields later

	return nil
}

// parseEnvironment handles both the list and map formats for environment variables
// TODO: Resolve .env files and exported envs?
func parseEnvironment(env interface{}) map[string]string {
	result := make(map[string]string)

	switch v := env.(type) {
	case []interface{}:
		for _, item := range v {
			strItem := item.(string)
			parts := strings.SplitN(strItem, "=", 2)
			if len(parts) == 2 {
				result[parts[0]] = parts[1]
			} else {
				result[parts[0]] = ""
			}
		}
	case map[interface{}]interface{}:
		for key, value := range v {
			strKey := key.(string)
			strValue := value.(string)
			result[strKey] = strValue
		}
	case map[string]interface{}:
		for key, value := range v {
			var strValue string
			switch v := value.(type) {
			case string:
				strValue = v
			case float64:
				strValue = fmt.Sprintf("%f", v)
			case bool:
				strValue = fmt.Sprintf("%t", v)
			default:
				strValue = fmt.Sprintf("%v", v)
			}
			result[key] = strValue
		}
	default:
		// Log unexpected types for debugging
		log.Printf("Unexpected environment format: %T\n", v)
	}

	return result
}

// parseStringList converts an interface{} to a []string if possible
func parseStringList(input interface{}) []string {
	var result []string

	switch v := input.(type) {
	case []interface{}:
		for _, item := range v {
			result = append(result, item.(string))
		}
	}

	return result
}
