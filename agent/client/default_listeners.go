package client

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/agent/proto"
)

func triggerListener(_ context.Context, _ *proto.TriggerRequest) error {
	fmt.Println("triggerListener")
	return nil
}

func pollListener(_ context.Context, _ *proto.PollingRequest) error {
	fmt.Println("pollListener")
	return nil
}

func shutdownListener(_ context.Context, _ *proto.ShutdownRequest) error {
	fmt.Println("shutdownListener")
	return nil
}

func dataStoreConnectionListener(_ context.Context, _ *proto.DataStoreConnectionTestRequest) error {
	fmt.Println("dataStoreConnectionListener")
	return nil
}
