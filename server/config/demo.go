package config

import "strings"

func (c *Config) DemoEnabled() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.config.Demo.Enabled
}

func (c *Config) DemoEndpoints() map[string]string {
	c.mu.Lock()
	defer c.mu.Unlock()

	fixed := map[string]string{}
	for k, v := range c.config.Demo.Endpoints {
		// uppercase first letter
		fixed[strings.ToUpper(k[:1])+k[1:]] = v
	}

	return fixed
}
