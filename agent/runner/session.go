package runner

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/collector"
	"github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/workers"
	"github.com/kubeshop/tracetest/agent/workers/poller"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var ErrOtlpServerStart = errors.New("OTLP server start error")

var logger *zap.Logger

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
func StartSession(ctx context.Context, cfg config.Config) (*Session, error) {
	traceCache := collector.NewTraceCache()
	controlPlaneClient, err := newControlPlaneClient(ctx, cfg, traceCache)
	if err != nil {
		return nil, err
	}

	err = controlPlaneClient.Start(ctx)
	if err != nil {
		return nil, err
	}

	err = StartCollector(ctx, cfg, traceCache)
	if err != nil {
		return nil, err
	}

	return &Session{
		client: controlPlaneClient,
		Token:  controlPlaneClient.SessionConfiguration().AgentIdentification.Token,
	}, nil
}

func StartCollector(ctx context.Context, config config.Config, traceCache collector.TraceCache) error {
	noopTracer := trace.NewNoopTracerProvider().Tracer("noop")
	collectorConfig := collector.Config{
		HTTPPort: config.OTLPServer.HTTPPort,
		GRPCPort: config.OTLPServer.GRPCPort,
	}

	opts := []collector.CollectorOption{
		collector.WithTraceCache(traceCache),
		collector.WithStartRemoteServer(false),
	}

	if enableLogging() {
		opts = append(opts, collector.WithLogger(logger))
	}

	_, err := collector.Start(
		ctx,
		collectorConfig,
		noopTracer,
		opts...,
	)
	if err != nil {
		return ErrOtlpServerStart
	}

	return nil
}

func enableLogging() bool {
	return os.Getenv("TRACETEST_DEV") == "true"
}

func newControlPlaneClient(ctx context.Context, config config.Config, traceCache collector.TraceCache) (*client.Client, error) {
	if enableLogging() {
		var err error
		logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, fmt.Errorf("could not create logger: %w", err)
		}
	}

	controlPlaneClient, err := client.Connect(ctx, config.ServerURL,
		client.WithAPIKey(config.APIKey),
		client.WithAgentName(config.Name),
		client.WithLogger(logger),
	)
	if err != nil {
		return nil, err
	}

	triggerWorker := workers.NewTriggerWorker(controlPlaneClient, workers.WithTraceCache(traceCache))
	pollingWorker := workers.NewPollerWorker(controlPlaneClient, workers.WithInMemoryDatastore(
		poller.NewInMemoryDatastore(traceCache),
	))
	dataStoreTestConnectionWorker := workers.NewTestConnectionWorker(controlPlaneClient)

	if enableLogging() {
		triggerWorker.SetLogger(logger)
		pollingWorker.SetLogger(logger)
		dataStoreTestConnectionWorker.SetLogger(logger)
	}

	controlPlaneClient.OnDataStoreTestConnectionRequest(dataStoreTestConnectionWorker.Test)
	controlPlaneClient.OnTriggerRequest(triggerWorker.Trigger)
	controlPlaneClient.OnPollingRequest(pollingWorker.Poll)
	controlPlaneClient.OnConnectionClosed(func(ctx context.Context, sr *proto.ShutdownRequest) error {
		fmt.Printf("Server terminated the connection with the agent. Reason: %s\n", sr.Reason)
		return controlPlaneClient.Close()
	})

	return controlPlaneClient, nil
}
