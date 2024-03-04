package pipeline

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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

	header := make(nats.Header)
	if propagator := otel.GetTextMapPropagator(); propagator != nil {
		propagator.Inject(ctx, propagation.HeaderCarrier(header))
	}

	err = d.conn.PublishMsg(&nats.Msg{
		Subject: d.topic,
		Header:  header,
		Data:    msgJson,
	})

	if err != nil {
		fmt.Printf("could not send publish message request: %s\n", err.Error())
	}
}

// SetListener implements QueueDriver.
func (d *NatsDriver[T]) SetListener(listener Listener[T]) {
	subscription, err := d.conn.QueueSubscribe(d.topic, "queue_worker", func(msg *nats.Msg) {
		var target T
		err := json.Unmarshal(msg.Data, &target)
		if err != nil {
			fmt.Printf(`could not unmarshal message got in queue "%s": %s\n`, d.topic, err.Error())
		}

		ctx := context.Background()
		if propagator := otel.GetTextMapPropagator(); propagator != nil {
			ctx = propagator.Extract(ctx, propagation.HeaderCarrier(msg.Header))
		}

		// TODO: We probably should return an error for acking or nacking this message
		listener.Listen(ctx, target)

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
