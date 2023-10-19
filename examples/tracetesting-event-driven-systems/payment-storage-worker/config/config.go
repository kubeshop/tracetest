package config

import "os"

type Config struct {
	OtelExporterEndpoint string
	OtelServiceName      string
	KafkaBrokerUrl       string
	KafkaTopic           string
}

func Current() *Config {
	return &Config{
		OtelExporterEndpoint: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		OtelServiceName:      os.Getenv("OTEL_SERVICE_NAME"),
		KafkaBrokerUrl:       os.Getenv("KAFKA_BROKER_URL"),
		KafkaTopic:           os.Getenv("KAFKA_TOPIC"),
	}
}
