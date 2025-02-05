package client

import (
	"fmt"

	"github.com/kthcloud/go-deploy/dto/v2/body"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

func (c *Client) Remove(data any) (*resty.Response, error) {
	var path string

	switch v := data.(type) {
	case *body.DeploymentRead:
		path = "/v2/deployments/" + v.ID
	case *body.VmRead:
		path = "/v2/vms/" + v.ID
	default:
		return nil, fmt.Errorf("unsupported data type: %T", v)
	}

	logrus.Debugln("removing deployment:", path)
	resp, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		Delete(c.baseURL + path)

	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	return resp, nil
}
