package subscription

import "github.com/google/uuid"

type Subscriber interface {
	ID() string
	Notify(message Message) error
}

type SubscriberFunction struct {
	id       string
	function func(Message) error
}

func (sf *SubscriberFunction) ID() string {
	return sf.id
}

func (sf *SubscriberFunction) Notify(message Message) error {
	return sf.function(message)
}

func NewSubscriberFunction(function func(Message) error) Subscriber {
	return &SubscriberFunction{
		id:       uuid.New().String(),
		function: function,
	}
}
