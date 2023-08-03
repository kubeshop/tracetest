package otlp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/trace"
	pb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
)

type runGetter interface {
	GetRunByTraceID(context.Context, trace.TraceID) (test.Run, error)
}

type tracePersister interface {
	SaveTrace(context.Context, traces.Trace) error
}

type ingester struct {
	tracePersister tracePersister
	runGetter      runGetter
	eventEmitter   executor.EventEmitter
	dsRepo         *datastore.Repository
}

func NewIngester(tracePersister tracePersister, runRepository runGetter, eventEmitter executor.EventEmitter, dsRepo *datastore.Repository) ingester {
	return ingester{
		tracePersister: tracePersister,
		runGetter:      runRepository,
		eventEmitter:   eventEmitter,
		dsRepo:         dsRepo,
	}
}

func (i ingester) Ingest(ctx context.Context, request *pb.ExportTraceServiceRequest, requestType string) (*pb.ExportTraceServiceResponse, error) {
	ds, err := i.dsRepo.Current(ctx)

	if err != nil || !ds.IsOTLPBasedProvider() {
		fmt.Println("OTLP server is not enabled. Ignoring request")
		return &pb.ExportTraceServiceResponse{}, nil
	}

	if len(request.ResourceSpans) == 0 {
		return &pb.ExportTraceServiceResponse{}, nil
	}

	modelTrace := traces.FromOtelResourceSpans(request.ResourceSpans)
	err = i.tracePersister.SaveTrace(ctx, modelTrace)
	if err != nil {
		return nil, fmt.Errorf("failed to save trace: %w", err)
	}

	err = i.notify(ctx, modelTrace, requestType)
	if err != nil {
		return nil, fmt.Errorf("failed to notify: %w", err)
	}

	return &pb.ExportTraceServiceResponse{}, nil
}

func (i ingester) notify(ctx context.Context, trace traces.Trace, requestType string) error {
	run, err := i.runGetter.GetRunByTraceID(ctx, trace.ID)
	if errors.Is(err, sql.ErrNoRows) {
		// trace is not part of any known test run, so no need to notify
		return nil
	}
	if err != nil {
		// there was an actual error accessing the DB
		return fmt.Errorf("error getting run by traceID: %w", err)
	}

	evt := events.TraceOtlpServerReceivedSpans(
		run.TestID,
		run.ID,
		len(trace.Flat),
		requestType,
	)
	err = i.eventEmitter.Emit(ctx, evt)
	if err != nil {
		// there was an actual error accessing the DB
		return fmt.Errorf("error getting run by traceID: %w", err)
	}

	return nil
}
