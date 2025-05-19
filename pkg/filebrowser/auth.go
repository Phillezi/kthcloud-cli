package filebrowser

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *Client) Auth() (bool, error) {
	if c.token == "" {
		requestBody := strings.NewReader(`{"username": "", "password": "", "recaptcha": ""}`)

		req, err := http.NewRequest(
			"POST",
			c.filebrowserURL+"/api/login",
			requestBody,
		)
		if err != nil {
			return false, err
		}

		req.Header.Set("Content-Type", "application/json")
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
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return false, err
		}

		c.token = string(bodyBytes)
	}

	return true, nil
}
