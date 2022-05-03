package selectors

import (
	"github.com/kubeshop/tracetest/traces"
)

type PseudoClass interface {
	Filter(spans []traces.Span) []traces.Span
}

type NthChildPseudoClass struct {
	N int64
}

func (nc NthChildPseudoClass) Filter(spans []traces.Span) []traces.Span {
	if len(spans) < int(nc.N) {
		return []traces.Span{}
	}

	return []traces.Span{spans[int(nc.N-1)]}
}
