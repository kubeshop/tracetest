package initialization

import (
	"context"
	"errors"
	"fmt"

	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/collector"
	"github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/workers"
	"github.com/kubeshop/tracetest/agent/workers/poller"
	"go.opentelemetry.io/otel/trace"
)

var ErrOtlpServerStart = errors.New("OTLP server start error")

func NewClient(ctx context.Context, config config.Config, traceCache collector.TraceCache) (*client.Client, error) {
	controlPlaneClient, err := client.Connect(ctx, config.ServerURL,
		client.WithAPIKey(config.APIKey),
		client.WithAgentName(config.Name),
	)
	if err != nil {
		return nil, err
	}

	triggerWorker := workers.NewTriggerWorker(controlPlaneClient, workers.WithTraceCache(traceCache))
	pollingWorker := workers.NewPollerWorker(controlPlaneClient, workers.WithInMemoryDatastore(
		poller.NewInMemoryDatastore(traceCache),
	))
	dataStoreTestConnectionWorker := workers.NewTestConnectionWorker(controlPlaneClient)

	controlPlaneClient.OnDataStoreTestConnectionRequest(dataStoreTestConnectionWorker.Test)
	controlPlaneClient.OnTriggerRequest(triggerWorker.Trigger)
	controlPlaneClient.OnPollingRequest(pollingWorker.Poll)
	controlPlaneClient.OnConnectionClosed(func(ctx context.Context, sr *proto.ShutdownRequest) error {
		fmt.Printf("Server terminated the connection with the agent. Reason: %s\n", sr.Reason)
		return controlPlaneClient.Close()
	})

	return controlPlaneClient, nil
}

// Start the agent with given configuration
func Start(ctx context.Context, cfg config.Config) (*Session, error) {
	traceCache := collector.NewTraceCache()
	controlPlaneClient, err := NewClient(ctx, cfg, traceCache)
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

	_, err := collector.Start(ctx, collectorConfig, noopTracer, collector.WithTraceCache(traceCache), collector.WithStartRemoteServer(false))
	if err != nil {
		return ErrOtlpServerStart
	}

	return nil
}
