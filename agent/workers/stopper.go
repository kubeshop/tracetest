package workers

import (
	"context"
	"errors"
	"fmt"

	"github.com/kubeshop/tracetest/agent/event"
	"github.com/kubeshop/tracetest/agent/proto"
	"go.uber.org/zap"
)

type StopperWorker struct {
	logger         *zap.Logger
	observer       event.Observer
	cancelContexts *cancelCauseFuncMap
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

func NewStopperWorker(opts ...StopperOption) *StopperWorker {
	worker := &StopperWorker{
		logger:   zap.NewNop(),
		observer: event.NewNopObserver(),
	}

	for _, opt := range opts {
		opt(worker)
	}

	return worker
}

func (w *StopperWorker) SetLogger(logger *zap.Logger) {
	w.logger = logger
}

func (w *StopperWorker) Stop(ctx context.Context, stopRequest *proto.StopRequest) error {
	w.logger.Debug("Stop request received", zap.Any("stopRequest", stopRequest))
	w.observer.StartStopRequest(stopRequest)

	cacheKey := key(stopRequest.TestID, stopRequest.RunID)
	cancelFn, found := w.cancelContexts.Get(cacheKey)
	if !found {
		err := fmt.Errorf("cancel func for StopRequest not found")
		w.logger.Error(err.Error(), zap.String("testID", stopRequest.TestID), zap.Int32("runID", stopRequest.RunID))
		w.observer.EndStopRequest(stopRequest, err)
		return err
	}

	cancelFn(errors.New(stopRequest.Type))

	w.observer.EndStopRequest(stopRequest, nil)

	return nil
}
