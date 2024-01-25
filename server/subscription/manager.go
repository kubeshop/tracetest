package subscription

import (
	"fmt"
	"sync"

	"github.com/nats-io/nats.go"
)

type Manager interface {
	Subscribe(resourceID string, subscriber Subscriber) Subscription
	GetSubscription(resourceID string, subscriptionID string) Subscription

	PublishUpdate(message Message)
	Publish(resourceID string, message any)
}

type Subscription interface {
	Unsubscribe() error
}

type optFn func(*options)

type options struct {
	conn *nats.Conn
}

func WithNats(conn *nats.Conn) optFn {
	return func(o *options) {
		o.conn = conn
	}
}

func NewManager(opts ...optFn) Manager {

	currentOptions := options{}
	for _, opt := range opts {
		opt(&currentOptions)
	}

	if currentOptions.conn != nil {
		return &natsManager{
			currentOptions.conn,
			newSubscriptionStorage(),
		}
	}

	return &inMemoryManager{
		subscribers:   make(map[string][]Subscriber),
		subscriptions: newSubscriptionStorage(),
		mutex:         sync.Mutex{},
	}
}

type subscriptionStorage struct {
	subscriptions map[string]Subscription
}

func newSubscriptionStorage() *subscriptionStorage {
	return &subscriptionStorage{
		subscriptions: make(map[string]Subscription),
	}
}

func (s *subscriptionStorage) Get(resourceID, subscriberID string) Subscription {
	key := s.key(resourceID, subscriberID)
	return s.subscriptions[key]
}

func (s *subscriptionStorage) Set(resourceID, subscriberID string, subscription Subscription) {
	key := s.key(resourceID, subscriberID)
	s.subscriptions[key] = subscription
}

func (s *subscriptionStorage) key(resourceID, subscriberID string) string {
	return fmt.Sprintf("%s-%s", resourceID, subscriberID)
}
