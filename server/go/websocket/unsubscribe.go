package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/kubeshop/tracetest/server/go/subscription"
)

type unsubscribeMessage struct {
	Resource       string `json:"resource"`
	SubscriptionId string `json:"subscriptionId"`
}

func HandleUnsubscribeCommand(conn *websocket.Conn, message []byte) {
	msg := unsubscribeMessage{}
	err := json.Unmarshal(message, &msg)
	if err != nil {
		conn.WriteJSON(ErrorMessage(fmt.Errorf("invalid unsubscription message: %w", err)))
		return
	}

	if msg.Resource == "" {
		conn.WriteJSON(ErrorMessage(fmt.Errorf("Resource cannot be empty")))
		return
	}

	manager := subscription.GetManager()
	manager.Unsubscribe(msg.Resource, msg.SubscriptionId)

	conn.WriteJSON(UnsubscribeSuccess())
}
