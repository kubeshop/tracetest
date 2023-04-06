package subscription

import "time"

type Manager struct {
	subscriptions *BufferedSubscriptions
	cleanupPeriod time.Duration
}

func NewManager() *Manager {
	manager := &Manager{
		subscriptions: NewBufferedSubscriptions(),
		cleanupPeriod: 10 * time.Second,
	}

	manager.startCleanupProcess()

	return manager
}

func (m *Manager) Subscribe(resourceID string, subscriber Subscriber) {
	m.subscriptions.AddSubscriber(resourceID, subscriber)
}

func (m *Manager) Unsubscribe(resourceID string, subscriptionID string) {
	m.subscriptions.RemoveSubscriber(resourceID, subscriptionID)
}

func (m *Manager) PublishUpdate(message Message) {
	m.subscriptions.NotifySubscribers(message.ResourceID, message)
}

func (m *Manager) Publish(resourceID string, message any) {
	m.subscriptions.NotifySubscribers(resourceID, Message{ResourceID: resourceID, Content: message})
}

func (m *Manager) startCleanupProcess() {
	go func() {
		ticker := time.NewTicker(m.cleanupPeriod)
		select {
		case <-ticker.C:
			m.subscriptions.CleanUp()
		}
	}()
}
