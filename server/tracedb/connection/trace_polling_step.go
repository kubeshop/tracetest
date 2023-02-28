package connection

import (
	"context"
	"errors"

	"github.com/kubeshop/tracetest/server/model"
	"go.opentelemetry.io/otel/trace"
)

type DataStore interface {
	GetTraceByID(context.Context, string) (model.Trace, error)
	GetTraceID() trace.TraceID
}

type tracePollingTestStep struct {
	dataStore DataStore
}

func (s *tracePollingTestStep) TestConnection(ctx context.Context) ConnectionTestStepResult {
	_, err := s.dataStore.GetTraceByID(ctx, s.dataStore.GetTraceID().String())
	if !errors.Is(err, ErrTraceNotFound) {
		return ConnectionTestStepResult{
			OperationDescription: "Tracetest could not get traces back from the data store",
			Error:                err,
			Status:               StatusFailed,
		}
	}

	return ConnectionTestStepResult{
		OperationDescription: "Traces were obtained successfully",
		Error:                nil,
		Status:               StatusPassed,
	}
}

func TracePollingTestStep(ds DataStore) TestStep {
	return &tracePollingTestStep{ds}
}
