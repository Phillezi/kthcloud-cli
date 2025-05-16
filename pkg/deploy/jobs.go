package deploy

import (
	"errors"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/kthcloud/go-deploy/dto/v2/body"
)

func (c *Client) JobByID(jobID string) (*body.JobRead, error) {
	if !c.HasValidSession() {
		return nil, errors.New("no active session, log in first")
	}
	req := c.client.R()

	resp, err := req.Get("/v2/jobs/" + jobID)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, errors.New("request to get job returned with status: " + resp.Status())
	}

	gpuGroup, err := util.ProcessResponse[body.JobRead](resp.String())
	if err != nil {
		return nil, err
	}

	return gpuGroup, nil
}
