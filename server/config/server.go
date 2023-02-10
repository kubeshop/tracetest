package config

import (
	"fmt"
)

var serverOptions = options{
	{"postgres.host", "postgres", "postgres DB host"},
	{"postgres.user", "postgres", "postgres DB user"},
	{"postgres.password", "postgres", "postgres DB password"},
	{"postgres.dbname", "tracetest", "postgres DB dbname"},
	{"postgres.port", 5432, "postgres DB port"},
	{"postgres.params", "sslmode=disable", "postgres DB connection params"},

	{"server.httpPort", 11633, "tracetest HTTP Port"},
	{"server.pathPrefix", "/", "tracetest HTTP Path prefix"},

	{"experimentalFeatures", []string{}, "enabled experimental features"},

	{"internalTelemetry.enabled", false, "enable internal telemetry (used for internal testing)"},
	{"internalTelemetry.otelCollectorEndpoint", "", "internal telemetry  otel collector (used for internal testing)"},
}

func init() {
	configOptions = append(configOptions, serverOptions...)
}

func (c *Config) PostgresConnString() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	str := fmt.Sprintf(
		"host=%s user=%s password=%s port=%d dbname=%s",
		c.vp.GetString("postgres.host"),
		c.vp.GetString("postgres.user"),
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

func (c *Config) ServerPort() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.vp.GetInt("server.httpPort")
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

	return c.vp.GetString("internalTelemetry.otelCollectorEndpoint")
}
