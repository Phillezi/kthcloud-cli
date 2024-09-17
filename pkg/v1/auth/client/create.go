package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-deploy/dto/v2/body"

	"github.com/go-resty/resty/v2"
)

func (c *Client) Create(data any) (*resty.Response, error) {
	var path string

	// Determine the path based on the type of data
	switch v := data.(type) {
	case *body.DeploymentCreate:
		path = "/v2/deployments"
	case *body.VmCreate:
		path = "/v2/vms"
	default:
		return nil, fmt.Errorf("unsupported data type: %T", v)
	}

	// Serialize the data to JSON
	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("could not marshal data to JSON: %v", err)
	}

	// Perform the POST request
	resp, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bytes.NewReader(body)).
		Post(c.baseURL + path)

	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	return resp, nil
}
