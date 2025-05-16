package deploy

import (
	"errors"
)

func (c *Client) Token() (string, error) {
	if !c.HasValidSession() {
		return "", errors.New("no active session, log in first")
	}

	if c.session.ApiKey != nil {
		return c.session.ApiKey.Key, nil
	}

	return c.session.Token.AccessToken, nil
}

func (c *Client) ApiKey() (string, error) {
	if !c.HasValidSession() {
		return "", errors.New("no active session, log in first")
	}

	if c.session.ApiKey == nil {
		return "", errors.New("no api key available")
	}
	return c.session.ApiKey.Key, nil
}
