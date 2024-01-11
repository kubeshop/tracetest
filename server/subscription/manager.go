package subscription

import "sync"

type Manager interface {
	Subscribe(resourceID string, subscriber Subscriber)
	Unsubscribe(resourceID string, subscriptionID string)

	PublishUpdate(message Message)
	Publish(resourceID string, message any)
}

func NewManager() Manager {
	return &manager{
		subscriptions: make(map[string][]Subscriber),
		mutex:         sync.Mutex{},
	}
}

type manager struct {
	subscriptions map[string][]Subscriber
	mutex         sync.Mutex
}

func (m *manager) getSubscribers(resourceID string) []Subscriber {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.subscriptions[resourceID]
}

func (m *manager) setSubscribers(resourceID string, subscribers []Subscriber) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.subscriptions[resourceID] = subscribers
}

func (m *manager) Subscribe(resourceID string, subscriber Subscriber) {
	subscribers := append(m.getSubscribers(resourceID), subscriber)
	m.setSubscribers(resourceID, subscribers)
}

func (m *manager) Unsubscribe(resourceID string, subscriptionID string) {
	subscribers := m.getSubscribers(resourceID)

	updated := make([]Subscriber, 0, len(subscribers)-1)
	for _, item := range subscribers {
		if item.ID() != subscriptionID {
			updated = append(updated, item)
		}
	}

	m.setSubscribers(resourceID, updated)
}

func (m *manager) PublishUpdate(message Message) {
	subscribers := m.getSubscribers(message.ResourceID)

	for _, subscriber := range subscribers {
		subscriber.Notify(message)
	}
}

func (m *manager) Publish(resourceID string, message any) {
	m.PublishUpdate(Message{
		ResourceID: resourceID,
		Content:    message,
	})
}
