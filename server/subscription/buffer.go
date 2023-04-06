package subscription

import (
	"sync"
	"time"
)

type TemporaryBuffer struct {
	messages             []Message
	lastMessageTimestamp time.Time
}

func NewTemporaryBuffer() *TemporaryBuffer {
	return &TemporaryBuffer{
		messages:             make([]Message, 0),
		lastMessageTimestamp: time.Now(),
	}
}

func (b *TemporaryBuffer) AddMessage(message Message) {
	b.messages = append(b.messages, message)
	b.lastMessageTimestamp = time.Now()
}

func (b *TemporaryBuffer) GetMessages() []Message {
	return b.messages
}

func (b *TemporaryBuffer) CleanUp() {
	if len(b.messages) > 0 && time.Since(b.lastMessageTimestamp) > 10*time.Second {
		b.messages = make([]Message, 0)
		b.lastMessageTimestamp = time.Now()
	}
}

type BufferedSubscriptions struct {
	subscriptions  map[string][]Subscriber
	messageBuffers map[string]*TemporaryBuffer
	mutex          *sync.Mutex
}

func NewBufferedSubscriptions() *BufferedSubscriptions {
	var mutex sync.Mutex
	return &BufferedSubscriptions{
		subscriptions:  make(map[string][]Subscriber),
		messageBuffers: make(map[string]*TemporaryBuffer),
		mutex:          &mutex,
	}
}

func (bs *BufferedSubscriptions) CleanUp() {
	for _, buffer := range bs.messageBuffers {
		buffer.CleanUp()
	}
}

func (bs *BufferedSubscriptions) NotifySubscribers(resourceID string, message Message) {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	bs.initResource(resourceID)

	bs.messageBuffers[resourceID].AddMessage(message)

	if subscribers, ok := bs.subscriptions[resourceID]; ok {
		for _, subscriber := range subscribers {
			subscriber.Notify(message)
		}
	}
}

func (bs *BufferedSubscriptions) AddSubscriber(resourceID string, subscriber Subscriber) {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	bs.initResource(resourceID)

	bs.subscriptions[resourceID] = append(bs.subscriptions[resourceID], subscriber)

	for _, message := range bs.messageBuffers[resourceID].GetMessages() {
		subscriber.Notify(message)
	}
}

func (bs *BufferedSubscriptions) RemoveSubscriber(resourceID string, subscriptionID string) {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	if _, ok := bs.subscriptions[resourceID]; !ok {
		return
	}

	subscribers := bs.subscriptions[resourceID]

	newArray := make([]Subscriber, 0, len(subscribers)-1)
	for _, item := range subscribers {
		if item.ID() != subscriptionID {
			newArray = append(newArray, item)
		}
	}

	bs.subscriptions[resourceID] = newArray
}

func (bs *BufferedSubscriptions) initResource(resourceID string) {
	if _, ok := bs.messageBuffers[resourceID]; !ok {
		bs.messageBuffers[resourceID] = NewTemporaryBuffer()
	}

	if _, ok := bs.subscriptions[resourceID]; !ok {
		bs.subscriptions[resourceID] = make([]Subscriber, 0)
	}

}
