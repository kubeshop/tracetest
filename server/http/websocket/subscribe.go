package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/kubeshop/tracetest/server/http/mappings"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/transactions"
)

type subscriptionMessage struct {
	Resource string `json:"resource"`
}

type subscribeCommandExecutor struct {
	subscriptionManager *subscription.Manager
	mappers             mappings.Mappings
}

func NewSubscribeCommandExecutor(manager *subscription.Manager, mappers mappings.Mappings) MessageExecutor {
	return subscribeCommandExecutor{
		subscriptionManager: manager,
		mappers:             mappers,
	}
}

func (e subscribeCommandExecutor) Execute(conn *websocket.Conn, message []byte) {
	msg := subscriptionMessage{}
	err := json.Unmarshal(message, &msg)
	if err != nil {
		conn.WriteJSON(ErrorMessage(fmt.Errorf("invalid subscription message: %w", err)))
		return
	}

	if msg.Resource == "" {
		conn.WriteJSON(ErrorMessage(fmt.Errorf("resource cannot be empty")))
		return
	}

	messageConverter := subscription.NewSubscriberFunction(func(m subscription.Message) error {
		event := e.ResourceUpdatedEvent(m.Content)
		event.Resource = msg.Resource
		err := conn.WriteJSON(event)
		if err != nil {
			return fmt.Errorf("could not send update message: %w", err)
		}

		return nil
	})

	e.subscriptionManager.Subscribe(msg.Resource, messageConverter)

	conn.WriteJSON(SubscriptionSuccess(msg.Resource, messageConverter.ID()))
}

func (e subscribeCommandExecutor) ResourceUpdatedEvent(resource interface{}) Event {
	var mapped interface{}
	switch v := resource.(type) {
	case model.Run:
		mapped = e.mappers.Out.Run(&v)
	case *model.Run:
		mapped = e.mappers.Out.Run(v)
	case transactions.TransactionRun:
		mapped = e.mappers.Out.TransactionRun(v)
	case *transactions.TransactionRun:
		mapped = e.mappers.Out.TransactionRun(*v)
	case model.TestRunEvent:
		mapped = e.mappers.Out.TestRunEvent(v)
	default:
		fmt.Printf("type %T mapping not supported\n", v)
		mapped = v
	}

	return Event{
		Type:  "update",
		Event: mapped,
	}
}
