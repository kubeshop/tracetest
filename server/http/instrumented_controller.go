package http

import (
	"context"

	"github.com/kubeshop/tracetest/server/openapi"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type instrumentedServicer struct {
	tracer   trace.Tracer
	servicer openapi.ApiApiServicer
}

func NewInstrumentedServicer(tracer trace.Tracer, servicer openapi.ApiApiServicer) openapi.ApiApiServicer {
	return instrumentedServicer{
		tracer:   tracer,
		servicer: servicer,
	}
}

type instrumentedFunction func(context.Context) (openapi.ImplResponse, error)

func (is instrumentedServicer) instrumentFunction(ctx context.Context, name string, f instrumentedFunction) (openapi.ImplResponse, error) {
	instrumentedCtx, span := is.tracer.Start(ctx, name)
	defer span.End()
	response, err := f(instrumentedCtx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}

	return response, err
}

func (is instrumentedServicer) CreateTest(ctx context.Context, test openapi.Test) (openapi.ImplResponse, error) {
	return is.instrumentFunction(ctx, "CreateTest", func(ctx context.Context) (openapi.ImplResponse, error) {
		return is.servicer.CreateTest(ctx, test)
	})
}

func (is instrumentedServicer) DeleteTest(ctx context.Context, id string) (openapi.ImplResponse, error) {
	return is.instrumentFunction(ctx, "DeleteTest", func(ctx context.Context) (openapi.ImplResponse, error) {
		return is.servicer.DeleteTest(ctx, id)
	})
}

func (is instrumentedServicer) DeleteTestRun(ctx context.Context, testID string, runID string) (openapi.ImplResponse, error) {
	return is.instrumentFunction(ctx, "DeleteTestRun", func(ctx context.Context) (openapi.ImplResponse, error) {
		return is.servicer.DeleteTestRun(ctx, testID, runID)
	})
}

func (is instrumentedServicer) DryRunAssertion(ctx context.Context, testID string, runID string, definition openapi.TestDefinition) (openapi.ImplResponse, error) {
	return is.instrumentFunction(ctx, "DryRunAssertion", func(ctx context.Context) (openapi.ImplResponse, error) {
		return is.servicer.DryRunAssertion(ctx, testID, runID, definition)
	})
}

func (is instrumentedServicer) GetTest(ctx context.Context, id string) (openapi.ImplResponse, error) {
	return is.instrumentFunction(ctx, "GetTest", func(ctx context.Context) (openapi.ImplResponse, error) {
		return is.servicer.GetTest(ctx, id)
	})
}

func (is instrumentedServicer) GetTestDefinition(ctx context.Context, id string) (openapi.ImplResponse, error) {
	return is.instrumentFunction(ctx, "GetTestDefinition", func(ctx context.Context) (openapi.ImplResponse, error) {
		return is.servicer.GetTestDefinition(ctx, id)
	})
}

func (is instrumentedServicer) GetTestResultSelectedSpans(ctx context.Context, testID string, runID string, selector string) (openapi.ImplResponse, error) {
	return is.instrumentFunction(ctx, "GetTestResultSelectedSpans", func(ctx context.Context) (openapi.ImplResponse, error) {
		return is.servicer.GetTestResultSelectedSpans(ctx, testID, runID, selector)
	})
}

func (is instrumentedServicer) GetTestRun(ctx context.Context, testID string, runID string) (openapi.ImplResponse, error) {
	return is.instrumentFunction(ctx, "GetTestRun", func(ctx context.Context) (openapi.ImplResponse, error) {
		return is.servicer.GetTestRun(ctx, testID, runID)
	})
}

func (is instrumentedServicer) GetTestRuns(ctx context.Context, testID string, take int32, skip int32) (openapi.ImplResponse, error) {
	return is.instrumentFunction(ctx, "GetTestRuns", func(ctx context.Context) (openapi.ImplResponse, error) {
		return is.servicer.GetTestRuns(ctx, testID, take, skip)
	})
}

func (is instrumentedServicer) GetTests(ctx context.Context, take int32, skip int32) (openapi.ImplResponse, error) {
	return is.instrumentFunction(ctx, "GetTests", func(ctx context.Context) (openapi.ImplResponse, error) {
		return is.servicer.GetTests(ctx, take, skip)
	})
}

func (is instrumentedServicer) RerunTestRun(ctx context.Context, testID string, runID string) (openapi.ImplResponse, error) {
	return is.instrumentFunction(ctx, "RerunTestRun", func(ctx context.Context) (openapi.ImplResponse, error) {
		return is.servicer.RerunTestRun(ctx, testID, runID)
	})
}

func (is instrumentedServicer) RunTest(ctx context.Context, testID string) (openapi.ImplResponse, error) {
	return is.instrumentFunction(ctx, "RunTest", func(ctx context.Context) (openapi.ImplResponse, error) {
		return is.servicer.RunTest(ctx, testID)
	})
}

func (is instrumentedServicer) SetTestDefinition(ctx context.Context, testID string, definition openapi.TestDefinition) (openapi.ImplResponse, error) {
	return is.instrumentFunction(ctx, "SetTestDefinition", func(ctx context.Context) (openapi.ImplResponse, error) {
		return is.servicer.SetTestDefinition(ctx, testID, definition)
	})
}

func (is instrumentedServicer) UpdateTest(ctx context.Context, testID string, test openapi.Test) (openapi.ImplResponse, error) {
	return is.instrumentFunction(ctx, "UpdateTest", func(ctx context.Context) (openapi.ImplResponse, error) {
		return is.servicer.UpdateTest(ctx, testID, test)
	})
}

func (is instrumentedServicer) GetRunResultJUnit(ctx context.Context, testID string, runID string) (openapi.ImplResponse, error) {
	return is.instrumentFunction(ctx, "GetRunResultJUnit", func(ctx context.Context) (openapi.ImplResponse, error) {
		return is.servicer.GetRunResultJUnit(ctx, testID, runID)
	})
}
