package trigger

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/agent/workers/trigger/kafka"
	"go.opentelemetry.io/otel/propagation"
)

func KAFKA() Triggerer {
	return &kafkaTriggerer{}
}

type kafkaTriggerer struct{}

func (te *kafkaTriggerer) Trigger(ctx context.Context, triggerConfig Trigger, opts *Options) (Response, error) {
	response := Response{
		Result: TriggerResult{
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

	response.Result.Kafka = &KafkaResponse{
		Partition: result.Partition,
		Offset:    result.Offset,
	}

	response.SpanAttributes = map[string]string{
		"tracetest.run.trigger.kafka.partition": result.Partition,
		"tracetest.run.trigger.kafka.offset":    result.Offset,
	}

	return response, nil
}

func (t *kafkaTriggerer) Type() TriggerType {
	return TriggerTypeKafka
}

func (t *kafkaTriggerer) getConfig(request *KafkaRequest) kafka.Config {
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

const TriggerTypeKafka TriggerType = "kafka"

type KafkaRequest struct {
	BrokerURLs      []string             `expr_enabled:"true" json:"brokerUrls"`
	Topic           string               `expr_enabled:"true" json:"topic"`
	Headers         []KafkaMessageHeader `json:"headers"`
	Authentication  *KafkaAuthenticator  `json:"authetication,omitempty"`
	MessageKey      string               `expr_enabled:"true" json:"messageKey"`
	MessageValue    string               `expr_enabled:"true" json:"messageValue"`
	SSLVerification bool                 `json:"sslVerification,omitempty"`
}

func (r KafkaRequest) GetHeaderAsMap() map[string]string {
	headerAsMap := make(map[string]string, len(r.Headers))

	for _, item := range r.Headers {
		headerAsMap[item.Key] = item.Value
	}

	return headerAsMap
}

type KafkaMessageHeader struct {
	Key   string `expr_enabled:"true" json:"key,omitempty"`
	Value string `expr_enabled:"true" json:"value,omitempty"`
}

type KafkaAuthenticator struct {
	Type  string                   `json:"type,omitempty" expr_enabled:"true"`
	Plain *KafkaPlainAuthenticator `json:"plain,omitempty"`
}

type KafkaPlainAuthenticator struct {
	Username string `json:"username,omitempty" expr_enabled:"true"`
	Password string `json:"password,omitempty" expr_enabled:"true"`
}

type KafkaResponse struct {
	Partition string
	Offset    string
}
