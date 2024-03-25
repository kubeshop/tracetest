package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
)

const producerMaxRetries = 5

type Config struct {
	BrokerURLs      []string
	Topic           string
	Authentication  *AuthenticationConfig
	SSLVerification bool
}

type AuthenticationConfig struct {
	Username string
	Password string
}

type Producer struct {
	internalProducer sarama.SyncProducer
	topic            string
}

type ProduceMessageResult struct {
	Topic     string
	Partition string
	Offset    string
}

func GetProducer(kafkaConfig Config) (*Producer, error) {
	saramaConfig := getSaramaConfig(kafkaConfig)

	conn, err := sarama.NewSyncProducer(kafkaConfig.BrokerURLs, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("error when setting up kafka publisher: %w", err)
	}

	return &Producer{
		internalProducer: conn,
		topic:            kafkaConfig.Topic,
	}, nil
}

func (p *Producer) ProduceSyncMessage(ctx context.Context, messageKey, messageValue string, headers map[string]string) (ProduceMessageResult, error) {
	msg := &sarama.ProducerMessage{
		Topic:   p.topic,
		Key:     sarama.StringEncoder(messageKey),
		Value:   sarama.StringEncoder(messageValue),
		Headers: []sarama.RecordHeader{},
	}

	for headerKey, headerValue := range headers {
		msg.Headers = append(msg.Headers, sarama.RecordHeader{
			Key:   []byte(headerKey),
			Value: []byte(headerValue),
		})
	}

	partition, offset, err := p.internalProducer.SendMessage(msg)
	if err != nil {
		return ProduceMessageResult{}, fmt.Errorf("failed to send message to kafka: %w", err)
	}

	return ProduceMessageResult{
		Topic:     p.topic,
		Partition: fmt.Sprintf("%d", partition),
		Offset:    fmt.Sprintf("%d", offset),
	}, nil
}

func (p *Producer) Close() {
	p.internalProducer.Close()
}

func getSaramaConfig(kafkaConfig Config) *sarama.Config {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = producerMaxRetries

	if kafkaConfig.SSLVerification {
		saramaConfig.Net.TLS.Enable = true
	}

	if kafkaConfig.Authentication != nil {
		saramaConfig.Net.SASL.Enable = true
		saramaConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		saramaConfig.Net.SASL.User = kafkaConfig.Authentication.Username
		saramaConfig.Net.SASL.Password = kafkaConfig.Authentication.Password
	}

	return saramaConfig
}
