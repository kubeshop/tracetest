package mappings

import (
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
)

func (m *OpenAPI) ConnectionTestResult(in connection.ConnectionTestResult) openapi.ConnectionResult {
	result := openapi.ConnectionResult{}

	if in.ConnectivityTestResult.IsSet() {
		result.Connectivity = m.ConnectionTestStep(in.ConnectivityTestResult)
	}

	if in.AuthenticationTestResult.IsSet() {
		result.Authentication = m.ConnectionTestStep(in.AuthenticationTestResult)
	}

	if in.TraceRetrievalTestResult.IsSet() {
		result.FetchTraces = m.ConnectionTestStep(in.TraceRetrievalTestResult)
	}

	return result
}

func (m *OpenAPI) ConnectionTestStep(in connection.ConnectionTestStepResult) openapi.ConnectionTestStep {
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
