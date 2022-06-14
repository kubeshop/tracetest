package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/kubeshop/tracetest/server/subscription"
)

type subscriptionMessage struct {
	Resource string `json:"resource"`
}

type subscribeCommandExecutor struct {
	subscriptionManager *subscription.Manager
}

func NewSubscribeCommandExecutor(manager *subscription.Manager) MessageExecutor {
	return subscribeCommandExecutor{
		subscriptionManager: manager,
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
		err := conn.WriteJSON(ResourceUpdatedEvent(m.Content))
		if err != nil {
			return fmt.Errorf("could not send update message: %w", err)
		}

		return nil
	})

	e.subscriptionManager.Subscribe(msg.Resource, messageConverter)

	conn.WriteJSON(SubscriptionSuccess(messageConverter.ID()))
}
