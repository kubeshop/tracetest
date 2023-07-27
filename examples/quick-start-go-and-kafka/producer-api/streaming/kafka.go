package streaming

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"go.opentelemetry.io/otel"
)

type Publisher struct {
	producer sarama.SyncProducer
	topic    string
}

func GetKafkaPublisher(brokerUrl, topic string) (*Publisher, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	conn, err := sarama.NewSyncProducer([]string{brokerUrl}, config)
	if err != nil {
		return nil, fmt.Errorf("error when setting up kafka publisher: %w", err)
	}

	// Wrap instrumentation
	conn = otelsarama.WrapSyncProducer(config, conn)

	return &Publisher{
		producer: conn,
		topic:    topic,
	}, nil
}

func (p *Publisher) Publish(ctx context.Context, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(message),
	}

	otel.GetTextMapPropagator().Inject(ctx, otelsarama.NewProducerMessageCarrier(msg))

	fmt.Printf("Sending message to topic: %s\n", p.topic)

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to kafka: %w", err)
	}
	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", p.topic, partition, offset)

	return nil
}

func (p *Publisher) Close() {
	p.producer.Close()
}
