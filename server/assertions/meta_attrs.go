package assertions

import (
	"errors"
	"strconv"

	"github.com/kubeshop/tracetest/server/traces"
)

var (
	errMetaAttrNotDefined = errors.New("meta attribute not defined")
	metaAttrsRegistry     = map[string]MetaAttribute{
		"spans_collection.count": count{},
	}
)

func metaAttr(name string) (MetaAttribute, error) {
	a, ok := metaAttrsRegistry[name]
	if !ok {
		return nil, errMetaAttrNotDefined
	}

	return a, nil
}

type MetaAttribute interface {
	Value(spans []traces.Span) string
}

var _ MetaAttribute = count{}

type count struct{}

func (c count) Value(spans []traces.Span) string {
	return strconv.Itoa(len(spans))
}
