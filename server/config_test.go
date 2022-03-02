package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configtls"
)

func TestLoadConfig(t *testing.T) {
	got, err := LoadConfig("testdata/config.yaml")
	assert.NoError(t, err)

	expected := &Config{
		PostgresConnString: "host=postgres user=postgres password=postgres port=5432 sslmode=disable",
		JaegerConnectionConfig: configgrpc.GRPCClientSettings{
			Endpoint:   "jaeger-query:16685",
			TLSSetting: configtls.TLSClientSetting{Insecure: true},
		},
	}
	assert.Equal(t, expected, got)
}
