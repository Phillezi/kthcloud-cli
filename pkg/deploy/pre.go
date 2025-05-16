package deploy

func (c *Client) HasValidSession() bool {
	if c.session == nil && c.authClient != nil {
		c.session = c.authClient.Session()
		c.client.SetAuthToken(c.session.Token.AccessToken)
	}
	return c.session != nil && !c.session.IsExpired()
}
