package util

import (
	"encoding/json"
	"fmt"
	"go-deploy/dto/v2/body"
)

func ProcessUserReadResponse(responseBody string) ([]body.UserRead, error) {
	var users []body.UserRead
	if err := json.Unmarshal([]byte(responseBody), &users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return users, nil
}

func ProcessResponseArr[T any](responseBody string) ([]T, error) {
	var items []T
	if err := json.Unmarshal([]byte(responseBody), &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return items, nil
}

func ProcessResponse[T any](responseBody string) (*T, error) {
	var item T
	if err := json.Unmarshal([]byte(responseBody), &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &item, nil
}

func DeploymentCreateToUpdate(create *body.DeploymentCreate) body.DeploymentUpdate {

	return body.DeploymentUpdate{
		// skipping name here
		CpuCores:        create.CpuCores,
		RAM:             create.RAM,
		Replicas:        create.Replicas,
		Envs:            toPointerSlice(create.Envs),
		Volumes:         toPointerSlice(create.Volumes),
		InitCommands:    toPointerSlice(create.InitCommands),
		Args:            toPointerSlice(create.Args),
		Visibility:      &create.Visibility,
		Private:         &create.Private,
		Image:           create.Image,
		HealthCheckPath: create.HealthCheckPath,
		CustomDomain:    create.CustomDomain,
	}
}

func toPointerSlice[T any](slice []T) *[]T {
	if len(slice) == 0 {
		return nil
	}
	return &slice
}
