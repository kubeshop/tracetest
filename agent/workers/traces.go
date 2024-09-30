package workers

import (
	"context"
	"errors"
	"fmt"

	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/telemetry"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	pb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"go.uber.org/zap"
)

var (
	ErrNotSupportedDataStore = fmt.Errorf("datastore not supported, only OTLP based datastores are supported")
)

type TracesWorker struct {
	client *client.Client
	logger *zap.Logger
	tracer trace.Tracer
	meter  metric.Meter
}

type TracesOption func(*TracesWorker)

func WithTracesLogger(logger *zap.Logger) TracesOption {
	return func(w *TracesWorker) {
		w.logger = logger
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

func (w *TracesWorker) Export(ctx context.Context, r *pb.ExportTraceServiceRequest) error {
	return w.client.SendTraces(ctx, &proto.ExportRequest{
		Request: r,
	})
}
