package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func init() {
	defaultSetters = append(defaultSetters, serverDefaultSetter)
}

func serverDefaultSetter(vp *viper.Viper) {
	vp.SetDefault("postgres.host", "postgres")
	vp.SetDefault("postgres.username", "postgres")
	vp.SetDefault("postgres.password", "postgres")
	vp.SetDefault("postgres.dbname", "tracetest")
	vp.SetDefault("postgres.port", 5432)
	vp.SetDefault("postgres.params", "")

	vp.SetDefault("server.port", 11633)
	vp.SetDefault("server.pathPrefix", "/")

	vp.SetDefault("experimentalFeatures", []string{})

	vp.SetDefault("internalTelemetry.enabled", false)
	vp.SetDefault("internalTelemetry.otelCollectorAddress", "")
}

func (c *Config) PostgresConnString() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	str := fmt.Sprintf(
		"host=%s user=%s password=%s port=%d dbname=%s",
		c.vp.GetString("postgres.host"),
		c.vp.GetString("postgres.username"),
		c.vp.GetString("postgres.password"),
		c.vp.GetInt("postgres.port"),
		c.vp.GetString("postgres.dbname"),
	)

	if params := c.vp.GetString("postgres.params"); params != "" {
		str += " " + params
	}

	return str
}

func (c *Config) ServerPathPrefix() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.vp.GetString("server.pathPrefix")
}

func (c *Config) SetServerPathPrefix(prefix string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.vp.Set("server.pathPrefix", prefix)
}

func (c *Config) ServerPort() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.vp.GetInt("server.port")
}

func (c *Config) SetServerPort(port int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.vp.Set("server.port", port)
}

func (c *Config) ExperimentalFeatures() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.vp.GetStringSlice("experimentalFeatures")
}

func (c *Config) InternalTelemetryEnabled() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.vp.GetBool("internalTelemetry.enabled")
}

func (c *Config) InternalTelemetryOtelCollectorAddress() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.vp.GetString("internalTelemetry.otelCollectorAddress")
}
