package config

import "os"

type Config struct {
	OtelExporterEndpoint string
	OtelServiceName      string
}

func Current() *Config {
	return &Config{
		OtelExporterEndpoint: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		OtelServiceName:      os.Getenv("OTEL_SERVICE_NAME"),
	}
}
