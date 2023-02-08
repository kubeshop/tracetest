package config

func (c *Config) DemoEnabled() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.config.Demo.Enabled
}

func (c *Config) DemoEndpoints() map[string]string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.config.Demo.Endpoints
}
