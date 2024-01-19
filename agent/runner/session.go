package runner

import (
	"context"
	"errors"
	"fmt"

	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/collector"
	"github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/agent/event"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/workers"
	"github.com/kubeshop/tracetest/agent/workers/poller"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var ErrOtlpServerStart = errors.New("OTLP server start error")

type Session struct {
	Token  string
	client *client.Client
}

func (s *Session) Close() {
	s.client.Close()
}

func (s *Session) WaitUntilDisconnected() {
	s.client.WaitUntilDisconnected()
}

// Start the agent session with given configuration
func StartSession(ctx context.Context, cfg config.Config, observer event.Observer, logger *zap.Logger) (*Session, error) {
	observer = event.WrapObserver(observer)

	traceCache := collector.NewTraceCache()
	controlPlaneClient, err := newControlPlaneClient(ctx, cfg, traceCache, observer, logger)
	if err != nil {
		return nil, err
	}

	err = controlPlaneClient.Start(ctx)
	if err != nil {
		return nil, err
	}

	agentCollector, err := StartCollector(ctx, cfg, traceCache, observer, logger)
	if err != nil {
		return nil, err
	}

	controlPlaneClient.OnOTLPConnectionTest(func(ctx context.Context, otr *proto.OTLPConnectionTestRequest) error {
		if otr.ResetCounter {
			agentCollector.ResetStatistics()
			return nil
		}

		statistics := agentCollector.Statistics()
		controlPlaneClient.SendOTLPConnectionResult(ctx, &proto.OTLPConnectionTestResponse{
			RequestID:         otr.RequestID,
			SpanCount:         int64(statistics.SpanCount),
			LastSpanTimestamp: statistics.LastSpanTimestamp.UnixMilli(),
		})
		return nil
	})

	return &Session{
		client: controlPlaneClient,
		Token:  controlPlaneClient.SessionConfiguration().AgentIdentification.Token,
	}, nil
}

func StartCollector(ctx context.Context, config config.Config, traceCache collector.TraceCache, observer event.Observer, logger *zap.Logger) (collector.Collector, error) {
	noopTracer := trace.NewNoopTracerProvider().Tracer("noop")
	collectorConfig := collector.Config{
		HTTPPort: config.OTLPServer.HTTPPort,
		GRPCPort: config.OTLPServer.GRPCPort,
	}

	opts := []collector.CollectorOption{
		collector.WithTraceCache(traceCache),
		collector.WithStartRemoteServer(false),
		collector.WithObserver(observer),
		collector.WithLogger(logger),
	}

	collector, err := collector.Start(
		ctx,
		collectorConfig,
		noopTracer,
		opts...,
	)
	if err != nil {
		return nil, ErrOtlpServerStart
	}

	return collector, nil
}

func newControlPlaneClient(ctx context.Context, config config.Config, traceCache collector.TraceCache, observer event.Observer, logger *zap.Logger) (*client.Client, error) {
	controlPlaneClient, err := client.Connect(ctx, config.ServerURL,
		client.WithAPIKey(config.APIKey),
		client.WithAgentName(config.Name),
		client.WithLogger(logger),
	)
	if err != nil {
		observer.Error(err)
		return nil, err
	}

	cancelFuncs := workers.NewCancelFuncMap()

	stopWorker := workers.NewStopperWorker(
		workers.WithStopperObserver(observer),
		workers.WithStopperCancelFuncList(cancelFuncs),
	)

	triggerWorker := workers.NewTriggerWorker(
		controlPlaneClient,
		workers.WithTraceCache(traceCache),
		workers.WithTriggerObserver(observer),
		workers.WithTriggerCancelFuncList(cancelFuncs),
	)

	pollingWorker := workers.NewPollerWorker(
		controlPlaneClient,
		workers.WithInMemoryDatastore(poller.NewInMemoryDatastore(traceCache)),
		workers.WithObserver(observer),
	)

	dataStoreTestConnectionWorker := workers.NewTestConnectionWorker(controlPlaneClient, observer)

	triggerWorker.SetLogger(logger)
	pollingWorker.SetLogger(logger)
	dataStoreTestConnectionWorker.SetLogger(logger)

	controlPlaneClient.OnDataStoreTestConnectionRequest(dataStoreTestConnectionWorker.Test)
	controlPlaneClient.OnStopRequest(stopWorker.Stop)
	controlPlaneClient.OnTriggerRequest(triggerWorker.Trigger)
	controlPlaneClient.OnPollingRequest(pollingWorker.Poll)
	controlPlaneClient.OnConnectionClosed(func(ctx context.Context, sr *proto.ShutdownRequest) error {
		fmt.Printf("Server terminated the connection with the agent. Reason: %s\n", sr.Reason)
		return controlPlaneClient.Close()
	})

	return controlPlaneClient, nil
}
