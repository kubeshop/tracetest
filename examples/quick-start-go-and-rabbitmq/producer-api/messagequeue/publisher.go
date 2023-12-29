package messagequeue

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"go.opentelemetry.io/otel"
)

const PUBLISH_TIMEOUT = 1 * time.Minute

type Publisher struct {
	connection amqp.Connection
	channel    amqp.Channel
	queue      amqp.Queue
}

func GetMessageQueuePublisher(queueConnectionString, queueName string) (*Publisher, error) {
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

	return &Publisher{
		connection: connection,
		channel:    channel,
		queue:      queue,
	}, nil
}

func (p *Publisher) Publish(ctx context.Context, message string) error {
	ctx, cancel := context.WithTimeout(ctx, PUBLISH_TIMEOUT)
	defer cancel()

	messageToPublish := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	}

	otel.GetTextMapPropagator().Inject(ctx, messageToPublish.Headers)

	fmt.Printf("Sending message to queue: %s\n", p.queue.Name)
	fmt.Println("Headers: ", messageToPublish.Headers)

	err := p.channel.PublishWithContext(ctx,
		"",           // exchange
		p.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		messageToPublish,
	)
	if err != nil {
		return fmt.Errorf("failed to send message to message queue: %w", err)
	}

	fmt.Printf("Message published\n")

	return nil
}

func (p *Publisher) Close() {
	p.channel.Close()
	p.connection.Close()
}
