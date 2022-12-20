package mappings

import (
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/tracedb"
)

func (m *OpenAPI) ConnectionTestResult(in tracedb.ConnectionTestResult) openapi.ConnectionResult {
	result := openapi.ConnectionResult{}

	if in.ConnectivityTestResult.IsSet() {
		result.Connectivity = m.ConnectionTestStep(in.ConnectivityTestResult)
	}

	if in.AuthenticationTestResult.IsSet() {
		result.Authentication = m.ConnectionTestStep(in.AuthenticationTestResult)
	}

	if in.TraceRetrivalTestResult.IsSet() {
		result.FetchTraces = m.ConnectionTestStep(in.TraceRetrivalTestResult)
	}

	return result
}

func (m *OpenAPI) ConnectionTestStep(in tracedb.ConnectionTestStepResult) openapi.ConnectionTestStep {
	errMessage := ""
	if in.Error != nil {
		errMessage = in.Error.Error()
	}

	return openapi.ConnectionTestStep{
		Passed:  in.Error == nil,
		Message: in.OperationDescription,
		Error:   errMessage,
	}
}
