package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

func (c *Client) fetchOAuthToken(redirectURI, code string) (*http.Response, error) {
	tokenURL, err := c.constructOAuthTokenURL()
	if err != nil {
		return nil, err
	}

	logrus.Debug(tokenURL)
	logrus.Debug(redirectURI)
	logrus.Debug(code)
	logrus.Debug(c.keycloakClientID)

	form := url.Values{}
	form.Add("client_id", c.keycloakClientID)
	form.Add("redirect_uri", redirectURI)
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)

	tokenCtx, cancelTokenCtx := context.WithTimeout(c.ctx, c.requestTimeout)
	defer cancelTokenCtx()
	req, err := http.NewRequestWithContext(tokenCtx,
		"POST",
		tokenURL,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	return resp, nil
}
