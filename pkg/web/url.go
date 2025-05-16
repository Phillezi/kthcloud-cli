package web

import (
	"fmt"
	"net/url"
)

func (c *Server) constructRedirectURI() (string, error) {
	redirectURI := fmt.Sprintf("%s%s", c.redirectHost, c.redirectPath)
	_, err := url.Parse(redirectURI)
	if err != nil {
		return "", err
	}

	return redirectURI, nil
}
