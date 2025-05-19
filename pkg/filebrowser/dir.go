package filebrowser

import (
	"fmt"
	"net/http"
)

func (c *Client) CreateDir(filePath string) (bool, error) {
	if c.token == "" {
		_, err := c.Auth()
		if err != nil {
			return false, err
		}
	}

	endpointURL := fmt.Sprintf("%s/api/resources/%s/?override=false", c.filebrowserURL, filePath)
	req, err := http.NewRequest(
		"POST",
		endpointURL,
		nil,
	)
	if err != nil {
		return false, err
	}

	req.Header.Set("X-Auth", c.token)
	req.Header.Set("Content-Type", "text/plain;charset=UTF-8")

	if c.session != nil && c.session.Token.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.session.Token.AccessToken)
	} else {
		return false, fmt.Errorf("no active session")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("request error: %s", resp.Status)
	}

	return true, nil
}
