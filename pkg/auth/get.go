package auth

import "github.com/Phillezi/kthcloud-cli/pkg/session"

func (c *Client) KeycloakBaseURL() string {
	return c.keycloakBaseURL
}

func (c *Client) Session() *session.Session {
	return c.session
}
