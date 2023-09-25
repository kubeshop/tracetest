package analytics_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/analytics"
	"github.com/stretchr/testify/assert"
)

func TestReadyness(t *testing.T) {
	analytics.Init(false, "serverID", "1.0", "env", "123", "123")
	assert.True(t, analytics.Ready())

	analytics.Init(true, "serverID", "1.0", "env", "123", "123")
	assert.True(t, analytics.Ready())

	analytics.Init(true, "serverID", "", "env", "", "")
	assert.False(t, analytics.Ready())

	analytics.Init(true, "", "1.0", "env", "", "")
	assert.False(t, analytics.Ready())

	analytics.Init(true, "serverID", "1.0", "", "", "")
	assert.False(t, analytics.Ready())

}
