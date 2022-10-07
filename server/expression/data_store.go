package expression

import "github.com/kubeshop/tracetest/server/traces"

type DataStore interface {
	Get(name string) (string, error)
}

type AttributeDataStore struct {
	Span traces.Span
}

func (ads AttributeDataStore) Get(name string) (string, error) {
	return ads.Span.Attributes.Get(name), nil
}
