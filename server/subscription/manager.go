package subscription

import "sync"

type Manager struct {
	subscriptions map[string][]Subscriber
	mutex         sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		subscriptions: make(map[string][]Subscriber),
		mutex:         sync.Mutex{},
	}
}

func (m *Manager) Subscribe(resourceID string, subscriber Subscriber) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	array := make([]Subscriber, 0)

	if existingArray, ok := m.subscriptions[resourceID]; ok {
		array = existingArray
	}

	array = append(array, subscriber)
	m.subscriptions[resourceID] = array
}

func (m *Manager) Unsubscribe(resourceID string, subscriptionID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	array, exists := m.subscriptions[resourceID]
	if !exists {
		return
	}

	newArray := make([]Subscriber, 0, len(array)-1)
	for _, item := range array {
		if item.ID() != subscriptionID {
			newArray = append(newArray, item)
		}
	}

	m.subscriptions[resourceID] = newArray
}

func (m *Manager) PublishUpdate(message Message) {
	if subscribers, ok := m.subscriptions[message.ResourceID]; ok {
		for _, subscriber := range subscribers {
			subscriber.Notify(message)
		}
	}
}

func (m *Manager) Publish(resourceID string, message any) {
	if subscribers, ok := m.subscriptions[resourceID]; ok {
		for _, subscriber := range subscribers {
			subscriber.Notify(Message{
				ResourceID: resourceID,
				Content:    message,
			})
		}
	}
}
