package initialization

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/collector"
	"github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/workers"
	"github.com/kubeshop/tracetest/agent/workers/poller"
	"go.opentelemetry.io/otel/trace"
)

func NewClient(ctx context.Context, config config.Config, traceCache collector.TraceCache) (*client.Client, error) {
	client, err := client.Connect(ctx, config.ServerURL,
		client.WithAPIKey(config.APIKey),
		client.WithAgentName(config.Name),
	)
	if err != nil {
		return nil, err
	}

	triggerWorker := workers.NewTriggerWorker(client, workers.WithTraceCache(traceCache))
	pollingWorker := workers.NewPollerWorker(client, workers.WithInMemoryDatastore(
		poller.NewInMemoryDatastore(traceCache),
	))
	dataStoreTestConnectionWorker := workers.NewTestConnectionWorker(client)

	client.OnDataStoreTestConnectionRequest(dataStoreTestConnectionWorker.Test)
	client.OnTriggerRequest(triggerWorker.Trigger)
	client.OnPollingRequest(pollingWorker.Poll)
	client.OnConnectionClosed(func(ctx context.Context, sr *proto.ShutdownRequest) error {
		fmt.Printf("Server terminated the connection with the agent. Reason: %s\n", sr.Reason)
		return client.Close()
	})

	return client, nil
}

// Start the agent with given configuration
func Start(ctx context.Context, config config.Config) (*Session, error) {
	traceCache := collector.NewTraceCache()
	client, err := NewClient(ctx, config, traceCache)
	if err != nil {
		return nil, err
	}

	err = client.Start(ctx)
	if err != nil {
		return nil, err
	}

	err = StartCollector(ctx, config, traceCache)
	if err != nil {
		return nil, err
	}

	return &Session{
		client: client,
		Token:  client.SessionConfiguration().AgentIdentification.Token,
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
		return err
	}

	return nil
}
