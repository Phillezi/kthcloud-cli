package parser

import (
	"testing"

	"github.com/kthcloud/cli/pkg/deploy"
	"github.com/stretchr/testify/assert"
)

func TestParseKeyValuePairs(t *testing.T) {
	input := []string{"FOO=bar", "BAZ=qux", "INVALID"}
	expected := []deploy.BodyEnv{
		{Name: "FOO", Value: "bar"},
		{Name: "BAZ", Value: "qux"},
	}

	result := parseKeyValuePairs(input)
	assert.Equal(t, expected, result)
}

func TestParseVolumePairs(t *testing.T) {
	input := []string{"/local:/remote", "/tmp:/app"}
	expected := []deploy.BodyVolume{
		{Name: "cli-0", AppPath: "/local", ServerPath: "/remote"},
		{Name: "cli-1", AppPath: "/tmp", ServerPath: "/app"},
	}

	result := parseVolumePairs(input)
	assert.Equal(t, expected, result)
}

func TestParsePorts(t *testing.T) {
	input := []string{"8080:80", "443"}
	result := parsePorts(input)

	// Check that PORT is first port
	assert.Equal(t, "80", result[0].Value)
	// INTERNAL_PORTS includes both
	assert.Equal(t, "80, 443", result[1].Value)
}

func TestParseDeployment_Minimal(t *testing.T) {
	flags := DeploymentFlags{
		Name: "myapp",
	}
	image := "nginx:latest"

	body, err := ParseDeployment(image, flags)
	assert.NoError(t, err)
	assert.Equal(t, "myapp", body.Name)
	assert.Equal(t, "nginx:latest", *body.Image)
	assert.Nil(t, body.Envs)
}

func TestParseDeployment_WithEnvVolumePort(t *testing.T) {
	flags := DeploymentFlags{
		Name:   "testapp",
		Env:    []string{"FOO=bar", "BAZ=qux"},
		Volume: []string{"/local:/remote"},
		Port:   []string{"8080:80", "443"},
	}
	image := "myimage:v1"

	body, err := ParseDeployment(image, flags)
	assert.NoError(t, err)
	assert.NotNil(t, body.Envs)
	assert.NotNil(t, body.Volumes)
	assert.Len(t, *body.Envs, 4)    // FOO + BAZ + PORT + INTERNAL_PORTS
	assert.Len(t, *body.Volumes, 1) // single volume

	// Env vars
	assert.Equal(t, "FOO", (*body.Envs)[0].Name)
	assert.Equal(t, "bar", (*body.Envs)[0].Value)
	assert.Equal(t, "BAZ", (*body.Envs)[1].Name)
	assert.Equal(t, "qux", (*body.Envs)[1].Value)

	// PORT env (first port)
	assert.Equal(t, "PORT", (*body.Envs)[2].Name)
	assert.Equal(t, "80", (*body.Envs)[2].Value)

	// INTERNAL_PORTS env
	assert.Equal(t, "INTERNAL_PORTS", (*body.Envs)[3].Name)
	assert.Equal(t, "80, 443", (*body.Envs)[3].Value)
}

func TestParseDeployment_ValidationFailure(t *testing.T) {
	flags := DeploymentFlags{
		Name: "", // invalid
	}
	image := "nginx:latest"

	body, err := ParseDeployment(image, flags)
	assert.Nil(t, body)
	assert.Error(t, err)
}

func TestToBackendDeployment(t *testing.T) {
	body := &deploy.BodyDeploymentCreate{
		Name: "myapp",
	}
	backend, err := toBackendDeployment(body)
	assert.NoError(t, err)
	assert.Equal(t, "myapp", backend.Name)
}
