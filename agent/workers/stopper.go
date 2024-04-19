package workers

import (
	"context"
	"errors"
	"fmt"

	"github.com/kubeshop/tracetest/agent/telemetry"

	"github.com/kubeshop/tracetest/agent/event"
	"github.com/kubeshop/tracetest/agent/proto"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type StopperWorker struct {
	logger         *zap.Logger
	observer       event.Observer
	cancelContexts *cancelCauseFuncMap
	tracer         trace.Tracer
	meter          metric.Meter
}

type StopperOption func(*StopperWorker)

func WithStopperCancelFuncList(cancelContexts *cancelCauseFuncMap) StopperOption {
	return func(tw *StopperWorker) {
		tw.cancelContexts = cancelContexts
	}
}

func WithStopperObserver(observer event.Observer) StopperOption {
	return func(tw *StopperWorker) {
		tw.observer = observer
	}
}

func WithStopperTracer(tracer trace.Tracer) StopperOption {
	return func(tw *StopperWorker) {
		tw.tracer = tracer
	}
}

func WithStopperMeter(meter metric.Meter) StopperOption {
	return func(tw *StopperWorker) {
		tw.meter = meter
	}
}

func WithStopperLogger(logger *zap.Logger) StopperOption {
	return func(tw *StopperWorker) {
		tw.logger = logger
	}
}

func NewStopperWorker(opts ...StopperOption) *StopperWorker {
	worker := &StopperWorker{
		logger:   zap.NewNop(),
		observer: event.NewNopObserver(),
		tracer:   telemetry.GetNoopTracer(),
		meter:    telemetry.GetNoopMeter(),
	}

	for _, opt := range opts {
		opt(worker)
	}

	return worker
}

func (w *StopperWorker) Stop(ctx context.Context, request *proto.StopRequest) error {
	ctx, span := w.tracer.Start(ctx, "StopRequest Worker operation")
	defer span.End()

	runCounter, _ := w.meter.Int64Counter("tracetest.agent.stopworker.runs")
	runCounter.Add(ctx, 1)

	errorCounter, _ := w.meter.Int64Counter("tracetest.agent.stopworker.errors")

	w.logger.Debug("Stop request received", zap.Any("stopRequest", request))
	w.observer.StartStopRequest(request)

	cacheKey := key(request.TestID, request.RunID)
	cancelFn, found := w.cancelContexts.Get(cacheKey)
	if !found {
		err := fmt.Errorf("cancel func for StopRequest not found")
		w.logger.Error(err.Error(), zap.String("testID", request.TestID), zap.Int32("runID", request.RunID))
		w.observer.EndStopRequest(request, err)
		span.RecordError(err)

		errorCounter.Add(ctx, 1)

		return err
	}

	cancelFn(errors.New(request.Type))

	w.observer.EndStopRequest(request, nil)

	return nil
}
