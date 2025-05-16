package deploy

func (c *Client) BaseURL() string {
	return c.baseURL
}

func (c *Client) KeycloakBaseURL() string {
	return c.authClient.KeycloakBaseURL()
}
