package initialization

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/workers"
)

func NewClient(ctx context.Context, config config.Config) (*client.Client, error) {
	client, err := client.Connect(ctx, config.ServerURL,
		client.WithAPIKey(config.APIKey),
		client.WithAgentName(config.Name),
	)
	if err != nil {
		return nil, err
	}

	triggerWorker := workers.NewTriggerWorker(client)
	pollingWorker := workers.NewPollerWorker(client)

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

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-ch
		fmt.Println("SIGTERM received: closing grpc connection")
		client.Close()
		fmt.Println("client connection was closed")
	}()

	client.WaitUntilDisconnected()
	return nil
}
