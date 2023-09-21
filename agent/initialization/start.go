package initialization

import (
	"context"
	"fmt"

	gocache "github.com/Code-Hex/go-generics-cache"
	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/workers"
	"github.com/kubeshop/tracetest/server/traces"
)

func NewClient(ctx context.Context, config config.Config) (*client.Client, error) {
	client, err := client.Connect(ctx, config.ServerURL,
		client.WithAPIKey(config.APIKey),
		client.WithAgentName(config.Name),
	)
	if err != nil {
		return nil, err
	}

	tracesCache := gocache.New[string, []traces.Span]()

	triggerWorker := workers.NewTriggerWorker(client, workers.WithTraceCache(tracesCache))
	pollingWorker := workers.NewPollerWorker(client)
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
func Start(ctx context.Context, config config.Config) error {
	client, err := NewClient(ctx, config)
	if err != nil {
		return err
	}

	err = client.Start(ctx)
	if err != nil {
		return err
	}

	client.WaitUntilDisconnected()
	return nil
}
