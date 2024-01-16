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
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/testconnection"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	pb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

type RequestType string

var (
	RequestTypeHTTP RequestType = "HTTP"
	RequestTypeGRPC RequestType = "gRPC"
)

type Ingester interface {
	Ingest(ctx context.Context, request *pb.ExportTraceServiceRequest, requestType RequestType) (*pb.ExportTraceServiceResponse, error)
}

type runGetter interface {
	GetRunByTraceID(context.Context, trace.TraceID) (test.Run, error)
}

type tracePersister interface {
	UpdateTraceSpans(context.Context, *traces.Trace) error
}

func NewIngester(
	tracePersister tracePersister,
	runRepository runGetter,
	eventEmitter executor.EventEmitter,
	dsRepo *datastore.Repository,
	subManager subscription.Manager,
	tracer trace.Tracer,
) ingester {
	ingester := ingester{
		log: func(format string, args ...interface{}) {
			log.Printf("[OTLP] "+format, args...)
		},
		tracePersister: tracePersister,
		runGetter:      runRepository,
		eventEmitter:   eventEmitter,
		dsRepo:         dsRepo,
		subManager:     subManager,
		tracer:         tracer,
	}

	ingester.startTesterListener()

	return ingester
}

type ingester struct {
	log            func(string, ...interface{})
	tracePersister tracePersister
	runGetter      runGetter
	eventEmitter   executor.EventEmitter
	dsRepo         *datastore.Repository
	subManager     subscription.Manager
	tracer         trace.Tracer

	captureTracesForConnectionTest bool
}

func (i *ingester) startTesterListener() {
	i.subManager.Subscribe("start_otlp_connection_test_", subscription.NewSubscriberFunction(func(m subscription.Message) error {
		i.captureTracesForConnectionTest = true
		return nil
	}))

	i.subManager.Subscribe("end_otlp_connection_test_", subscription.NewSubscriberFunction(func(m subscription.Message) error {
		i.captureTracesForConnectionTest = false
		return nil
	}))
}

func (i ingester) Ingest(ctx context.Context, request *pb.ExportTraceServiceRequest, requestType RequestType) (*pb.ExportTraceServiceResponse, error) {
	ds, err := i.dsRepo.Current(ctx)

	if err != nil || !ds.IsOTLPBasedProvider() {
		i.log("OTLP server is not enabled. Ignoring request")
		return &pb.ExportTraceServiceResponse{}, nil
	}

	if len(request.ResourceSpans) == 0 {
		i.log("no spans to ingest")
		return &pb.ExportTraceServiceResponse{}, nil
	}

	ctx, span := i.tracer.Start(ctx, "convert received spans")
	defer span.End()
	receivedTraces := i.traces(request.ResourceSpans)
	i.log("received %d traces", len(receivedTraces))

	if i.captureTracesForConnectionTest {
		i.subManager.Publish("otlp_connection_test_spans_incoming", testconnection.OTLPConnectionTestResponse{
			NumberTraces: len(receivedTraces),
		})
	}

	// each request can have different traces so we need to go over each individual trace
	for ix, modelTrace := range receivedTraces {
		i.log("processing trace %d/%d traceID %s", ix+1, len(receivedTraces), modelTrace.ID.String())
		err = i.processTrace(ctx, modelTrace, ix, requestType)
		if err != nil {
			span.RecordError(err, trace.WithAttributes(attribute.String("tracetest.ingestor.trace_id", modelTrace.ID.String())))
			return nil, err
		}
	}

	return &pb.ExportTraceServiceResponse{
		PartialSuccess: &pb.ExportTracePartialSuccess{
			RejectedSpans: 0,
		},
	}, nil
}

func (i ingester) processTrace(ctx context.Context, modelTrace traces.Trace, ix int, requestType RequestType) error {
	ctx, span := i.tracer.Start(ctx, "process otlp trace")
	span.SetAttributes(
		attribute.String("tracetest.ingestor.trace_id", modelTrace.ID.String()),
	)
	defer span.End()

	// if at some point we want to save all traces, not just the ones related to a running test
	// we can just remove this check
	run, err := i.getOngoinTestRunForTrace(ctx, modelTrace)
	if errors.Is(err, errNoTestRun) {
		i.log("trace %s is not part of any ongoing test run", modelTrace.ID.String())
		span.RecordError(err)
		return nil
	}
	if err != nil {
		i.log("[%s] failed to get test run: %s", modelTrace.ID.String(), err.Error())
		span.RecordError(err)
		return fmt.Errorf("failed to get test run: %w", err)
	}

	err = i.tracePersister.UpdateTraceSpans(ctx, &modelTrace)
	if err != nil {
		i.log("[%s] failed to save trace: %s", modelTrace.ID.String(), err.Error())
		span.RecordError(err)
		return fmt.Errorf("failed to save trace: %w", err)
	}

	err = i.notify(ctx, run, modelTrace, requestType)
	if err != nil {
		i.log("[%s] failed to notify: %s", modelTrace.ID.String(), err.Error())
		span.RecordError(err)
		return fmt.Errorf("failed to notify: %w", err)
	}

	return nil
}

func (i ingester) traces(input []*v1.ResourceSpans) []traces.Trace {
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

func (i ingester) getOngoinTestRunForTrace(ctx context.Context, trace traces.Trace) (test.Run, error) {
	run, err := i.runGetter.GetRunByTraceID(ctx, trace.ID)
	if errors.Is(err, sql.ErrNoRows) {
		// trace is not part of any known test run, no need to notify
		return test.Run{}, errNoTestRun
	}
	if err != nil {
		// there was an actual error accessing the DB
		return test.Run{}, fmt.Errorf("error getting run by traceID: %w", err)
	}

	if run.State != test.RunStateAwaitingTrace && run.State != test.RunStateExecuting {
		return test.Run{}, fmt.Errorf("test run is not awaiting trace nor executing. Actual state: %s", run.State)
	}

	return run, nil
}

func (i ingester) notify(ctx context.Context, run test.Run, trace traces.Trace, requestType RequestType) error {
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
