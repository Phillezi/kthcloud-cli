package service

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Build struct {
	// context defines either a path to a directory containing a Dockerfile, or a URL to a git repository.
	// When the value supplied is a relative path, it is interpreted as relative to the project directory.
	// If not set explicitly, context defaults to project directory (.).
	Context string `yaml:"context,omitempty"`
	// The name of the dockerfile to use inside the context directory (can be relative path to the context directory).
	Dockerfile string `yaml:"dockerfile,omitempty"`
	// Build args to pass to the dockerfile.
	Args []string `yaml:"args,omitempty"`
}

func (b *Build) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		// If the node is a scalar (string), assign it to Raw
		return value.Decode(&b.Context)
	} else if value.Kind == yaml.MappingNode {
		// If the node is a mapping, decode into the struct fields
		type buildAlias Build // Create an alias to avoid recursion
		return value.Decode((*buildAlias)(b))
	}
	return fmt.Errorf("invalid type for Build: %v", value.Kind)
}
func (b Build) String() string {
	if b.Context != "" && b.Args == nil && b.Dockerfile == "" {
		return b.Context
	}

	var sb strings.Builder
	sb.WriteString("\n")
	if b.Context != "" {
		sb.WriteString(fmt.Sprintf("  Context: %s\n", b.Context))
	}
	if b.Dockerfile != "" {
		sb.WriteString(fmt.Sprintf("  Dockerfile: %s\n", b.Dockerfile))
	}
	if len(b.Args) > 0 {
		sb.WriteString(fmt.Sprintf("  Args: %s\n", strings.Join(b.Args, ", ")))
	}

	return sb.String()
}
