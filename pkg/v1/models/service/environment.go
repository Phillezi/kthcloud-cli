package service

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// EnvVars is a custom type to handle key=value pairs in YAML
type EnvVars map[string]string

// UnmarshalYAML implements the yaml.Unmarshaler interface for EnvVars
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
