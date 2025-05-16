package auth

func (c *Client) Logout() error {
	c.session = nil
	return c.session.Save(c.sessionPath)
}
