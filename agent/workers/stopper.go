package workers

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/agent/event"
	"github.com/kubeshop/tracetest/agent/proto"
	"go.uber.org/zap"
)

type StopperWorker struct {
	logger         *zap.Logger
	observer       event.Observer
	cancelContexts *cancelChannelMap
}

type StopperOption func(*StopperWorker)

func WithStopperCancelContextsList(cancelContexts *cancelChannelMap) StopperOption {
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
	cancelChan, found := w.cancelContexts.Get(cacheKey)
	if !found {
		err := fmt.Errorf("cancel func for StopRequest not found")
		w.logger.Error(err.Error(), zap.String("testID", stopRequest.TestID), zap.Int32("runID", stopRequest.RunID))
		w.observer.EndStopRequest(stopRequest, err)
		return err
	}

	cancelChan <- true

	w.observer.EndStopRequest(stopRequest, nil)

	return nil
}
