package analytics_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/analytics"
	"github.com/kubeshop/tracetest/server/config"
	"github.com/stretchr/testify/assert"
)

func TestReadyness(t *testing.T) {
	analytics.Init(config.GoogleAnalytics{Enabled: false}, "test", "1.0")
	assert.True(t, analytics.Ready())

	analytics.Init(config.GoogleAnalytics{Enabled: true}, "test", "1.0")
	assert.False(t, analytics.Ready())

	analytics.Init(config.GoogleAnalytics{
		Enabled:       true,
		MeasurementID: "1",
		SecretKey:     "2",
	}, "test", "1.0")
	assert.True(t, analytics.Ready())

	analytics.Init(config.GoogleAnalytics{
		Enabled:       true,
		MeasurementID: "1",
		SecretKey:     "2",
	}, "test", "1.0")
	assert.True(t, analytics.Ready())

	analytics.Init(config.GoogleAnalytics{
		Enabled:       true,
		MeasurementID: "1",
		SecretKey:     "2",
	}, "", "1.0")
	assert.False(t, analytics.Ready())

	analytics.Init(config.GoogleAnalytics{
		Enabled:       true,
		MeasurementID: "1",
		SecretKey:     "2",
	}, "test", "")
	assert.False(t, analytics.Ready())

}
