package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/kubeshop/tracetest/server/go/subscription"
)

type subscriptionMessage struct {
	ResourceType string `json:"resourceType"`
	ResourceID   string `json:"resourceId"`
}

func HandleSubscribeCommand(conn *websocket.Conn, message []byte) {
	msg := subscriptionMessage{}
	err := json.Unmarshal(message, &msg)
	if err != nil {
		conn.WriteJSON(Error(fmt.Errorf("invalid subscription message: %w", err)))
		return
	}

	if msg.ResourceID == "" || msg.ResourceType == "" {
		conn.WriteJSON(Error(fmt.Errorf("either ResourceType or ResourceID is empty")))
		return
	}

	messageConverter := subscription.NewSubscriberFunction(func(m *subscription.Message) error {
		err := conn.WriteJSON(UpdateMessage(m.Content))
		if err != nil {
			return fmt.Errorf("could not send update message: %w", err)
		}

		return nil
	})

	manager := subscription.GetManager()
	resourceName := fmt.Sprintf("%s:%s", msg.ResourceType, msg.ResourceID)
	manager.Subscribe(resourceName, messageConverter)

	conn.WriteJSON(SuccessMessage("susbcribe"))
}
