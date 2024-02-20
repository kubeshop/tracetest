package expression

import (
	"fmt"
	"strconv"

	"github.com/kubeshop/tracetest/server/traces"
	"github.com/kubeshop/tracetest/server/variableset"
)

type DataStore interface {
	Source() string
	Get(name string) (string, error)
}

var attributeAlias = map[string]string{
	"name": "tracetest.span.name",
}

type AttributeDataStore struct {
	Span traces.Span
}

func (ds AttributeDataStore) Source() string {
	return "attr"
}

func (ds AttributeDataStore) getFromAlias(name string) (string, error) {
	alias, found := attributeAlias[name]

	if !found {
		return "", fmt.Errorf(`attribute "%s" not found`, name)
	}

	value := ds.Span.Attributes.Get(alias)
	if value == "" {
		return "", fmt.Errorf(`attribute "%s" not found`, name)
	}

	return value, nil
}

func (ds AttributeDataStore) Get(name string) (string, error) {
	if !ds.Span.ID.IsValid() {
		// It's probably a nil span and we never got a matching selector,
		// so instead of returning a resolution error, let's return a non-matching
		// span error instead
		return "", fmt.Errorf(`there are no matching spans to retrieve the attribute "%s" from`, name)
	}

	value := ds.Span.Attributes.Get(name)
	if value == "" {
		return ds.getFromAlias(name)
	}

	return value, nil
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

type VariableDataStore struct {
	Values []variableset.VariableSetValue
}

func (ds VariableDataStore) Source() string {
	return "env"
}

func (ds VariableDataStore) Get(name string) (string, error) {
	for _, v := range ds.Values {
		if v.Key == name {
			return v.Value, nil
		}
	}

	return "", fmt.Errorf(`variable "%s" not found`, name)
}
