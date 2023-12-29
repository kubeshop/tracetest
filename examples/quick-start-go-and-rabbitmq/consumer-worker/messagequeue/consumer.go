package streaming

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
)

const consumerGroupName = "reader"

type Reader struct {
	connection amqp.Connection
	channel    amqp.Channel
	queue      amqp.Queue
}

func GetMessageQueueReader(queueConnectionString, queueName string) (*Reader, error) {
	connection, err := amqp.Dial(queueConnectionString)
	if err != nil {
		return nil, fmt.Errorf("error when connecting on queue system to set up message queue publisher: %w", err)
	}

	channel, err := connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("error when creating channel to set up message queue publisher: %w", err)
	}

	queue, err := channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	return &Reader{
		connection: connection,
		channel:    channel,
		queue:      queue,
	}, nil
}

func (r *Reader) Read(ctx context.Context, messageReader func(context.Context, string)) error {
	deliveryChannel, err := r.channel.ConsumeWithContext(
		ctx,          // context
		r.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return fmt.Errorf("error when reading message queue: %w", err)
	}

	// read messages in a separated go proc
	go func() {
		for message := range deliveryChannel {
			ctx = otel.GetTextMapPropagator().Extract(ctx, message.Headers)

			fmt.Println("Message received from queue %s", r.queue.Name)
			fmt.Println("Headers: ", message.Headers)

			messageReader(ctx, string(message.Body))
		}
	}()

	return nil
}
