package ps

func (c *Command) getGPUName(gpuNames *map[string]string, groupID string) (string, error) {
	if name, ok := (*gpuNames)[groupID]; ok {
		return name, nil
	}

	group, err := c.client.GpuGroupByID(groupID)
	if err != nil {
		return "", err
	}

	return group.DisplayName, nil
}
