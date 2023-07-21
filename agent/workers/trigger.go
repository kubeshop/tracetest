package workers

import (
	"context"

	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/proto"
)

type TriggerWorker struct {
	client *client.Client
}

func NewTriggerWorker(client *client.Client) *TriggerWorker {
	return &TriggerWorker{client}
}

func (w *TriggerWorker) Trigger(ctx context.Context, trigger *proto.TriggerRequest) error {
	return nil
}
