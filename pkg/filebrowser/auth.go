package filebrowser

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *Client) Auth() (bool, error) {
	_, err := c.loadCookies()
	if err != nil {
		// log err here
		return false, err
	}

	resp, err := c.client.Get(c.filebrowserURL + "/oauth2/auth")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 202 {
		return false, fmt.Errorf("request error: %s", resp.Status)
	}

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
