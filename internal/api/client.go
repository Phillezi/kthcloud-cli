package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	apiURL   string
	token    string
	apiToken string
}

func NewClient(apiURL, token string) *Client {
	return &Client{
		apiURL:   apiURL,
		token:    token,
		apiToken: "",
	}
}

func NewAPIClient(apiURL, apiToken string) *Client {
	return &Client{
		apiURL:   apiURL,
		token:    "",
		apiToken: apiToken,
	}
}

func (c *Client) FetchResource(resource string, method string) (string, error) {
	client := resty.New()

	request := client.R()
	if c.apiToken != "" {
		request.SetHeader("X-Api-Key", c.apiToken)
	} else {
		request.SetAuthToken(c.token)
	}
	url := fmt.Sprintf("%s/%s", c.apiURL, resource)

	var resp *resty.Response
	var err error

	switch method {
	case "GET":
		resp, err = request.Get(url)
	case "POST":
		resp, err = request.Post(url)
	case "PUT":
		resp, err = request.Put(url)
	case "DELETE":
		resp, err = request.Delete(url)
	default:
		return "", fmt.Errorf("unsupported method: %s", method)
	}

	if err != nil {
		return "", fmt.Errorf("failed to fetch resource using %s: %w", method, err)
	}

	return resp.String(), nil
}

func (c *Client) Req(resource string, method string, body interface{}) (string, error) {
	client := resty.New()

	request := client.R()
	if c.apiToken != "" {
		request.SetHeader("X-Api-Key", c.apiToken)
	} else {
		request.SetAuthToken(c.token)
	}
	url := fmt.Sprintf("%s/%s", c.apiURL, resource)

	var resp *resty.Response
	var err error

	switch method {
	case "GET":
		resp, err = request.Get(url)
	case "POST":
		resp, err = request.SetBody(body).Post(url)
	case "PUT":
		resp, err = request.SetBody(body).Put(url)
	case "DELETE":
		resp, err = request.Delete(url)
	default:
		return "", fmt.Errorf("unsupported method: %s", method)
	}

	if err != nil {
		return "", fmt.Errorf("failed to fetch resource using %s: %w", method, err)
	}

	return resp.String(), nil
}
