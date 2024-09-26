package workers

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/collector"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/telemetry"
	"github.com/kubeshop/tracetest/agent/tracedb"
	"github.com/kubeshop/tracetest/agent/tracedb/connection"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var (
	ErrNotSupportedDataStore = fmt.Errorf("datastore not supported, only OTLP based datastores are supported")
)

type TracesWorker struct {
	client            *client.Client
	inmemoryDatastore tracedb.TraceDB
	logger            *zap.Logger
	tracer            trace.Tracer
	meter             metric.Meter
	traceCache        collector.TraceCache
}

type TracesOption func(*TracesWorker)

func WithTracesLogger(logger *zap.Logger) TracesOption {
	return func(w *TracesWorker) {
		w.logger = logger
	}
}

func WithTracesInMemoryDatastore(datastore tracedb.TraceDB) TracesOption {
	return func(pw *TracesWorker) {
		pw.inmemoryDatastore = datastore
	}
}

func WithTracesTraceCache(cache collector.TraceCache) TracesOption {
	return func(tw *TracesWorker) {
		tw.traceCache = cache
	}
}

func WithTracesTracer(tracer trace.Tracer) TracesOption {
	return func(w *TracesWorker) {
		w.tracer = tracer
	}
}

func WithTracesMeter(meter metric.Meter) TracesOption {
	return func(w *TracesWorker) {
		w.meter = meter
	}
}

func NewTracesWorker(client *client.Client, opts ...TracesOption) *TracesWorker {
	worker := &TracesWorker{
		client: client,
		tracer: telemetry.GetNoopTracer(),
		logger: zap.NewNop(),
		meter:  telemetry.GetNoopMeter(),
	}

	for _, opt := range opts {
		opt(worker)
	}

	return worker
}

var (
	ErrInvalidRequest = errors.New("invalid request")
)

func (w *TracesWorker) Handler(ctx context.Context, request *proto.TraceModeRequest) error {
	if request.ListTracesRequest != nil {
		return w.listTraces(ctx, request)
	}

	if request.GetTraceRequest != nil {
		return w.getTrace(ctx, request)
	}

	err := ErrInvalidRequest
	errorResponse := &proto.TraceModeResponse{
		RequestID:           request.RequestID,
		AgentIdentification: w.client.SessionConfiguration().AgentIdentification,
		Error: &proto.Error{
			Message: err.Error(),
		},
	}

	err = w.client.SendTraceModeResponse(ctx, errorResponse)
	if err != nil {
		w.logger.Error("Error sending polling error", zap.Error(err))
		return err
	}

	return err
}

func (w *TracesWorker) listTraces(ctx context.Context, request *proto.TraceModeRequest) error {
	ctx, span := w.tracer.Start(ctx, "List Traces Worker operation")
	defer span.End()

	runCounter, _ := w.meter.Int64Counter("tracetest.agent.traces.list_traces")
	runCounter.Add(ctx, 1)

	errorCounter, _ := w.meter.Int64Counter("tracetest.agent.traces.list_traces")

	w.logger.Debug("List Traces request received", zap.Any("list_traces", request))

	var err error
	ds, err := w.getTraceDB(request.ListTracesRequest.Datastore)
	if err != nil {
		errorResponse := &proto.TraceModeResponse{
			RequestID:           request.RequestID,
			AgentIdentification: w.client.SessionConfiguration().AgentIdentification,
			Error: &proto.Error{
				Message: err.Error(),
			},
		}

		sendErr := w.client.SendTraceModeResponse(ctx, errorResponse)
		if sendErr != nil {
			w.logger.Error("Error sending polling error", zap.Error(sendErr))

			formattedErr := fmt.Errorf("could not report polling error back to the server: %w. Original error: %s", sendErr, err.Error())
			span.RecordError(formattedErr)
			errorCounter.Add(ctx, 1)

			return formattedErr
		}

		span.RecordError(err)
		errorCounter.Add(ctx, 1)
		return err
	}

	defer ds.Close()

	take := 20
	if request.ListTracesRequest.Take != 0 {
		take = int(request.ListTracesRequest.Take)
	}

	skip := 0
	if request.ListTracesRequest.Skip != 0 {
		skip = int(request.ListTracesRequest.Skip)
	}

	list, err := ds.List(ctx, take, skip)
	if err != nil {
		errorResponse := &proto.TraceModeResponse{
			RequestID:           request.RequestID,
			AgentIdentification: w.client.SessionConfiguration().AgentIdentification,
			Error: &proto.Error{
				Message: err.Error(),
			},
		}

		sendErr := w.client.SendTraceModeResponse(ctx, errorResponse)
		if sendErr != nil {
			w.logger.Error("Error sending polling error", zap.Error(sendErr))

			formattedErr := fmt.Errorf("could not report polling error back to the server: %w. Original error: %s", sendErr, err.Error())
			span.RecordError(formattedErr)
			errorCounter.Add(ctx, 1)

			return formattedErr
		}

		span.RecordError(err)
		errorCounter.Add(ctx, 1)
	}

	traces := convertTraceListInToProtoListTraces(list)

	response := &proto.TraceModeResponse{
		RequestID:           request.RequestID,
		AgentIdentification: w.client.SessionConfiguration().AgentIdentification,
		ListTracesResponse: &proto.ListTracesResponse{
			Traces: traces,
		},
	}

	w.logger.Debug("Converted traces", zap.Any("trace", traces), zap.Any("listTracesResponse", spew.Sdump(response)))
	err = w.client.SendTraceModeResponse(ctx, response)

	return err
}

func (w *TracesWorker) getTrace(ctx context.Context, request *proto.TraceModeRequest) error {
	ctx, span := w.tracer.Start(ctx, "Get Trace Worker operation")
	defer span.End()

	runCounter, _ := w.meter.Int64Counter("tracetest.agent.traces.get_trace")
	runCounter.Add(ctx, 1)
	errorCounter, _ := w.meter.Int64Counter("tracetest.agent.traces.get_trace")

	w.logger.Debug("Get Trace request received", zap.Any("get_trace", request))

	ds, err := w.getTraceDB(request.GetTraceRequest.Datastore)
	if err != nil {
		errorResponse := &proto.TraceModeResponse{
			RequestID:           request.RequestID,
			AgentIdentification: w.client.SessionConfiguration().AgentIdentification,
			Error: &proto.Error{
				Message: err.Error(),
			},
		}

		sendErr := w.client.SendTraceModeResponse(ctx, errorResponse)
		if sendErr != nil {
			w.logger.Error("Error sending polling error", zap.Error(sendErr))

			formattedErr := fmt.Errorf("could not report polling error back to the server: %w. Original error: %s", sendErr, err.Error())
			span.RecordError(formattedErr)
			errorCounter.Add(ctx, 1)

			return formattedErr
		}

		span.RecordError(err)
		errorCounter.Add(ctx, 1)
		return err
	}

	defer ds.Close()

	getTraceResponse := &proto.TraceModeResponse{
		RequestID:        request.RequestID,
		GetTraceResponse: &proto.GetTraceResponse{},
	}

	trace, err := ds.GetTraceByID(ctx, request.GetTraceRequest.TraceID)
	if err != nil {
		w.logger.Debug("Trace not found")
		getTraceResponse.Error = &proto.Error{Message: connection.ErrTraceNotFound.Error()}
	} else {
		w.logger.Debug("Trace found")
		getTraceResponse.GetTraceResponse.Spans = convertTraceInToProtoSpans(trace)
		w.logger.Debug("Converted trace", zap.Any("trace", trace), zap.Any("pollingResponse", spew.Sdump(getTraceResponse)))
	}

	err = w.client.SendTraceModeResponse(ctx, getTraceResponse)
	if err != nil {
		w.logger.Error("Cannot send trace to server", zap.Error(err))
		log.Printf("cannot send trace to server: %s", err.Error())
		return err
	}

	return nil
}

func (w *TracesWorker) getTraceDB(protoDataStore *proto.DataStore) (tracedb.TraceDB, error) {
	datastoreConfig, err := convertProtoToDataStore(protoDataStore)
	if err != nil {
		w.logger.Error("Invalid datastore", zap.Error(err))
		return nil, err
	}
	w.logger.Debug("Converted datastore", zap.Any("datastore", datastoreConfig), zap.Any("originalDatastore", protoDataStore))

	if datastoreConfig == nil {
		w.logger.Error("Invalid datastore: nil")
		return nil, fmt.Errorf("invalid datastore: nil")
	}

	// Only OTLP based datastores are supported for V1
	if !datastoreConfig.IsOTLPBasedProvider() {
		return nil, ErrNotSupportedDataStore
	}

	dsFactory := tracedb.Factory(nil)
	ds, err := dsFactory(*datastoreConfig)
	if err != nil {
		w.logger.Error("Invalid datastore", zap.Error(err))
		log.Printf("Invalid datastore: %s", err.Error())
		return nil, err
	}

	w.logger.Debug("Created datastore", zap.Any("datastore", ds), zap.Bool("isOTLPBasedProvider", datastoreConfig.IsOTLPBasedProvider()))

	if datastoreConfig.IsOTLPBasedProvider() && w.inmemoryDatastore != nil {
		w.logger.Debug("Using in-memory datastore")
		ds = w.inmemoryDatastore
	}

	return ds, nil
}

func convertTraceListInToProtoListTraces(traces []traces.TraceMetadata) []*proto.TraceMetadata {
	protTraces := make([]*proto.TraceMetadata, 0, len(traces))
	for _, trace := range traces {
		protoTrace := proto.TraceMetadata{
			TraceID:           trace.TraceID,
			RootServiceName:   trace.RootServiceName,
			RootTraceName:     trace.RootTraceName,
			StartTimeUnixNano: trace.StartTimeUnixNano,
			DurationMs:        int32(trace.DurationMs),
			SpanCount:         int32(trace.SpanCount),
		}

		protTraces = append(protTraces, &protoTrace)
	}

	return protTraces
}
