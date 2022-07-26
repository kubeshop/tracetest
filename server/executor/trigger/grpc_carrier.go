package trigger

import (
	"github.com/kubeshop/tracetest/server/model"
	"go.opentelemetry.io/otel/propagation"
)

type GRPCHeaderCarrier struct {
	headers *[]model.GRPCHeader
}

func NewGRPCHeaderCarrier(headers *[]model.GRPCHeader) GRPCHeaderCarrier {
	return GRPCHeaderCarrier{headers: headers}
}

func (c GRPCHeaderCarrier) Get(key string) string {
	for _, header := range *c.headers {
		if header.Key == key {
			return header.Value
		}
	}

	return ""
}

func (c GRPCHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(*c.headers))
	for _, header := range *c.headers {
		keys = append(keys, header.Key)
	}

	return keys
}

// Set implements propagation.TextMapCarrier
func (c GRPCHeaderCarrier) Set(key string, value string) {
	valueSet := false
	for i, header := range *c.headers {
		if header.Key == key {
			header.Value = value
			(*c.headers)[i] = header
			valueSet = true
		}
	}

	if !valueSet {
		*c.headers = append(*c.headers, model.GRPCHeader{
			Key:   key,
			Value: value,
		})
	}
}

var _ propagation.TextMapCarrier = &GRPCHeaderCarrier{}
