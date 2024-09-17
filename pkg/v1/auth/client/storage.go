package client

import (
	"errors"
	"fmt"
	"strings"
)

func (c *Client) StorageAuth() (bool, error) {
	user, err := c.User()
	if err != nil {
		return false, err
	}
	req := c.client.R()

	resp, err := req.Get(fmt.Sprintf("%s/oauth2/auth", *user.StorageURL))
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, fmt.Errorf("request error: %s", resp.Status())
	}

	return true, nil
}

func (c *Client) StorageCreateDir(filePath string) (bool, error) {
	user, err := c.User()
	if err != nil {
		return false, err
	}
	req := c.client.R()

	resp, err := req.Post(fmt.Sprintf("%s/api/resources/%s/?override=false", *user.StorageURL, filePath))
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, fmt.Errorf("request error: %s", resp.Status())
	}

	return true, nil
}

func (c *Client) StorageCreateFile(filePath string, content []byte) (bool, error) {
	user, err := c.User()
	if err != nil {
		return false, err
	}
	req := c.client.R()

	// 1. Create initial file
	resp, err := req.Post(fmt.Sprintf("%s/api/tus/%s/?override=false", *user.StorageURL, filePath))
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, fmt.Errorf("request error: %s", resp.Status())
	}

	// 2. Check contents of file to make sure it is empty
	resp, err = req.Head(fmt.Sprintf("%s/api/tus/%s/?override=false", *user.StorageURL, filePath))
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, fmt.Errorf("request error: %s", resp.Status())
	}

	if len(strings.Trim(resp.String(), " \n")) > 0 {
		return false, errors.New("file not empty")
	}

	// 3. Upload contents
	req.Body = content
	resp, err = req.Patch(fmt.Sprintf("%s/api/tus/%s/?override=false", *user.StorageURL, filePath))
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, fmt.Errorf("request error: %s", resp.Status())
	}

	return true, nil
}
