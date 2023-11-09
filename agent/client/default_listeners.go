package client

import (
	"context"

	"github.com/kubeshop/tracetest/agent/proto"
)

func triggerListener(_ context.Context, _ *proto.TriggerRequest) error {
	return nil
}

func pollListener(_ context.Context, _ *proto.PollingRequest) error {
	return nil
}

func shutdownListener(_ context.Context, _ *proto.ShutdownRequest) error {
	return nil
}

func dataStoreConnectionListener(_ context.Context, _ *proto.DataStoreConnectionTestRequest) error {
	return nil
}
