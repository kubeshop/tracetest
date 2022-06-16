package analytics_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/analytics"
	"github.com/stretchr/testify/assert"
)

func TestReadyness(t *testing.T) {
	analytics.Init(false, "serverID", "test", "1.0")
	assert.True(t, analytics.Ready())

	analytics.Init(true, "serverID", "test", "1.0")
	assert.True(t, analytics.Ready())

	analytics.Init(true, "serverID", "", "1.0")
	assert.False(t, analytics.Ready())

	analytics.Init(true, "serverID", "test", "")
	assert.False(t, analytics.Ready())

	analytics.Init(true, "", "test", "1.0")
	assert.False(t, analytics.Ready())

}
