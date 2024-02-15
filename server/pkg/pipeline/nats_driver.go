package pipeline

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

type NatsDriver[T any] struct {
	log          loggerFn
	conn         *nats.Conn
	topic        string
	subscription *nats.Subscription
}

func NewNatsDriver[T any](conn *nats.Conn, topic string) *NatsDriver[T] {
	return &NatsDriver[T]{
		log:   newLoggerFn(fmt.Sprintf("NatsDriver - %s", topic)),
		conn:  conn,
		topic: topic,
	}
}

func (d *NatsDriver[T]) Enqueue(ctx context.Context, msg T) {
	msgJson, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("could not marshal message: %s\n", err.Error())
	}

	err = d.conn.PublishMsg(&nats.Msg{
		Subject: d.topic,
		Data:    msgJson,
	})

	if err != nil {
		fmt.Printf("could not send publish message request: %s\n", err.Error())
	}
}

// SetListener implements QueueDriver.
func (d *NatsDriver[T]) SetListener(listener Listener[T]) {
	subscription, err := d.conn.Subscribe(d.topic, func(msg *nats.Msg) {
		var target T
		err := json.Unmarshal(msg.Data, &target)
		if err != nil {
			fmt.Printf(`could not unmarshal message got in queue "%s": %s\n`, d.topic, err.Error())
		}

		// TODO: We probably should return an error for acking or nacking this message
		listener.Listen(target)

		msg.Ack()
	})

	if err != nil {
		panic(err)
	}

	d.subscription = subscription
}

func (d *NatsDriver[T]) Start() {
	d.log("start")
}

func (d *NatsDriver[T]) Stop() {
	err := d.subscription.Unsubscribe()
	if err != nil {
		d.log(`could not unsubscribe to topic "%s"\n`, d.topic)
	}
	d.subscription = nil
}
