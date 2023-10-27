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
	c, err := client.Connect(ctx, config.ServerURL,
		client.WithAPIKey(config.APIKey),
		client.WithAgentName(config.Name),
	)
	if err != nil {
		return nil, err
	}

	triggerWorker := workers.NewTriggerWorker(c, workers.WithTraceCache(traceCache))
	pollingWorker := workers.NewPollerWorker(c, workers.WithInMemoryDatastore(
		poller.NewInMemoryDatastore(traceCache),
	))
	dataStoreTestConnectionWorker := workers.NewTestConnectionWorker(c)

	c.OnDataStoreTestConnectionRequest(dataStoreTestConnectionWorker.Test)
	c.OnTriggerRequest(triggerWorker.Trigger)
	c.OnPollingRequest(pollingWorker.Poll)
	c.OnConnectionClosed(func(ctx context.Context, sr *proto.ShutdownRequest) error {
		fmt.Printf("Server terminated the connection with the agent. Reason: %s\n", sr.Reason)
		return c.Close()
	})

	return c, nil
}

// Start the agent with given configuration
func Start(ctx context.Context, cfg config.Config) (*Session, error) {
	traceCache := collector.NewTraceCache()
	c, err := NewClient(ctx, cfg, traceCache)
	if err != nil {
		return nil, err
	}

	err = c.Start(ctx)
	if err != nil {
		return nil, err
	}

	err = StartCollector(ctx, cfg, traceCache)
	if err != nil {
		return nil, err
	}

	return &Session{
		client: c,
		Token:  c.SessionConfiguration().AgentIdentification.Token,
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
