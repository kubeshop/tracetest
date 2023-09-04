package initialization

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/workers"
)

// Start the agent with given configuration
func Start(ctx context.Context, config config.Config) error {
	fmt.Println("Starting agent")
	client, err := client.Connect(ctx, config.ServerURL,
		client.WithAPIKey(config.APIKey),
		client.WithAgentName(config.Name),
	)
	if err != nil {
		return err
	}

	triggerWorker := workers.NewTriggerWorker(client)
	pollingWorker := workers.NewPollerWorker(client)

	client.OnTriggerRequest(triggerWorker.Trigger)
	client.OnPollingRequest(pollingWorker.Poll)
	client.OnConnectionClosed(func(ctx context.Context, sr *proto.ShutdownRequest) error {
		fmt.Printf("Server terminated the connection with the agent. Reason: %s\n", sr.Reason)
		return client.Close()
	})

	err = client.Start(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Agent started! Do not close the terminal.")
	client.WaitUntilDisconnected()
	return nil
}
