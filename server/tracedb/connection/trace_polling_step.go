package connection

import (
	"context"
	"errors"

	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
)

type DataStore interface {
	GetTraceByID(context.Context, string) (model.Trace, error)
}

type tracePollingTestStep struct {
	dataStore DataStore
}

func (s *tracePollingTestStep) TestConnection(ctx context.Context) ConnectionTestStepResult {
	_, err := s.dataStore.GetTraceByID(ctx, id.NewRandGenerator().TraceID().String())
	if !errors.Is(err, ErrTraceNotFound) {
		return ConnectionTestStepResult{
			OperationDescription: "Tracetest could not get traces back from the data store",
			Error:                err,
		}
	}

	return ConnectionTestStepResult{
		OperationDescription: "Traces were obtained successfully",
		Error:                nil,
	}
}

func TracePollingTestStep(ds DataStore) TestStep {
	return &tracePollingTestStep{ds}
}
