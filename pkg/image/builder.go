package builder

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kthcloud/cli/pkg/ui/io"
)

var ErrNoBuilderInstalled = errors.New("no builder found, please install docker")

type Builder interface {
	Login(ctx context.Context, registry, user, password string, tag string, name string) error
	Build(ctx context.Context, containerFilePath string, context string, tag string, push bool, name string) error
}

func Get() (Builder, error) {
	if _, err := exec.LookPath("docker"); err == nil {
		return DockerBuildxImpl{}, nil
	}

	// TODO: Also make a podman impl
	return nil, ErrNoBuilderInstalled
}

type DockerBuildxImpl struct{}

func (d DockerBuildxImpl) Login(ctx context.Context, registry, user, password string, tag string, name string) error {
	dockerContext := d.getDockerContextForTag(tag)

	cmd := exec.CommandContext(ctx,
		"docker",
		"--context", dockerContext,
		"login", registry,
		"-u", user,
		"--password-stdin",
	)

	cmd.Stdin = bytes.NewBufferString(password)
	cmd.Stdout = &io.PrefixWriter{
		Prefix: fmt.Sprintf("[%s]\t", name),
		Writer: os.Stdout,
	}
	cmd.Stderr = &io.PrefixWriter{
		Prefix: fmt.Sprintf("[%s]\t", name),
		Writer: os.Stderr,
	}

	return cmd.Run()
}

func (d DockerBuildxImpl) Build(ctx context.Context, containerFilePath string, context string, tag string, push bool, name string) error {
	dockerContext := d.getDockerContextForTag(tag)

	args := []string{
		"--context", dockerContext,
		"buildx",
		"build",
		"-f=" + containerFilePath,
		"--tag=" + tag,
	}

	if push {
		args = append(args, "--push")
	}

	args = append(args, context)

	cmd := exec.CommandContext(ctx, "docker", args...)

	cmd.Stdout = &io.PrefixWriter{
		Prefix: fmt.Sprintf("[%s]\t", name),
		Writer: os.Stdout,
	}
	cmd.Stderr = &io.PrefixWriter{
		Prefix: fmt.Sprintf("[%s]\t", name),
		Writer: os.Stderr,
	}

	return cmd.Run()
}

func (DockerBuildxImpl) getDockerContextForTag(tag string) string {
	withoutTag := strings.Split(tag, ":")[0]

	return strings.ReplaceAll(withoutTag, "/", "-")
}
