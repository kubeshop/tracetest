package workers

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/proto"
	agentTrigger "github.com/kubeshop/tracetest/agent/workers/trigger"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/test/trigger"
	"go.opentelemetry.io/otel/trace"
)

type TriggerWorker struct {
	client   *client.Client
	registry *agentTrigger.Registry
}

func NewTriggerWorker(client *client.Client) *TriggerWorker {
	// TODO: use a real tracer
	tracer := trace.NewNoopTracerProvider().Tracer("noop")

	registry := agentTrigger.NewRegistry(tracer)
	registry.Add(agentTrigger.HTTP())
	registry.Add(agentTrigger.GRPC())
	registry.Add(agentTrigger.TRACEID())

	return &TriggerWorker{client, registry}
}

func (w *TriggerWorker) Trigger(ctx context.Context, triggerRequest *proto.TriggerRequest) error {
	fmt.Println("Trigger handled by agent")
	triggerConfig := convertProtoToTrigger(triggerRequest.Trigger)
	triggerer, err := w.registry.Get(triggerConfig.Type)
	if err != nil {
		return err
	}

	traceID, err := trace.TraceIDFromHex(triggerRequest.TraceID)
	if err != nil {
		return fmt.Errorf("invalid traceID was received in TriggerRequest: %w", err)
	}

	response, err := triggerer.Trigger(ctx, triggerConfig, &agentTrigger.Options{
		TraceID: traceID,
		SpanID:  id.NewRandGenerator().SpanID(),
		TestID:  id.ID(triggerRequest.TestID),
	})

	if err != nil {
		return fmt.Errorf("could not trigger test: %w", err)
	}

	protoResponse := convertResponseToProtoResponse(triggerRequest, response)
	protoResponse.RequestID = triggerRequest.RequestID
	err = w.client.SendTriggerResponse(ctx, protoResponse)
	if err != nil {
		return fmt.Errorf("could not send trigger response to server: %w", err)
	}

	return nil
}

func convertProtoToTrigger(pt *proto.Trigger) trigger.Trigger {
	return trigger.Trigger{
		Type:    trigger.TriggerType(pt.Type),
		HTTP:    convertProtoHttpTriggerToHttpTrigger(pt.Http),
		GRPC:    convertProtoGrpcTriggerToGrpcTrigger(pt.Grpc),
		TraceID: convertProtoTraceIDTriggerToTraceIDTrigger(pt.TraceID),
	}
}

func convertProtoHttpTriggerToHttpTrigger(pt *proto.HttpRequest) *trigger.HTTPRequest {
	if pt == nil {
		return nil
	}

	headers := make([]trigger.HTTPHeader, 0, len(pt.Headers))

	for _, header := range pt.Headers {
		headers = append(headers, trigger.HTTPHeader{Key: header.Key, Value: header.Value})
	}

	return &trigger.HTTPRequest{
		Method:          trigger.HTTPMethod(pt.Method),
		URL:             pt.Url,
		Body:            pt.Body,
		Headers:         headers,
		Auth:            convertProtoHttpAuthToHttpAuth(pt.Authentication),
		SSLVerification: pt.SSLVerification,
	}
}

func convertProtoHttpAuthToHttpAuth(httpAuthentication *proto.HttpAuthentication) *trigger.HTTPAuthenticator {
	if httpAuthentication == nil {
		return nil
	}

	var (
		apiKey *trigger.APIKeyAuthenticator
		basic  *trigger.BasicAuthenticator
		bearer *trigger.BearerAuthenticator
	)

	if httpAuthentication.ApiKey != nil {
		apiKey = &trigger.APIKeyAuthenticator{
			Key:   httpAuthentication.ApiKey.Key,
			Value: httpAuthentication.ApiKey.Value,
			In:    trigger.APIKeyPosition(httpAuthentication.ApiKey.In),
		}
	}

	if httpAuthentication.Basic != nil {
		basic = &trigger.BasicAuthenticator{
			Username: httpAuthentication.Basic.Username,
			Password: httpAuthentication.Basic.Password,
		}
	}

	if httpAuthentication.Bearer != nil {
		bearer = &trigger.BearerAuthenticator{
			Bearer: httpAuthentication.Bearer.Token,
		}
	}

	return &trigger.HTTPAuthenticator{
		Type:   httpAuthentication.Type,
		APIKey: apiKey,
		Basic:  basic,
		Bearer: bearer,
	}
}

func convertProtoGrpcTriggerToGrpcTrigger(grpcRequest *proto.GrpcRequest) *trigger.GRPCRequest {
	if grpcRequest == nil {
		return nil
	}

	metadata := make([]trigger.GRPCHeader, 0, len(grpcRequest.Metadata))
	for _, keyValuePair := range grpcRequest.Metadata {
		metadata = append(metadata, trigger.GRPCHeader{Key: keyValuePair.Key, Value: keyValuePair.Value})
	}

	return &trigger.GRPCRequest{
		ProtobufFile: grpcRequest.ProtobufFile,
		Address:      grpcRequest.Address,
		Service:      grpcRequest.Service,
		Method:       grpcRequest.Method,
		Request:      grpcRequest.Request,
		Metadata:     metadata,
		Auth:         convertProtoHttpAuthToHttpAuth(grpcRequest.Authentication),
	}
}

func convertProtoTraceIDTriggerToTraceIDTrigger(traceIDRequest *proto.TraceIDRequest) *trigger.TraceIDRequest {
	if traceIDRequest == nil {
		return nil
	}

	return &trigger.TraceIDRequest{
		ID: traceIDRequest.Id,
	}
}

func convertResponseToProtoResponse(request *proto.TriggerRequest, response agentTrigger.Response) *proto.TriggerResponse {
	return &proto.TriggerResponse{
		TestID: request.TestID,
		RunID:  request.RunID,
		TriggerResult: &proto.TriggerResult{
			Type:    string(response.Result.Type),
			Http:    convertHttpResponseToProto(response.Result.HTTP),
			Grpc:    convertGrpcResponseToProto(response.Result.GRPC),
			TraceID: convertTraceIDResponseToProto(response.Result.TraceID),
		},
	}
}

func convertHttpResponseToProto(http *trigger.HTTPResponse) *proto.HttpResponse {
	if http == nil {
		return nil
	}

	headers := make([]*proto.HttpHeader, 0, len(http.Headers))
	for _, header := range http.Headers {
		headers = append(headers, &proto.HttpHeader{Key: header.Key, Value: header.Value})
	}

	return &proto.HttpResponse{
		StatusCode: int32(http.StatusCode),
		Status:     http.Status,
		Headers:    headers,
		Body:       http.Body,
	}
}

func convertGrpcResponseToProto(grpc *trigger.GRPCResponse) *proto.GrpcResponse {
	if grpc == nil {
		return nil
	}

	headers := make([]*proto.GrpcHeader, 0, len(grpc.Metadata))
	for _, header := range grpc.Metadata {
		headers = append(headers, &proto.GrpcHeader{Key: header.Key, Value: header.Value})
	}

	return &proto.GrpcResponse{
		StatusCode: int32(grpc.StatusCode),
		Metadata:   headers,
		Body:       grpc.Body,
	}
}

func convertTraceIDResponseToProto(traceID *trigger.TraceIDResponse) *proto.TraceIdResponse {
	if traceID == nil {
		return nil
	}

	return &proto.TraceIdResponse{
		Id: traceID.ID,
	}
}
