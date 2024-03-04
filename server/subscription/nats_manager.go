package subscription

import (
	"log"

	"github.com/nats-io/nats.go"
)

type natsManager struct {
	conn          *nats.Conn
	subscriptions *subscriptionStorage
}

func (m *natsManager) Subscribe(resourceID string, subscriber Subscriber) Subscription {
	subscription, err := m.conn.QueueSubscribe(resourceID, "subscriptions", func(msg *nats.Msg) {
		decoded, err := DecodeMessage(msg.Data)
		if err != nil {
			log.Printf("cannot unmarshall incoming nats message: %s", err.Error())
			return
		}
		err = subscriber.Notify(decoded)
		if err != nil {
			log.Printf("cannot handle notification of nats message: %s", err.Error())
			return
		}
	})
	if err != nil {
		log.Printf("cannot subscribe to nats topic: %s", err.Error())
		return nil
	}

	return subscription
}

func (m *natsManager) GetSubscription(resourceID string, subscriptionID string) Subscription {
	return m.subscriptions.Get(resourceID, subscriptionID)
}

func (m *natsManager) PublishUpdate(message Message) {
	bytes, err := message.Encode()
	if err != nil {
		log.Printf("cannot marshal message to publish nats message: %s", err.Error())
		return
	}

	err = m.conn.Publish(message.ResourceID, bytes)
	if err != nil {
		log.Printf("cannot publish nats message: %s", err.Error())
		return
	}
}

func (m *natsManager) Publish(resourceID string, message any) {
	m.PublishUpdate(Message{
		ResourceID: resourceID,
		Content:    message,
	})
}
