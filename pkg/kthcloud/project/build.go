package project

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/kthcloud/cli/internal/body"
	"github.com/kthcloud/cli/pkg/deploy"
	builder "github.com/kthcloud/cli/pkg/image"
	"github.com/kthcloud/go-deploy/models/model"

	// FIXME: update to v3?
	"gopkg.in/yaml.v2"
)

// Build builds all required / specified services
// If buildServices is nil / empty array all services that have build specified will be built
// TODO: HEAD request registry to see if we have to build the image for each buildable service.
func (proj *Project) Build(ctx context.Context, client deploy.ClientWithResponsesInterface, buildServices []string) error {
	buildLen := len(buildServices)
	requiresBuildLen := len(proj.Sevices)
	if buildLen > 0 {
		requiresBuildLen = buildLen
	}
	requiresBuild := make(map[string]Service, requiresBuildLen)

	for name, svc := range proj.Sevices {
		if svc.Build != nil && (buildLen < 1 || slices.Contains(buildServices, svc.Name)) {
			requiresBuild[name] = svc
		}
	}

	if len(requiresBuild) < 1 {
		return nil
	}

	build, err := builder.Get()
	if err != nil {
		return err
	}

	resp, err := client.GetV2DeploymentsWithResponse(ctx, &deploy.GetV2DeploymentsParams{})
	if err != nil {
		return err
	}

	obj, err := deploy.HandleAndAssert[*[]deploy.BodyDeploymentRead](resp, "get")
	if err != nil {
		return err
	}
	if obj == nil {
		return fmt.Errorf("get deployments got nil response")
	}

	userDeployments := make(map[string]deploy.BodyDeploymentRead, len(*obj))
	for _, depl := range *obj {
		userDeployments[*depl.Name] = depl
	}

	for name, svc := range requiresBuild {
		id, err := getOrCreateDeployment(ctx, name, svc.BodyDeploymentCreate, client, userDeployments)
		if err != nil {
			return err
		}

		ci, err := getCiConfigForDeployment(ctx, id, client)
		if err != nil {
			return err
		}

		if ci == nil {
			return fmt.Errorf("ci is nil")
		}

		user, password, tag, err := parseCiConfig(*ci)
		if err != nil {
			return err
		}

		registry, err := parseTag(tag)
		if err != nil {
			return err
		}

		if err := build.Login(ctx, registry, user, password, tag, name); err != nil {
			return err
		}

		if err := build.Build(ctx, svc.Build.ContainerFile, svc.Build.Context, tag, true, name); err != nil {
			return err
		}
	}

	return nil
}

func getOrCreateDeployment(ctx context.Context, name string, cfg deploy.BodyDeploymentCreate, client deploy.ClientWithResponsesInterface, userDeployments map[string]deploy.BodyDeploymentRead) (id string, err error) {
	if existing, found := userDeployments[name]; found {
		// FIXME: check diff in config and update
		if existing.Id != nil {
			return *existing.Id, nil
		}
	}

	return createDeployment(ctx, cfg, client)
}

func createDeployment(ctx context.Context, cfg deploy.BodyDeploymentCreate, client deploy.ClientWithResponsesInterface) (id string, err error) {
	r, err := body.Reader(cfg)
	if err != nil {
		return "", err
	}
	resp, err := client.PostV2DeploymentsWithBodyWithResponse(ctx, "application/json", r)
	if err != nil {
		return "", err
	}

	obj, err := deploy.HandleAndAssert[*deploy.BodyDeploymentRead](resp, "create")
	if err != nil {
		return "", err
	}

	if obj != nil {
		if obj.Id != nil {
			return *obj.Id, nil
		}
	}

	return
}

func getCiConfigForDeployment(ctx context.Context, id string, client deploy.ClientWithResponsesInterface) (*deploy.BodyCiConfig, error) {
	const maxAttempts = 10
	const maxBackoff = 10
	backoff := 0 * time.Second
	for attempt := 0; attempt < maxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(backoff):
			resp, err := client.GetV2DeploymentsDeploymentIdCiConfigWithResponse(ctx, id)
			if err != nil {
				return nil, err
			}

			if resp.StatusCode() != http.StatusOK {
				// TODO: log here
				backoff = time.Duration(max((attempt+1)*2, maxBackoff)) * time.Second
				continue
			}

			obj, err := deploy.HandleAndAssert[*deploy.BodyCiConfig](resp, "get")
			if err != nil {
				return nil, err
			}
			return obj, nil
		}
	}

	return nil, fmt.Errorf(
		"max attempts exhausted when getting ci config")
}

func parseCiConfig(ci deploy.BodyCiConfig) (user, password, tag string, err error) {
	var gha model.GithubActionConfig

	err = yaml.Unmarshal([]byte(*ci.Config), &gha)
	if err != nil {
		return
	}

	if gha.Jobs.Docker.Steps != nil {
		for _, step := range gha.Jobs.Docker.Steps {
			if step.With.Password != "" {
				if user == "" {
					user = step.With.Username
				}
				if password == "" {
					password = step.With.Password
				}

				// gha.Jobs.Docker.Steps[i].With.Password = "${{ secrets.DOCKER_PASSWORD }}"
				// gha.Jobs.Docker.Steps[i].With.Username = "${{ secrets.DOCKER_USERNAME }}"
			}
			if step.With.Tags != "" {
				if tag == "" {
					tag = step.With.Tags
				}
				// gha.Jobs.Docker.Steps[i].With.Tags = "${{ secrets.DOCKER_TAG }}"
			}
		}
	}

	return
}

func parseTag(tag string) (string, error) {
	parts := strings.SplitN(strings.TrimPrefix(strings.TrimPrefix(tag, "https://"), "http://"), "/", 2)
	if len(parts) < 1 {
		return "", fmt.Errorf("invalid tag")
	}
	return parts[0], nil
}
