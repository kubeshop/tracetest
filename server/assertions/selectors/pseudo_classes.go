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
	if int(nc.N) < 1 || int(nc.N) > len(spans) {
		return []traces.Span{}
	}

	return []traces.Span{spans[int(nc.N-1)]}
}

type FirstPseudoClass struct{}

func (fpc FirstPseudoClass) Filter(spans []traces.Span) []traces.Span {
	if len(spans) == 0 {
		return []traces.Span{}
	}

	return []traces.Span{spans[0]}
}

type LastPseudoClass struct{}

func (lpc LastPseudoClass) Filter(spans []traces.Span) []traces.Span {
	length := len(spans)
	if length == 0 {
		return []traces.Span{}
	}

	return []traces.Span{spans[length-1]}
}
