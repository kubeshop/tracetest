package trigger

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/test/trigger"
)

func Kafka() Triggerer {
	return &KafkaTriggerer{}
}

type KafkaTriggerer struct{}

func (t *KafkaTriggerer) Trigger(ctx context.Context, test test.Test, opts *TriggerOptions) (Response, error) {
	// do logic

	response := Response{
		Result: trigger.TriggerResult{
			Type: t.Type(),
			Kafka: &trigger.KafkaResponse{
				Partition: "",
				Offset:    "",
				Error:     "",
			},
		},
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
