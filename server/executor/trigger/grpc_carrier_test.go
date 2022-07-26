package trigger_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/executor/trigger"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/stretchr/testify/assert"
)

func TestGrpcCarrierSet(t *testing.T) {
	headers := make([]model.GRPCHeader, 0)
	carrier := trigger.NewGRPCHeaderCarrier(&headers)

	carrier.Set("my-key", "my-value")

	assert.Len(t, headers, 1)
	assert.Equal(t, "my-value", headers[0].Value)
}

func TestGrpcCarrierSetWithExistingKey(t *testing.T) {
	headers := []model.GRPCHeader{
		{Key: "existing-key", Value: "old value"},
	}
	carrier := trigger.NewGRPCHeaderCarrier(&headers)

	carrier.Set("existing-key", "new value")

	assert.Len(t, headers, 1)
	assert.Equal(t, "new value", headers[0].Value)
}

func TestGrpcCarrierKeys(t *testing.T) {
	headers := []model.GRPCHeader{
		{Key: "key 1", Value: "value 1"},
		{Key: "key 2", Value: "value 2"},
		{Key: "key 3", Value: "value 3"},
	}

	carrier := trigger.NewGRPCHeaderCarrier(&headers)
	keys := carrier.Keys()

	assert.Len(t, keys, 3)
}

func TestGrpcCarrierGet(t *testing.T) {
	headers := []model.GRPCHeader{
		{Key: "key 1", Value: "value 1"},
		{Key: "key 2", Value: "value 2"},
		{Key: "key 3", Value: "value 3"},
	}

	carrier := trigger.NewGRPCHeaderCarrier(&headers)
	value := carrier.Get("key 2")

	assert.Equal(t, "value 2", value)
}
