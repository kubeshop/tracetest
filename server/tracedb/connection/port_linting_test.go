package connection_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/tracedb/connection"
	"github.com/stretchr/testify/assert"
)

func TestPortLinter(t *testing.T) {
	testCases := []struct {
		Name           string
		Endpoints      []string
		ExpectedPorts  []string
		ExpectedStatus connection.Status
	}{
		{
			Name:           "shouldSucceedIfPortIsExpected",
			Endpoints:      []string{"jaeger:16685"},
			ExpectedPorts:  []string{"16685"},
			ExpectedStatus: connection.StatusPassed,
		},
		{
			Name:           "shouldShowWarningInCaseOfDifferentPort",
			Endpoints:      []string{"jaeger:16686"},
			ExpectedPorts:  []string{"16685"},
			ExpectedStatus: connection.StatusWarning,
		},
		{
			Name:           "shouldSupportSchemas",
			Endpoints:      []string{"https://us2.endpoint:9200"},
			ExpectedPorts:  []string{"9200"},
			ExpectedStatus: connection.StatusPassed,
		},
		{
			Name:           "shouldSupportTwoPorts",
			Endpoints:      []string{"https://us2.endpoint:9100"},
			ExpectedPorts:  []string{"9200", "9250"},
			ExpectedStatus: connection.StatusWarning,
		},
		{
			Name:           "shouldSupportTwoPorts",
			Endpoints:      []string{"https://us2.endpoint:9100"},
			ExpectedPorts:  []string{"9200", "9250", "9300"},
			ExpectedStatus: connection.StatusWarning,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			linter := connection.PortLinter("Jaeger", testCase.ExpectedPorts, testCase.Endpoints...)
			result := linter.TestConnection(context.Background())
			assert.Equal(t, testCase.ExpectedStatus, result.Status)
		})
	}
}
