package deploy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/kthcloud/go-deploy/dto/v2/body"
)

func (c *Client) Create(data any) (*resty.Response, error) {
	if !c.HasValidSession() {
		return nil, errors.New("no active session, log in first")
	}
	var path string

	switch v := data.(type) {
	case *body.DeploymentCreate, body.DeploymentCreate:
		path = "/v2/deployments"
	case *body.VmCreate, body.VmCreate:
		path = "/v2/vms"
	case *body.ApiKeyCreate, body.ApiKeyCreate:
		user, err := c.User()
		if err != nil {
			return nil, err
		}
		path = "/v2/users/" + user.ID + "/apiKeys"
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
