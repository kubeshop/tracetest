package runner

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/collector"
	"github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/agent/event"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/telemetry"
	"github.com/kubeshop/tracetest/agent/workers"
	"github.com/kubeshop/tracetest/agent/workers/poller"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var ErrOtlpServerStart = errors.New("OTLP server start error")

type Session struct {
	Token  string
	client *client.Client
}

func (s *Session) Close() {
	if s != nil && s.client != nil {
		s.client.Close()
	}
}

func (s *Session) WaitUntilDisconnected() {
	s.client.WaitUntilDisconnected()
}

// Start the agent session with given configuration
func StartSession(ctx context.Context, cfg config.Config, observer event.Observer, logger *zap.Logger) (*Session, error) {
	logger.Debug("Starting agent session")
	observer = event.WrapObserver(observer)

	tracer, err := telemetry.GetTracer(ctx, cfg.CollectorEndpoint, cfg.Name)
	if err != nil {
		logger.Error("Failed to create tracer", zap.Error(err))
		observer.Error(err)
		return nil, err
	}

	meter, err := telemetry.GetMeter(ctx, cfg.CollectorEndpoint, cfg.Name)
	if err != nil {
		observer.Error(err)
		return nil, err
	}

	traceCache := collector.NewTraceCache()
	controlPlaneClient, err := newControlPlaneClient(ctx, cfg, traceCache, observer, logger, tracer, meter)
	if err != nil {
		return nil, err
	}

	err = controlPlaneClient.Start(ctx)
	if err != nil {
		return nil, err
	}

	// Trace Mode
	traceModeWorker := workers.NewTracesWorker(
		controlPlaneClient,
		workers.WithTracesLogger(logger),
		workers.WithTracesTracer(tracer),
		workers.WithTracesMeter(meter),
	)

	agentCollector, err := StartCollector(ctx, cfg, traceCache, traceModeWorker, observer, logger, tracer)
	if err != nil {
		return nil, err
	}

	controlPlaneClient.OnOTLPConnectionTest(func(ctx context.Context, otr *proto.OTLPConnectionTestRequest) error {
		logger.Debug("Received OTLP connection test request", zap.Bool("resetCounter", otr.ResetCounter))
		if otr.ResetCounter {
			logger.Debug("Resetting statistics")
			agentCollector.ResetStatistics()
			return nil
		}

		statistics := agentCollector.Statistics()
		logger.Debug("Sending OTLP connection test response",
			zap.Int64("spanCount", statistics.SpanCount),
			zap.Int64("lastSpanTimestamp", time.Since(statistics.LastSpanTimestamp).Milliseconds()),
		)

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

func StartCollector(ctx context.Context, config config.Config, traceCache collector.TraceCache, traceModeWorker collector.TraceModeForwarder, observer event.Observer, logger *zap.Logger, tracer trace.Tracer) (collector.Collector, error) {

	collectorConfig := collector.Config{
		HTTPPort: config.OTLPServer.HTTPPort,
		GRPCPort: config.OTLPServer.GRPCPort,
	}

	opts := []collector.CollectorOption{
		collector.WithTraceCache(traceCache),
		collector.WithStartRemoteServer(false),
		collector.WithObserver(observer),
		collector.WithLogger(logger),
		collector.WithTraceMode(config.TraceMode),
		collector.WithTraceModeForwarder(traceModeWorker),
	}

	collector, err := collector.Start(
		ctx,
		collectorConfig,
		tracer,
		opts...,
	)
	if err != nil {
		return nil, ErrOtlpServerStart
	}

	return collector, nil
}

func newControlPlaneClient(ctx context.Context, config config.Config, traceCache collector.TraceCache, observer event.Observer, logger *zap.Logger, tracer trace.Tracer, meter metric.Meter) (*client.Client, error) {
	opts := []client.Option{
		client.WithAPIKey(config.APIKey),
		client.WithAgentName(config.Name),
		client.WithEnvironmentID(config.EnvironmentID),
		client.WithLogger(logger),
	}
	if config.Insecure {
		opts = append(opts, client.WithInsecure())
	}

	if config.SkipVerify {
		opts = append(opts, client.WithSkipVerify())
	}
	controlPlaneClient, err := client.Connect(ctx, config.ServerURL,
		opts...,
	)
	if err != nil {
		observer.Error(err)
		return nil, err
	}

	processStopper := workers.NewProcessStopper()

	stopWorker := workers.NewStopperWorker(
		workers.WithStopperObserver(observer),
		workers.WithStopperCancelFuncList(processStopper.CancelMap()),
		workers.WithStopperLogger(logger),
		workers.WithStopperTracer(tracer),
		workers.WithStopperMeter(meter),
	)

	triggerWorker := workers.NewTriggerWorker(
		controlPlaneClient,
		workers.WithTraceCache(traceCache),
		workers.WithTriggerObserver(observer),
		workers.WithTriggerStoppableProcessRunner(processStopper.RunStoppableProcess),
		workers.WithTriggerLogger(logger),
		workers.WithTriggerTracer(tracer),
		workers.WithTriggerMeter(meter),
	)

	inMemoryDs := poller.NewInMemoryDatastore(traceCache)

	pollingWorker := workers.NewPollerWorker(
		controlPlaneClient,
		workers.WithInMemoryDatastore(inMemoryDs),
		workers.WithPollerTraceCache(traceCache),
		workers.WithPollerObserver(observer),
		workers.WithPollerStoppableProcessRunner(processStopper.RunStoppableProcess),
		workers.WithPollerLogger(logger),
		workers.WithPollerTracer(tracer),
		workers.WithPollerMeter(meter),
	)

	dataStoreTestConnectionWorker := workers.NewTestConnectionWorker(
		controlPlaneClient,
		workers.WithTestConnectionLogger(logger),
		workers.WithTestConnectionObserver(observer),
		workers.WithTestConnectionTracer(tracer),
		workers.WithTestConnectionMeter(meter),
	)

	graphqlIntrospectionWorker := workers.NewGraphqlIntrospectWorker(
		controlPlaneClient,
		workers.WithGraphqlIntrospectLogger(logger),
		workers.WithGraphqlIntrospectTracer(tracer),
		workers.WithGraphqlIntrospectMeter(meter),
	)

	controlPlaneClient.OnDataStoreTestConnectionRequest(dataStoreTestConnectionWorker.Test)
	controlPlaneClient.OnStopRequest(stopWorker.Stop)
	controlPlaneClient.OnTriggerRequest(triggerWorker.Trigger)
	controlPlaneClient.OnPollingRequest(pollingWorker.Poll)
	controlPlaneClient.OnGraphqlIntrospectionRequest(graphqlIntrospectionWorker.Introspect)
	controlPlaneClient.OnConnectionClosed(func(ctx context.Context, sr *proto.ShutdownRequest) error {
		logger.Info(fmt.Sprintf("Server terminated the connection with the agent. Reason: %s\n", sr.Reason))
		return controlPlaneClient.Close()
	})

	return controlPlaneClient, nil
}
