package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// EnvVars is a custom type to handle key=value pairs in YAML
type EnvVars map[string]string

func (e *EnvVars) UnmarshalYAML(node *yaml.Node) error {
	// Handle if node is a mapping or a sequence
	switch node.Kind {
	case yaml.MappingNode:
		// Handle case where the YAML data is in mapping format
		*e = make(map[string]string)
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]
			key := keyNode.Value
			value := valueNode.Value
			(*e)[key] = value
		}
	case yaml.SequenceNode:
		// Handle case where the YAML data is in sequence format
		*e = make(map[string]string)
		for _, item := range node.Content {
			if item.Kind != yaml.ScalarNode {
				return fmt.Errorf("expected scalar node, got %v", item.Kind)
			}
			parts := strings.SplitN(item.Value, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid key=value format: %s", item.Value)
			}
			(*e)[parts[0]] = parts[1]
		}
	default:
		return fmt.Errorf("expected mapping or sequence node, got %v", node.Kind)
	}

	return nil
}

func (envs EnvVars) ResolveConfigEnvs() (cpucores, ram *float64, replicas *int, healthPath *string) {
	return envs.getFloat64ptr("KTHCLOUD_CORES"), envs.getFloat64ptr("KTHCLOUD_RAM"), envs.getIntptr("KTHCLOUD_REPLICAS"), envs.getStrptr("KTHCLOUD_HEALTH_PATH")
}

func (envs EnvVars) getFloat64ptr(key string) *float64 {
	if envValue, exists := envs[key]; exists {
		parsedValue, err := strconv.ParseFloat(envValue, 64)
		if err != nil {
			logrus.Errorf("Error parsing %s: %v\n", key, err)
		} else {
			return &parsedValue
		}
	}
	return nil
}

func (envs EnvVars) getIntptr(key string) *int {
	if envValue, exists := envs[key]; exists {
		parsedValue, err := strconv.ParseInt(envValue, 10, 64)
		if err != nil {
			logrus.Errorf("Error parsing %s: %v\n", key, err)
		} else {
			parsedInt := int(parsedValue)
			return &parsedInt
		}
	}
	return nil
}

func (envs EnvVars) getStrptr(key string) *string {
	if envValue, exists := envs[key]; exists {
		return &envValue
	}
	return nil
}
