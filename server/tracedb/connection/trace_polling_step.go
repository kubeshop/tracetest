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

func (s *tracePollingTestStep) TestConnection(ctx context.Context) model.ConnectionTestStep {
	_, err := s.dataStore.GetTraceByID(ctx, s.dataStore.GetTraceID().String())
	if !errors.Is(err, ErrTraceNotFound) {
		return model.ConnectionTestStep{
			Message: "Tracetest could not get traces back from the data store",
			Error:   err,
			Status:  model.StatusFailed,
		}
	}

	return model.ConnectionTestStep{
		Message: "Traces were obtained successfully",
		Error:   nil,
		Status:  model.StatusPassed,
	}
}

func TracePollingTestStep(ds DataStore) TestStep {
	return &tracePollingTestStep{ds}
}
