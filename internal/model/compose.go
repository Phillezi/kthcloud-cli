package model

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

type Service struct {
	Image       string            `yaml:"image,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Ports       []string          `yaml:"ports,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
	Command     []string          `yaml:"command,omitempty"`
}

type ComposeFile struct {
	Version  string                            `yaml:"version"`
	Services map[string]map[string]interface{} `yaml:"services"`
}

// Create unique hash by combining services names
func Hash(services map[string]Service) string {
	keys := make([]string, 0, len(services))

	for key := range services {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	combinedKeys := strings.Join(keys, "")

	hash := sha256.New()
	hash.Write([]byte(combinedKeys))
	hashBytes := hash.Sum(nil)

	return hex.EncodeToString(hashBytes)
}
