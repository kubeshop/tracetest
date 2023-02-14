package config_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/stretchr/testify/assert"
)

func TestExporter(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		expectedExporter := &config.TelemetryExporterOption{
			ServiceName: "tracetest",
			Sampling:    100,
			Exporter: config.ExporterConfig{
				Type: "collector",
				CollectorConfiguration: config.OTELCollectorConfig{
					Endpoint: "collector:8888",
				},
			},
		}

		cfg := configFromFile(t, "./testdata/exporter_config.yaml")
		selectedExporter, err := cfg.Exporter()
		assert.NoError(t, err)
		assert.Equal(t, expectedExporter, selectedExporter)

		selectedAppExporter, err := cfg.ApplicationExporter()
		assert.NoError(t, err)
		assert.Equal(t, expectedExporter, selectedAppExporter)
	})

	t.Run("Inexistent", func(t *testing.T) {
		cfg := configFromFile(t, "./testdata/inexistent_exporter.yaml")

		exporter, err := cfg.Exporter()
		assert.Error(t, err)
		assert.Nil(t, exporter)

		appExporter, err := cfg.ApplicationExporter()
		assert.Error(t, err)
		assert.Nil(t, appExporter)
	})

	t.Run("Empty", func(t *testing.T) {
		cfg := configFromFile(t, "./testdata/empty_exporter.yaml")

		exporter, err := cfg.Exporter()
		assert.NoError(t, err)
		assert.Nil(t, exporter)

		appExporter, err := cfg.ApplicationExporter()
		assert.NoError(t, err)
		assert.Nil(t, appExporter)
	})
}
