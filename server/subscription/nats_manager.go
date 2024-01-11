package subscription

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

type natsManager struct {
	conn *nats.Conn
}

func (m *natsManager) Subscribe(resourceID string, subscriber Subscriber) {
	_, err := m.conn.Subscribe(resourceID, func(msg *nats.Msg) {
		decoded := Message{}
		err := json.Unmarshal(msg.Data, &decoded)
		if err != nil {
			panic(fmt.Errorf("cannot unmarshall incoming nats message: %w", err))
		}
		err = subscriber.Notify(decoded)
		if err != nil {
			panic(fmt.Errorf("cannot handle notification of nats message: %w", err))
		}
	})
	if err != nil {
		panic(fmt.Errorf("cannot subscribe to nats topic: %w", err))
	}
}

func (m *natsManager) Unsubscribe(resourceID string, subscriptionID string) {
	panic("nats unsubscribe not implemented")
}

func (m *natsManager) PublishUpdate(message Message) {
	bytes, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}

	err = m.conn.Publish(message.ResourceID, bytes)
	if err != nil {
		panic(err)
	}
}

func (m *natsManager) Publish(resourceID string, message any) {
	m.PublishUpdate(Message{
		ResourceID: resourceID,
		Content:    message,
	})
}
