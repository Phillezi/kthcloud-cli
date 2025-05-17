package deploy

func (c *Client) HasValidSession() bool {
	if c.session == nil && c.authClient != nil {
		c.session = c.authClient.Session()
	}
	return c.session != nil && !c.session.IsExpired()
}
