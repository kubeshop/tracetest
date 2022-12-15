package selectors

import "github.com/kubeshop/tracetest/server/model"

type PseudoClass interface {
	Name() string
	Filter(spans []model.Span) []model.Span
}

type NthChildPseudoClass struct {
	N int64
}

func (nc NthChildPseudoClass) Name() string {
	return "nth_child"
}

func (nc NthChildPseudoClass) Filter(spans []model.Span) []model.Span {
	if int(nc.N) < 1 || int(nc.N) > len(spans) {
		return []model.Span{}
	}

	return []model.Span{spans[int(nc.N-1)]}
}

type FirstPseudoClass struct{}

func (fpc FirstPseudoClass) Name() string {
	return "first"
}

func (fpc FirstPseudoClass) Filter(spans []model.Span) []model.Span {
	if len(spans) == 0 {
		return []model.Span{}
	}

	return []model.Span{spans[0]}
}

type LastPseudoClass struct{}

func (lpc LastPseudoClass) Name() string {
	return "last"
}

func (lpc LastPseudoClass) Filter(spans []model.Span) []model.Span {
	length := len(spans)
	if length == 0 {
		return []model.Span{}
	}

	return []model.Span{spans[length-1]}
}
