package mappings

import (
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
)

func (m OpenAPI) TestRunEvents(in []model.TestRunEvent) []openapi.TestRunEvent {
	out := make([]openapi.TestRunEvent, 0, len(in))
	for _, event := range in {
		out = append(out, m.TestRunEvent(event))
	}

	return out
}

func (m OpenAPI) TestRunEvent(in model.TestRunEvent) openapi.TestRunEvent {
	return openapi.TestRunEvent{
		Type:                in.Type,
		Stage:               string(in.Stage),
		Description:         in.Description,
		CreatedAt:           in.CreatedAt,
		TestId:              string(in.TestID),
		RunId:               int32(in.RunID),
		DataStoreConnection: m.ConnectionTestResult(in.DataStoreConnection),
		Polling:             m.PollingInfo(in.Polling),
		Outputs:             m.OutputsInfo(in.Outputs),
	}
}

func (m OpenAPI) PollingInfo(in model.PollingInfo) openapi.PollingInfo {
	return openapi.PollingInfo{
		Type:                string(in.Type),
		ReasonNextIteration: in.ReasonNextIteration,
		IsComplete:          in.IsComplete,
		Periodic:            m.PeriodicPollingInfo(in.Periodic),
	}
}

func (m OpenAPI) PeriodicPollingInfo(in *model.PeriodicPollingConfig) openapi.PollingInfoPeriodic {
	if in == nil {
		return openapi.PollingInfoPeriodic{}
	}

	return openapi.PollingInfoPeriodic{
		NumberSpans:      int32(in.NumberSpans),
		NumberIterations: int32(in.NumberIterations),
	}
}

func (m OpenAPI) OutputsInfo(in []model.OutputInfo) []openapi.OutputInfo {
	out := make([]openapi.OutputInfo, 0, len(in))
	for _, output := range in {
		newOutput := openapi.OutputInfo{
			LogLevel:   string(output.LogLevel),
			Message:    output.Message,
			OutputName: output.OutputName,
		}

		out = append(out, newOutput)
	}

	return out
}
