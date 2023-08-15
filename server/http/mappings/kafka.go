package mappings

import (
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/test/trigger"
)

// out

func (m OpenAPI) KafkaRequest(in *trigger.KafkaRequest) openapi.KafkaRequest {
	if in == nil {
		return openapi.KafkaRequest{}
	}

	return openapi.KafkaRequest{
		BrokerUrls:      in.BrokerURLs,
		Authentication:  m.KafkaAuth(in.Authentication),
		Topic:           in.Topic,
		SslVerification: in.SSLVerification,
		Headers:         m.KafkaMessageHeaders(in.Headers),
		MessageKey:      in.MessageKey,
		MessageValue:    in.MessageValue,
	}
}

func (m OpenAPI) KafkaResponse(in *trigger.KafkaResponse) openapi.KafkaResponse {
	if in == nil {
		return openapi.KafkaResponse{}
	}
	return openapi.KafkaResponse{
		Partition: in.Partition,
		Offset:    in.Offset,
	}
}

func (m OpenAPI) KafkaMessageHeaders(in []trigger.KafkaMessageHeader) []openapi.KafkaMessageHeader {
	headers := make([]openapi.KafkaMessageHeader, len(in))
	for i, h := range in {
		headers[i] = openapi.KafkaMessageHeader{Key: h.Key, Value: h.Value}
	}

	return headers
}

func (m OpenAPI) KafkaAuth(in *trigger.KafkaAuthenticator) openapi.KafkaAuthentication {
	if in == nil {
		return openapi.KafkaAuthentication{}
	}

	auth := openapi.KafkaAuthentication{
		Type: in.Type,
	}
	switch in.Type {
	case "plain":
		auth.Plain = openapi.HttpAuthBasic{
			Username: in.Plain.Username,
			Password: in.Plain.Password,
		}
	}

	return auth
}

// in

func (m Model) KafkaRequest(in openapi.KafkaRequest) *trigger.KafkaRequest {
	// ignore unset kafka requests
	if in.BrokerUrls == nil || len(in.BrokerUrls) == 0 || in.BrokerUrls[0] == "" {
		return nil
	}

	return &trigger.KafkaRequest{
		BrokerURLs:      in.BrokerUrls,
		Authentication:  m.KafkaAuth(in.Authentication),
		Topic:           in.Topic,
		SSLVerification: in.SslVerification,
		Headers:         m.KafkaMessageHeaders(in.Headers),
		MessageKey:      in.MessageKey,
		MessageValue:    in.MessageValue,
	}
}

func (m Model) KafkaResponse(in openapi.KafkaResponse) *trigger.KafkaResponse {
	// ignore unset kafka responses
	if in.Offset == "" {
		return nil
	}

	return &trigger.KafkaResponse{
		Partition: in.Partition,
		Offset:    in.Offset,
	}
}

func (m Model) KafkaMessageHeaders(in []openapi.KafkaMessageHeader) []trigger.KafkaMessageHeader {
	headers := make([]trigger.KafkaMessageHeader, len(in))
	for i, h := range in {
		headers[i] = trigger.KafkaMessageHeader{Key: h.Key, Value: h.Value}
	}

	return headers
}

func (m Model) KafkaAuth(in openapi.KafkaAuthentication) *trigger.KafkaAuthenticator {
	return &trigger.KafkaAuthenticator{
		Type: in.Type,
		Plain: &trigger.KafkaPlainAuthenticator{
			Username: in.Plain.Username,
			Password: in.Plain.Password,
		},
	}
}
