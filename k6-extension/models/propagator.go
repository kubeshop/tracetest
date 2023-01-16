package models

import (
	"fmt"
	"net/http"

	utils "github.com/xoscar/xk6-tracetest-tracing/utils"
)

type PropagatorName string

type Propagator struct {
	propagators []PropagatorName
}

const (
	PropagatorW3C    PropagatorName = "w3c"
	HeaderNameW3C    PropagatorName = "traceparent"
	PropagatorB3     PropagatorName = "b3"
	HeaderNameB3     PropagatorName = "b3"
	PropagatorJaeger PropagatorName = "jaeger"
	HeaderNameJaeger PropagatorName = "uber-trace-id"
)

func NewPropagator(propagators []PropagatorName) Propagator {
	return Propagator{
		propagators: propagators,
	}
}

func (p Propagator) GenerateHeaders(traceID string) http.Header {
	header := http.Header{}

	for _, propagator := range p.propagators {
		switch propagator {
		case PropagatorW3C:
			header.Add(string(HeaderNameW3C), fmt.Sprintf("00-%s-%s-01", traceID, utils.RandHexStringRunes(16)))
		case PropagatorB3:
			header.Add(string(HeaderNameB3), fmt.Sprintf("%s-%s-1", traceID, utils.RandHexStringRunes(8)))
		case PropagatorJaeger:
			header.Add(string(HeaderNameJaeger), fmt.Sprintf("%s:%s:0:1", traceID, utils.RandHexStringRunes(8)))
		}

		// todo: add headers for the missing scenarios
	}

	return header
}
