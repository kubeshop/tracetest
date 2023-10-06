package trigger

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/pkg/kafka"
	"github.com/kubeshop/tracetest/server/test/trigger"
	"go.opentelemetry.io/otel/propagation"
)

func KAFKA() Triggerer {
	return &kafkaTriggerer{}
}

type kafkaTriggerer struct{}

func (te *kafkaTriggerer) Trigger(ctx context.Context, triggerConfig trigger.Trigger, opts *Options) (Response, error) {
	response := Response{
		Result: trigger.TriggerResult{
			Type: te.Type(),
		},
	}

	kafkaTriggerRequest := triggerConfig.Kafka
	kafkaConfig := te.getConfig(kafkaTriggerRequest)

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

func (t *kafkaTriggerer) Type() trigger.TriggerType {
	return trigger.TriggerTypeKafka
}

func (t *kafkaTriggerer) getConfig(request *trigger.KafkaRequest) kafka.Config {
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
