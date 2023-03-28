package mappings

import (
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
)

func (m *OpenAPI) ConnectionTestResult(in model.ConnectionResult) openapi.ConnectionResult {
	result := openapi.ConnectionResult{}

	if in.PortCheck.IsSet() {
		result.PortCheck = m.ConnectionTestStep(in.PortCheck)
	}

	if in.Connectivity.IsSet() {
		result.Connectivity = m.ConnectionTestStep(in.Connectivity)
	}

	if in.Authentication.IsSet() {
		result.Authentication = m.ConnectionTestStep(in.Authentication)
	}

	if in.FetchTraces.IsSet() {
		result.FetchTraces = m.ConnectionTestStep(in.FetchTraces)
	}

	return result
}

func (m *OpenAPI) ConnectionTestStep(in model.ConnectionTestStep) openapi.ConnectionTestStep {
	errMessage := ""
	if in.Error != nil {
		errMessage = in.Error.Error()
	}

	return openapi.ConnectionTestStep{
		Passed:  in.Status != model.StatusFailed,
		Message: in.Message,
		Status:  string(in.Status),
		Error:   errMessage,
	}
}
