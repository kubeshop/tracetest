package subscription

import (
	"fmt"
	"log"
	"sync"
)

type inMemoryManager struct {
	subscriptions map[string][]Subscriber
	mutex         sync.Mutex
}

func (m *inMemoryManager) getSubscribers(resourceID string) []Subscriber {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.subscriptions[resourceID]
}

func (m *inMemoryManager) setSubscribers(resourceID string, subscribers []Subscriber) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.subscriptions[resourceID] = subscribers
}

func (m *inMemoryManager) Subscribe(resourceID string, subscriber Subscriber) {
	subscribers := append(m.getSubscribers(resourceID), subscriber)
	m.setSubscribers(resourceID, subscribers)
}

func (m *inMemoryManager) Unsubscribe(resourceID string, subscriptionID string) {
	subscribers := m.getSubscribers(resourceID)

	updated := make([]Subscriber, 0, len(subscribers)-1)
	for _, item := range subscribers {
		if item.ID() != subscriptionID {
			updated = append(updated, item)
		}
	}

	m.setSubscribers(resourceID, updated)
}

func (m *inMemoryManager) PublishUpdate(message Message) {
	subscribers := m.getSubscribers(message.ResourceID)

	// in order to keep compatibility with the nats manager
	// we need to transcode the messages
	transcoded, err := message.EncodeContent()
	if err != nil {
		log.Printf("cannot transcode message to publish: %s", err.Error())
		return
	}
	fmt.Println("ACA", string(transcoded.Content.([]byte)))

	for _, subscriber := range subscribers {
		subscriber.Notify(transcoded)
	}
}

func (m *inMemoryManager) Publish(resourceID string, message any) {
	m.PublishUpdate(Message{
		ResourceID: resourceID,
		Content:    message,
	})
}
