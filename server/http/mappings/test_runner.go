package mappings

import (
	"github.com/kubeshop/tracetest/server/executor/testrunner"
	"github.com/kubeshop/tracetest/server/openapi"
)

func (m OpenAPI) RequiredGatesResult(in testrunner.RequiredGatesResult) openapi.RequiredGatesResult {
	parsedRequired := make([]openapi.SupportedGates, len(in.Required))
	for i, gate := range in.Required {
		parsedRequired[i] = openapi.SupportedGates(gate)
	}

	parsedFailed := make([]openapi.SupportedGates, len(in.Failed))
	for i, gate := range in.Failed {
		parsedFailed[i] = openapi.SupportedGates(gate)
	}

	return openapi.RequiredGatesResult{
		Passed:   in.Passed,
		Required: parsedRequired,
		Failed:   parsedFailed,
	}
}

func (m Model) RequiredGates(in *[]openapi.SupportedGates) *[]testrunner.RequiredGate {
	if in == nil {
		return nil
	}

	parsedGates := make([]testrunner.RequiredGate, len(*in))

	for i, gate := range *in {
		parsedGates[i] = testrunner.RequiredGate(gate)
	}

	return &parsedGates
}
