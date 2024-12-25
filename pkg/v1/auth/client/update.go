package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-deploy/dto/v2/body"

	"github.com/go-resty/resty/v2"
)

func (c *Client) Update(data any, id string) (*resty.Response, error) {
	var path string

	switch v := data.(type) {
	case *body.DeploymentUpdate:
		path = "/v2/deployments/" + id
	case *body.VmUpdate:
		path = "/v2/vms/" + id
	default:
		return nil, fmt.Errorf("unsupported data type: %T", v)
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("could not marshal data to JSON: %v", err)
	}

	resp, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bytes.NewReader(body)).
		Post(c.baseURL + path)

	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	return resp, nil
}
