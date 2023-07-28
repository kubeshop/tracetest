package otlp

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/trace"
	pb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

type runGetter interface {
	GetRunByTraceID(context.Context, trace.TraceID) (test.Run, error)
}

type tracePersister interface {
	UpdateTraceSpans(context.Context, *traces.Trace) error
}

type Ingester struct {
	log            func(string, ...interface{})
	tracePersister tracePersister
	runGetter      runGetter
	eventEmitter   executor.EventEmitter
	dsRepo         *datastore.Repository
}

func NewIngester(tracePersister tracePersister, runRepository runGetter, eventEmitter executor.EventEmitter, dsRepo *datastore.Repository) Ingester {
	return Ingester{
		log: func(format string, args ...interface{}) {
			log.Printf("[OTLP] "+format, args...)
		},
		tracePersister: tracePersister,
		runGetter:      runRepository,
		eventEmitter:   eventEmitter,
		dsRepo:         dsRepo,
	}
}

type RequestType string

var (
	RequestTypeHTTP RequestType = "HTTP"
	RequestTypeGRPC RequestType = "gRPC"
)

func (i Ingester) Ingest(ctx context.Context, request *pb.ExportTraceServiceRequest, requestType RequestType) (*pb.ExportTraceServiceResponse, error) {
	ds, err := i.dsRepo.Current(ctx)

	if err != nil || !ds.IsOTLPBasedProvider() {
		i.log("OTLP server is not enabled. Ignoring request")
		return &pb.ExportTraceServiceResponse{}, nil
	}

	if len(request.ResourceSpans) == 0 {
		i.log("no spans to ingest")
		return &pb.ExportTraceServiceResponse{}, nil
	}

	receivedTraces := i.traces(request.ResourceSpans)
	i.log("received %d traces", len(receivedTraces))

	// each request can have different traces so we need to go over each individual trace
	for ix, modelTrace := range receivedTraces {
		i.log("processing trace %d/%d traceID %s", ix+1, len(receivedTraces), modelTrace.ID.String())
		// if at some point we want to save all traces, not just the ones related to a running test
		// we can just remove this check
		run, err := i.getOngoinTestRunForTrace(ctx, modelTrace)
		if errors.Is(err, errNoTestRun) {
			i.log("trace %s is not part of any ongoing test run", modelTrace.ID.String())
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get test run: %w", err)
		}

		err = i.tracePersister.UpdateTraceSpans(ctx, &modelTrace)
		if err != nil {
			return nil, fmt.Errorf("failed to save trace: %w", err)
		}

		err = i.notify(ctx, run, modelTrace, requestType)
		if err != nil {
			return nil, fmt.Errorf("failed to notify: %w", err)
		}
	}

	return &pb.ExportTraceServiceResponse{
		PartialSuccess: &pb.ExportTracePartialSuccess{
			RejectedSpans: 0,
		},
	}, nil
}

func (i Ingester) traces(input []*v1.ResourceSpans) []traces.Trace {
	spansByTrace := map[string][]*v1.Span{}

	for _, rs := range input {
		for _, il := range rs.ScopeSpans {
			for _, span := range il.Spans {
				traceID := trace.TraceID(span.TraceId).String()
				i.log("adding span %s to trace %s", hex.EncodeToString(span.SpanId), traceID)
				spansByTrace[traceID] = append(spansByTrace[traceID], span)
			}
		}
	}

	i.log("sorted %d traces", len(spansByTrace))

	modelTraces := make([]traces.Trace, 0, len(spansByTrace))
	for traceID, spans := range spansByTrace {
		i.log("creating trace %s with %d spans", traceID, len(spans))
		modelTraces = append(
			modelTraces,
			traces.FromSpanList(spans),
		)
	}

	return modelTraces
}

var errNoTestRun = errors.New("no test run")

func (i Ingester) getOngoinTestRunForTrace(ctx context.Context, trace traces.Trace) (test.Run, error) {
	run, err := i.runGetter.GetRunByTraceID(ctx, trace.ID)
	if errors.Is(err, sql.ErrNoRows) {
		// trace is not part of any known test run, no need to notify
		return test.Run{}, errNoTestRun
	}
	if err != nil {
		// there was an actual error accessing the DB
		return test.Run{}, fmt.Errorf("error getting run by traceID: %w", err)
	}

	if run.State != test.RunStateAwaitingTrace {
		return test.Run{}, errNoTestRun
	}

	return run, nil
}

func (i Ingester) notify(ctx context.Context, run test.Run, trace traces.Trace, requestType RequestType) error {
	evt := events.TraceOtlpServerReceivedSpans(
		run.TestID,
		run.ID,
		len(trace.Flat),
		string(requestType),
	)
	err := i.eventEmitter.Emit(ctx, evt)
	if err != nil {
		// there was an actual error accessing the DB
		return fmt.Errorf("error getting run by traceID: %w", err)
	}

	return nil
}
