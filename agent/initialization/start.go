package initialization

import (
	"context"
	"fmt"
	"log"

	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/agent/workers"
)

// Start the agent with given configuration
func Start(config config.Config) {
	fmt.Println("Starting agent")
	ctx := context.Background()

	client, err := client.Connect(ctx, config.ServerURL,
		client.WithAPIKey(config.APIKey),
		client.WithAgentName(config.AgentName),
	)
	if err != nil {
		log.Fatal(err)
	}

	triggerWorker := workers.NewTriggerWorker(client)
	pollingWorker := workers.NewPollerWorker(client)

	client.OnTriggerRequest(triggerWorker.Trigger)
	client.OnPollingRequest(pollingWorker.Poll)

	err = client.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}