package compose

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go-deploy/dto/v2/body"
	"os"
	"sort"
	"strings"

	"github.com/Phillezi/kthcloud-cli/pkg/v1/models/service"
	"gopkg.in/yaml.v3"
)

type Compose struct {
	Services map[string]*service.Service `yaml:"services"`
}

func New(filePath string) (*Compose, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open compose file: %v", err)
	}
	defer file.Close()

	c := &Compose{}

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(c); err != nil {
		return nil, fmt.Errorf("could not decode YAML file: %v", err)
	}

	return c, nil
}

func (c *Compose) String() string {
	var sb strings.Builder

	sb.WriteString("Compose:\n")
	sb.WriteString(fmt.Sprintf("Hash: %s\n", c.Hash()))
	sb.WriteString("Services:\n")

	indent := "  "

	for name, service := range c.Services {
		sb.WriteString(fmt.Sprintf("  %s:\n", name))
		sb.WriteString(indent)
		sb.WriteString(strings.ReplaceAll(service.String(), "\n", "\n"+indent))
		sb.WriteString("\n")
	}

	return sb.String()
}

func (c *Compose) ToDeploymentsWDeps() ([]*body.DeploymentCreate, map[string][]string) {
	projectDirectory := c.Hash()

	var depls []*body.DeploymentCreate
	deps := make(map[string][]string)

	for name, service := range c.Services {
		depls = append(depls, service.ToDeployment(name, projectDirectory))
		deps[name] = service.Dependencies
		if deps[name] == nil {
			deps[name] = []string{}
		}
	}

	return depls, deps
}

func (c *Compose) ToDeployments() []*body.DeploymentCreate {
	projectDirectory := c.Hash()

	var depls []*body.DeploymentCreate

	for name, service := range c.Services {
		depls = append(depls, service.ToDeployment(name, projectDirectory))
	}

	return depls
}

func (c *Compose) Hash() string {
	keys := make([]string, 0, len(c.Services))

	for key := range c.Services {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	combinedKeys := strings.Join(keys, "")

	hash := sha256.New()
	hash.Write([]byte(combinedKeys))
	hashBytes := hash.Sum(nil)

	return hex.EncodeToString(hashBytes)
}
