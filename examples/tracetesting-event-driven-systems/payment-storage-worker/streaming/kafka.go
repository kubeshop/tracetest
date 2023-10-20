package streaming

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"go.opentelemetry.io/otel"
)

const consumerGroupName = "reader"

type Reader struct {
	consumerGroup sarama.ConsumerGroup
	topic         string
}

func GetKafkaReader(brokerUrl, topic string) (*Reader, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	config.Consumer.Return.Errors = true

	brokerUrls := []string{brokerUrl}

	consumerGroup, err := sarama.NewConsumerGroup(brokerUrls, consumerGroupName, config)
	if err != nil {
		return nil, fmt.Errorf("error when starting consumer group: %w", err)
	}

	return &Reader{
		consumerGroup: consumerGroup,
		topic:         topic,
	}, nil
}

func (r *Reader) Read(ctx context.Context, messageReader func(context.Context, string, string)) error {
	consumerGroupHandler := internalConsumer{
		messageReader: messageReader,
	}

	// Wrap instrumentation
	handler := otelsarama.WrapConsumerGroupHandler(&consumerGroupHandler)

	err := r.consumerGroup.Consume(ctx, []string{r.topic}, handler)
	if err != nil {
		return fmt.Errorf("error when consuming via handler: %w", err)
	}

	return nil
}

// based on https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/Shopify/sarama/otelsarama/example/consumer/consumer.go

// Represents a Sarama consumer group consumer.
type internalConsumer struct {
	messageReader func(context.Context, string, string)
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (consumer *internalConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited.
func (consumer *internalConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *internalConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29

	for message := range claim.Messages() {
		ctx := otel.GetTextMapPropagator().Extract(context.Background(), otelsarama.NewConsumerMessageCarrier(message))

		fmt.Println("Headers: ", message.Headers)
		consumer.messageReader(ctx, string(message.Topic), string(message.Value))
	}

	return nil
}
