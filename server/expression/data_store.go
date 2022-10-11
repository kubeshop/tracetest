package expression

import (
	"fmt"
	"strconv"

	"github.com/kubeshop/tracetest/server/traces"
)

type DataStore interface {
	Get(name string) (string, error)
}

type AttributeDataStore struct {
	Span traces.Span
}

func (ds AttributeDataStore) Get(name string) (string, error) {
	return ds.Span.Attributes.Get(name), nil
}

type MetaAttributesDataStore struct {
	SelectedSpans []traces.Span
}

func (ds MetaAttributesDataStore) Get(name string) (string, error) {
	switch name {
	case "count":
		return ds.count(), nil
	}
	return "", fmt.Errorf("unknown meta attribute %s%s", metaPrefix, name)
}

func (ds MetaAttributesDataStore) count() string {
	return strconv.Itoa(len(ds.SelectedSpans))
}
