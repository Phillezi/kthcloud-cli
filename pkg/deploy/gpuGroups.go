package deploy

import (
	"errors"

	"github.com/kthcloud/go-deploy/dto/v2/body"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
)

func (c *Client) GpuGroupByID(groupID string) (*body.GpuGroupRead, error) {
	if !c.HasValidSession() {
		return nil, errors.New("no active session, log in first")
	}
	req := c.client.R()

	resp, err := req.Get("/v2/gpuGroups/" + groupID)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, errors.New("request to get gpu group returned with status: " + resp.Status())
	}

	gpuGroup, err := util.ProcessResponse[body.GpuGroupRead](resp.String())
	if err != nil {
		return nil, err
	}

	return gpuGroup, nil
}
