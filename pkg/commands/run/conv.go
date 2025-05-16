package run

import (
	"fmt"
	"strings"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/kthcloud/go-deploy/dto/v2/body"
)

func ConvertToDeploymentCreate(
	name string,
	envs map[string]string,
	port []int,
	visibility string,
	image string,
	memory, cores float64,
	replicas int,
) *body.DeploymentCreate {
	envList := make([]body.Env, 0, len(envs))
	for k, v := range envs {
		envList = append(envList, body.Env{Name: k, Value: v})
	}
	if len(port) > 0 {
		envList = append(envList, body.Env{Name: "PORT", Value: fmt.Sprintf("%d", port[0])})
		if len(port) > 1 {
			envList = append(envList, body.Env{
				Name: "INTERNAL_PORTS",
				Value: strings.Trim(
					strings.Join(
						strings.Fields(
							fmt.Sprint(
								port[1:],
							)),
						",",
					), "[]")})
		}
	}

	var img *string
	if strings.TrimSpace(image) != "" {
		img = &image
	}

	return &body.DeploymentCreate{
		Name:       name,
		CpuCores:   &cores,
		RAM:        &memory,
		Replicas:   &replicas,
		Envs:       envList,
		Visibility: util.Or(visibility, "public"),
		Image:      img,
	}
}
