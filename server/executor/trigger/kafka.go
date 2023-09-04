package trigger

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/pkg/kafka"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/test/trigger"

	"go.opentelemetry.io/otel/propagation"
)

func Kafka() Triggerer {
	return &KafkaTriggerer{}
}

type KafkaTriggerer struct{}

func (t *KafkaTriggerer) Trigger(ctx context.Context, test test.Test) (Response, error) {
	response := Response{
		Result: trigger.TriggerResult{
			Type: t.Type(),
		},
	}

	kafkaTriggerRequest := test.Trigger.Kafka
	kafkaConfig := t.getConfig(kafkaTriggerRequest)

	kafkaProducer, err := kafka.GetProducer(kafkaConfig)
	if err != nil {
		return response, fmt.Errorf("error when creating kafka producer: %w", err)
	}
	defer kafkaProducer.Close()

	messageHeaders := kafkaTriggerRequest.GetHeaderAsMap()
	propagators().Inject(ctx, propagation.MapCarrier(messageHeaders))

	result, err := kafkaProducer.ProduceSyncMessage(ctx, kafkaTriggerRequest.MessageKey, kafkaTriggerRequest.MessageValue, messageHeaders)
	if err != nil {
		return response, fmt.Errorf("error when sending message to kafka producer: %w", err)
	}

	response.Result.Kafka = &trigger.KafkaResponse{
		Partition: result.Partition,
		Offset:    result.Offset,
	}

	response.SpanAttributes = map[string]string{
		"tracetest.run.trigger.kafka.partition": result.Partition,
		"tracetest.run.trigger.kafka.offset":    result.Offset,
	}

	return response, nil
}

func (t *KafkaTriggerer) Type() trigger.TriggerType {
	return trigger.TriggerTypeKafka
}

func (t *KafkaTriggerer) Resolve(ctx context.Context, test test.Test, opts *TriggerOptions) (test.Test, error) {
	kafkaConfig := test.Trigger.Kafka
	if kafkaConfig == nil {
		return test, fmt.Errorf("no settings provided for kafka triggerer")
	}

	brokerURLs := []string{}
	for _, url := range kafkaConfig.BrokerURLs {
		newUrl, err := opts.Executor.ResolveStatement(WrapInQuotes(url, "\""))
		if err != nil {
			return test, err
		}

		brokerURLs = append(brokerURLs, newUrl)
	}

	kafkaConfig.BrokerURLs = brokerURLs

	headers := []trigger.KafkaMessageHeader{}
	for _, header := range kafkaConfig.Headers {
		key, err := opts.Executor.ResolveStatement(WrapInQuotes(header.Key, "\""))
		if err != nil {
			return test, err
		}

		value, err := opts.Executor.ResolveStatement(WrapInQuotes(header.Value, "\""))
		if err != nil {
			return test, err
		}

		kafkaHeader := trigger.KafkaMessageHeader{
			Key:   key,
			Value: value,
		}

		headers = append(headers, kafkaHeader)
	}
	kafkaConfig.Headers = headers

	if kafkaConfig.MessageKey != "" {
		messageKey, err := opts.Executor.ResolveStatement(WrapInQuotes(kafkaConfig.MessageKey, "'"))
		if err != nil {
			return test, err
		}

		kafkaConfig.MessageKey = messageKey
	}

	if kafkaConfig.Authentication != nil && kafkaConfig.Authentication.Plain != nil {
		username, err := opts.Executor.ResolveStatement(WrapInQuotes(kafkaConfig.Authentication.Plain.Username, "'"))
		if err != nil {
			return test, err
		}

		password, err := opts.Executor.ResolveStatement(WrapInQuotes(kafkaConfig.Authentication.Plain.Password, "'"))
		if err != nil {
			return test, err
		}

		kafkaConfig.Authentication = &trigger.KafkaAuthenticator{
			Type: "plain",
			Plain: &trigger.KafkaPlainAuthenticator{
				Username: username,
				Password: password,
			},
		}
	}

	test.Trigger.Kafka = kafkaConfig

	return test, nil
}

func (t *KafkaTriggerer) getConfig(request *trigger.KafkaRequest) kafka.Config {
	config := kafka.Config{
		BrokerURLs:      request.BrokerURLs,
		Topic:           request.Topic,
		SSLVerification: request.SSLVerification,
	}

	if request.Authentication == nil || request.Authentication.Plain == nil {
		return config
	}

	config.Authentication = &kafka.AuthenticationConfig{
		Username: request.Authentication.Plain.Username,
		Password: request.Authentication.Plain.Password,
	}

	return config
}
