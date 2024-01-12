package connection_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
	"github.com/stretchr/testify/assert"
)

func TestPortLinter(t *testing.T) {
	testCases := []struct {
		Name           string
		Endpoints      []string
		ExpectedPorts  []string
		ExpectedStatus model.Status
	}{
		{
			Name:           "shouldSucceedIfPortIsExpected",
			Endpoints:      []string{"jaeger:16685"},
			ExpectedPorts:  []string{"16685"},
			ExpectedStatus: model.StatusPassed,
		},
		{
			Name:           "shouldShowWarningInCaseOfDifferentPort",
			Endpoints:      []string{"jaeger:16686"},
			ExpectedPorts:  []string{"16685"},
			ExpectedStatus: model.StatusWarning,
		},
		{
			Name:           "shouldSupportSchemas",
			Endpoints:      []string{"https://us2.endpoint:9200"},
			ExpectedPorts:  []string{"9200"},
			ExpectedStatus: model.StatusPassed,
		},
		{
			Name:           "shouldSupportTwoPorts",
			Endpoints:      []string{"https://us2.endpoint:9100"},
			ExpectedPorts:  []string{"9200", "9250"},
			ExpectedStatus: model.StatusWarning,
		},
		{
			Name:           "shouldSupportTwoPorts",
			Endpoints:      []string{"https://us2.endpoint:9100"},
			ExpectedPorts:  []string{"9200", "9250", "9300"},
			ExpectedStatus: model.StatusWarning,
		},
		{
			Name:           "shouldGetPort80InHTTP",
			Endpoints:      []string{"http://tempo-us-central1.grafana.net"},
			ExpectedPorts:  []string{"80"},
			ExpectedStatus: model.StatusPassed,
		},
		{
			Name:           "shouldGetPort443InHTTPS",
			Endpoints:      []string{"https://tempo-us-central1.grafana.net"},
			ExpectedPorts:  []string{"443"},
			ExpectedStatus: model.StatusPassed,
		},
		{
			Name:           "shouldGetSpecifiedInHTTP",
			Endpoints:      []string{"http://tempo-us-central1.grafana.net:8081"},
			ExpectedPorts:  []string{"8081"},
			ExpectedStatus: model.StatusPassed,
		},
		{
			Name:           "shouldGetPort443InHTTPS",
			Endpoints:      []string{"https://tempo-us-central1.grafana.net:8082"},
			ExpectedPorts:  []string{"8082"},
			ExpectedStatus: model.StatusPassed,
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
