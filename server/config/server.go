package config

func (c *Config) PostgresConnString() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.config.PostgresConnString
}

func (c *Config) ServerPathPrefix() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.config.Server.PathPrefix
}

func (c *Config) SetServerPathPrefix(prefix string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.config.Server.PathPrefix = prefix
}

func (c *Config) ServerPort() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.config.Server.HttpPort != 0 {
		return c.config.Server.HttpPort
	}

	return 11633
}

func (c *Config) SetServerPort(port int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.config.Server.HttpPort = port
}

func (c *Config) ExperimentalFeatures() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.config.ExperimentalFeatures
}
