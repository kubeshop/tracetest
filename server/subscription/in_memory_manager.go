package subscription

import (
	"log"
	"sync"
)

type inMemoryManager struct {
	subscribers   map[string][]Subscriber
	subscriptions *subscriptionStorage
	mutex         sync.Mutex
}

type inMemorySubscription struct {
	unsubscribeFn func()
}

func (s *inMemorySubscription) Unsubscribe() error {
	s.unsubscribeFn()
	return nil
}

func (m *inMemoryManager) getSubscribers(resourceID string) []Subscriber {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.subscribers[resourceID]
}

func (m *inMemoryManager) setSubscribers(resourceID string, subscribers []Subscriber) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.subscribers[resourceID] = subscribers
}

func (m *inMemoryManager) Subscribe(resourceID string, subscriber Subscriber) Subscription {
	subscribers := append(m.getSubscribers(resourceID), subscriber)
	m.setSubscribers(resourceID, subscribers)

	return &inMemorySubscription{
		unsubscribeFn: func() { m.unsubscribe(resourceID, subscriber.ID()) },
	}
}

func (m *inMemoryManager) GetSubscription(resourceID string, subscriptionID string) Subscription {
	return m.subscriptions.Get(resourceID, subscriptionID)
}

func (m *inMemoryManager) unsubscribe(resourceID string, subscriptionID string) {
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

	for _, subscriber := range subscribers {
		err := subscriber.Notify(transcoded)
		if err != nil {
			log.Println("error notifying subscriber: ", err.Error())
		}
	}
}

func (m *inMemoryManager) Publish(resourceID string, message any) {
	m.PublishUpdate(Message{
		ResourceID: resourceID,
		Content:    message,
	})
}
