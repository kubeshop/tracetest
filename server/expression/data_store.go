package expression

import (
	"fmt"
	"strconv"

	"github.com/kubeshop/tracetest/server/traces"
)

type DataStore interface {
	Source() string
	Get(name string) (string, error)
}

type AttributeDataStore struct {
	Span traces.Span
}

func (ds AttributeDataStore) Source() string {
	return "attr"
}

func (ds AttributeDataStore) Get(name string) (string, error) {
	return ds.Span.Attributes.Get(name), nil
}

type MetaAttributesDataStore struct {
	SelectedSpans []traces.Span
}

func (ds MetaAttributesDataStore) Source() string {
	return metaPrefix
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

type VariableDataStore map[string]string

func (ds VariableDataStore) Source() string {
	return "var"
}

func (ds VariableDataStore) Get(name string) (string, error) {
	value, found := ds[name]
	if !found {
		return "", fmt.Errorf(`variable "%s" is not set`, name)
	}

	return value, nil
}
