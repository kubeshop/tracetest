package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/kubeshop/tracetest/server/subscription"
)

type unsubscribeMessage struct {
	Resource       string `json:"resource"`
	SubscriptionId string `json:"subscriptionId"`
}

type unsubscribeCommandExecutor struct {
	subscriptionManager subscription.Manager
}

func NewUnsubscribeCommandExecutor(manager subscription.Manager) MessageExecutor {
	return unsubscribeCommandExecutor{
		subscriptionManager: manager,
	}
}

func (e unsubscribeCommandExecutor) Execute(conn *websocket.Conn, message []byte) {
	msg := unsubscribeMessage{}
	err := json.Unmarshal(message, &msg)
	if err != nil {
		conn.WriteJSON(ErrorMessage(fmt.Errorf("invalid unsubscription message: %w", err)))
		return
	}

	if msg.Resource == "" {
		conn.WriteJSON(ErrorMessage(fmt.Errorf("resource cannot be empty")))
		return
	}

	e.subscriptionManager.Unsubscribe(msg.Resource, msg.SubscriptionId)

	conn.WriteJSON(UnsubscribeSuccess())
}
