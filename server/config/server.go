package config

import (
	"fmt"
	"os"
	"strings"
)

var serverOptions = options{
	{
		key:                "postgresConnString",
		defaultValue:       "",
		description:        "Postgres connection string",
		validate:           nil,
		deprecated:         true,
		deprecationMessage: "Use the new postgres config structure instead.",
	},
	{
		key:          "postgres.host",
		defaultValue: "postgres",
		description:  "Postgres DB host",
	},
	{
		key:          "postgres.user",
		defaultValue: "postgres",
		description:  "Postgres DB user",
	},
	{
		key:          "postgres.password",
		defaultValue: "postgres",
		description:  "Postgres DB password",
	},
	{
		key:          "postgres.dbname",
		defaultValue: "tracetest",
		description:  "Postgres DB database name",
	},
	{
		key:          "postgres.port",
		defaultValue: 5432,
		description:  "Postgres DB port",
	},
	{
		key:          "postgres.params",
		defaultValue: "sslmode=disable",
		description:  "Postgres DB connection parameters",
	},
	{
		key:          "server.httpPort",
		defaultValue: 11633,
		description:  "Tracetest server HTTP Port",
	},
	{
		key:          "server.pathPrefix",
		defaultValue: "",
		description:  "Tracetest server HTTP Path prefix",
	},
	{
		key:          "experimentalFeatures",
		defaultValue: []string{},
		description:  "enabled experimental features",
	},
	{
		key:          "internalTelemetry.enabled",
		defaultValue: false,
		description:  "enable internal telemetry (used for internal testing)",
	},
	{
		key:          "internalTelemetry.otelCollectorEndpoint",
		defaultValue: "",
		description:  "internal telemetry  otel collector (used for internal testing)",
		validate:     nil,
	},
	{
		key:          "testPipelines.triggerExecute.enabled",
		defaultValue: "true",
		description:  "enable local trigger execution",
	},
	{
		key:          "testPipelines.traceFetch.enabled",
		defaultValue: "true",
		description:  "enable local trace fetching",
	},
}

func init() {
	configOptions = append(configOptions, serverOptions...)
}

func (c *AppConfig) PostgresConnString() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if postgresConnString := c.vp.GetString("postgresConnString"); postgresConnString != "" {
		fmt.Println("ERROR: postgresConnString was discontinued. Migrate to the new postgres format")
		os.Exit(1)
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s",
		c.vp.GetString("postgres.user"),
		c.vp.GetString("postgres.password"),
		c.vp.GetString("postgres.host"),
		c.vp.GetInt("postgres.port"),
		c.vp.GetString("postgres.dbname"),
		strings.ReplaceAll(c.vp.GetString("postgres.params"), " ", "&"),
	)
}

func (c *AppConfig) ServerPathPrefix() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.vp.GetString("server.pathPrefix")
}

func (c *AppConfig) ServerPort() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.vp.GetInt("server.httpPort")
}

func (c *AppConfig) ExperimentalFeatures() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.vp.GetStringSlice("experimentalFeatures")
}

func (c *AppConfig) InternalTelemetryEnabled() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.vp.GetBool("internalTelemetry.enabled")
}

func (c *AppConfig) InternalTelemetryOtelCollectorAddress() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.vp.GetString("internalTelemetry.otelCollectorEndpoint")
}

func (c *AppConfig) TestPipelineTriggerExecutionEnabled() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	// this config needs to be a string because pflags
	// has a strage bug that ignores this field when
	// it is set as false
	return c.vp.GetString("testPipelines.triggerExecute.enabled") == "true"
}

func (c *AppConfig) TestPipelineTraceFetchingEnabled() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.vp.GetString("testPipelines.traceFetch.enabled") == "true"
}
