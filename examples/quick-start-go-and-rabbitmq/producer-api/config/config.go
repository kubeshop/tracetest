package config

import "os"

type Config struct {
	OtelExporterEndpoint         string
	OtelServiceName              string
	MessageQueueConnectionString string
	MessageQueueName             string
}

func Current() *Config {
	return &Config{
		OtelExporterEndpoint:         os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		OtelServiceName:              os.Getenv("OTEL_SERVICE_NAME"),
		MessageQueueConnectionString: os.Getenv("MESSAGE_QUEUE_CONNECTION_STRING"),
		MessageQueueName:             os.Getenv("MESSAGE_QUEUE_NAME"),
	}
}
