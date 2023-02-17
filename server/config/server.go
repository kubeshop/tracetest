package config

import (
	"fmt"
)

var serverOptions = options{
	{
		key:          "postgres.host",
		defaultValue: "postgres",
		description:  "Postgres DB host",
		validate:     nil,
	},
	{
		key:          "postgres.user",
		defaultValue: "postgres",
		description:  "Postgres DB user",
		validate:     nil,
	},
	{
		key:          "postgres.password",
		defaultValue: "postgres",
		description:  "Postgres DB password",
		validate:     nil,
	},
	{
		key:          "postgres.dbname",
		defaultValue: "tracetest",
		description:  "Postgres DB database name",
		validate:     nil,
	},
	{
		key:          "postgres.port",
		defaultValue: 5432,
		description:  "Postgres DB port",
		validate:     nil,
	},
	{
		key:          "postgres.params",
		defaultValue: "sslmode=disable",
		description:  "Postgres DB connection parameters",
		validate:     nil,
	},
	{
		key:          "server.httpPort",
		defaultValue: 11633,
		description:  "Tracetest server HTTP Port",
		validate:     nil,
	},
	{
		key:          "server.pathPrefix",
		defaultValue: "",
		description:  "Tracetest server HTTP Path prefix",
		validate:     nil,
	},
	{
		key:          "experimentalFeatures",
		defaultValue: []string{},
		description:  "enabled experimental features",
		validate:     nil,
	},
	{
		key:          "internalTelemetry.enabled",
		defaultValue: false,
		description:  "enable internal telemetry (used for internal testing)",
		validate:     nil,
	},
	{
		key:          "internalTelemetry.otelCollectorEndpoint",
		defaultValue: "",
		description:  "internal telemetry  otel collector (used for internal testing)",
		validate:     nil,
	},
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
