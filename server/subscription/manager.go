package subscription

import (
	"sync"

	"github.com/nats-io/nats.go"
)

type Manager interface {
	Subscribe(resourceID string, subscriber Subscriber)
	Unsubscribe(resourceID string, subscriptionID string)

	PublishUpdate(message Message)
	Publish(resourceID string, message any)
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
		return &natsManager{currentOptions.conn}
	}

	return &inMemoryManager{
		subscriptions: make(map[string][]Subscriber),
		mutex:         sync.Mutex{},
	}
}
