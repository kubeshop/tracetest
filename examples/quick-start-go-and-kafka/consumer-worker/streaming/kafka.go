package streaming

import (
	"fmt"

	"github.com/Shopify/sarama"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
)

const fixedPartition = 0

type Reader struct {
	consumer          sarama.Consumer
	partitionConsumer sarama.PartitionConsumer
	topic             string
}

func GetKafkaReader(brokerUrl, topic string) (*Reader, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	conn, err := sarama.NewConsumer([]string{brokerUrl}, config)
	if err != nil {
		return nil, fmt.Errorf("error when setting up kafka publisher: %w", err)
	}

	// Wrap instrumentation
	conn = otelsarama.WrapConsumer(conn)

	partitionConsumer, err := conn.ConsumePartition(topic, fixedPartition, sarama.OffsetOldest)
	if err != nil {
		return nil, fmt.Errorf("error when getting partition consumer: %w", err)
	}

	// Wrap instrumentation
	partitionConsumer = otelsarama.WrapPartitionConsumer(partitionConsumer)

	return &Reader{
		consumer:          conn,
		topic:             topic,
		partitionConsumer: partitionConsumer,
	}, nil
}

func (p *Reader) PartitionConsumer() sarama.PartitionConsumer {
	return p.partitionConsumer
}

func (p *Reader) Close() {
	p.partitionConsumer.Close()
	p.consumer.Close()
}
