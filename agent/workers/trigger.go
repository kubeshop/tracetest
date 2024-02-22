package workers

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/collector"
	"github.com/kubeshop/tracetest/agent/event"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/sensors"
	agentTrigger "github.com/kubeshop/tracetest/agent/workers/trigger"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/test/trigger"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
	"go.uber.org/zap"
)

type TriggerWorker struct {
	logger                 *zap.Logger
	client                 *client.Client
	registry               *agentTrigger.Registry
	traceCache             collector.TraceCache
	observer               event.Observer
	sensor                 sensors.Sensor
	stoppableProcessRunner StoppableProcessRunner
	tracer                 trace.Tracer
	meter                  metric.Meter
}

type TriggerOption func(*TriggerWorker)

func WithTriggerStoppableProcessRunner(stoppableProcessRunner StoppableProcessRunner) TriggerOption {
	return func(tw *TriggerWorker) {
		tw.stoppableProcessRunner = stoppableProcessRunner
	}
}

func WithTraceCache(cache collector.TraceCache) TriggerOption {
	return func(tw *TriggerWorker) {
		tw.traceCache = cache
	}
}

func WithTriggerObserver(observer event.Observer) TriggerOption {
	return func(tw *TriggerWorker) {
		tw.observer = observer
	}
}

func WithSensor(sensor sensors.Sensor) TriggerOption {
	return func(tw *TriggerWorker) {
		tw.sensor = sensor
	}
}

func WithTriggerLogger(logger *zap.Logger) TriggerOption {
	return func(tw *TriggerWorker) {
		tw.logger = logger
	}
}

func WithTriggerTracer(tracer trace.Tracer) TriggerOption {
	return func(tw *TriggerWorker) {
		tw.tracer = tracer
	}
}

func WithTriggerMeter(meter metric.Meter) TriggerOption {
	return func(tw *TriggerWorker) {
		tw.meter = meter
	}
}

func NewTriggerWorker(client *client.Client, opts ...TriggerOption) *TriggerWorker {
	worker := &TriggerWorker{
		client:   client,
		logger:   zap.NewNop(),
		observer: event.NewNopObserver(),
		tracer:   trace.NewNoopTracerProvider().Tracer("noop"),
	}

	for _, opt := range opts {
		opt(worker)
	}

	registry := agentTrigger.NewRegistry(worker.tracer)
	registry.Add(agentTrigger.HTTP())
	registry.Add(agentTrigger.GRPC())
	registry.Add(agentTrigger.TRACEID())
	registry.Add(agentTrigger.KAFKA())

	// Assign registry into worker
	worker.registry = registry

	return worker
}

func (w *TriggerWorker) Trigger(ctx context.Context, triggerRequest *proto.TriggerRequest) error {
	ctx, span := w.tracer.Start(ctx, "TriggerRequest Worker operation")
	defer span.End()

	runCounter, _ := w.meter.Int64Counter("tracetest.agent.triggerworker.runs")
	runCounter.Add(ctx, 1)

	errorCounter, _ := w.meter.Int64Counter("tracetest.agent.triggerworker.errors")

	w.logger.Debug("Trigger request received", zap.Any("triggerRequest", triggerRequest))
	w.observer.StartTriggerExecution(triggerRequest)

	var err error
	w.stoppableProcessRunner(ctx, triggerRequest.TestID, triggerRequest.RunID, func(subcontext context.Context) {
		err = w.trigger(subcontext, triggerRequest)
	}, func(_ string) {
		err = executor.ErrUserCancelled
	})

	if err != nil {
		w.logger.Error("Trigger error", zap.Error(err))
		w.observer.EndTriggerExecution(triggerRequest, err)

		sendErr := w.client.SendTriggerResponse(ctx, &proto.TriggerResponse{
			RequestID:           triggerRequest.RequestID,
			AgentIdentification: w.client.SessionConfiguration().AgentIdentification,
			TestID:              triggerRequest.GetTestID(),
			RunID:               triggerRequest.GetRunID(),
			TriggerResult: &proto.TriggerResult{
				Error: &proto.Error{
					Message: err.Error(),
				},
			},
		})

		if sendErr != nil {
			w.logger.Error("Could not report trigger error back to the server", zap.Error(sendErr))
			w.observer.Error(sendErr)

			formattedErr := fmt.Errorf("could not report trigger error back to the server: %w. Original error: %s", sendErr, err.Error())
			span.RecordError(formattedErr)
			errorCounter.Add(ctx, 1)

			return formattedErr
		}

		span.RecordError(err)
		errorCounter.Add(ctx, 1)
	}

	w.observer.EndTriggerExecution(triggerRequest, err)

	return nil
}

func (w *TriggerWorker) trigger(ctx context.Context, triggerRequest *proto.TriggerRequest) error {
	triggerConfig := convertProtoToTrigger(triggerRequest.Trigger)
	w.logger.Debug("Triggering test", zap.Any("triggerConfig", triggerConfig))
	triggerer, err := w.registry.Get(triggerConfig.Type)
	if err != nil {
		w.logger.Error("Could not get triggerer", zap.Error(err))
		return err
	}
	w.logger.Debug("Triggerer found", zap.Any("triggerer", triggerer))

	traceID, err := trace.TraceIDFromHex(triggerRequest.TraceID)
	if err != nil {
		w.logger.Error("Invalid traceID was received in TriggerRequest", zap.Error(err))
		return fmt.Errorf("invalid traceID was received in TriggerRequest: %w", err)
	}
	w.logger.Debug("TraceID parsed", zap.Any("traceID", traceID))

	if w.traceCache != nil {
		// Set traceID to cache so the collector starts watching for incoming traces
		// with same id
		w.logger.Debug("Appending traceID to trace cache", zap.Any("traceID", traceID))
		w.traceCache.Append(triggerRequest.TraceID, []*v1.Span{})
	}

	response, err := triggerer.Trigger(ctx, triggerConfig, &agentTrigger.Options{
		TraceID: traceID,
		SpanID:  id.NewRandGenerator().SpanID(),
		TestID:  id.ID(triggerRequest.TestID),
	})
	if err != nil {
		w.logger.Error("Could not trigger test", zap.Error(err))
		return fmt.Errorf("could not trigger test: %w", err)
	}

	w.logger.Debug("Test triggered", zap.Any("response", response))

	protoResponse := convertResponseToProtoResponse(triggerRequest, response)
	protoResponse.RequestID = triggerRequest.RequestID
	w.logger.Debug("Sending trigger response to server", zap.Any("protoResponse", protoResponse))
	err = w.client.SendTriggerResponse(ctx, protoResponse)
	if err != nil {
		w.logger.Error("Could not send trigger response to server", zap.Error(err))
		return fmt.Errorf("could not send trigger response to server: %w", err)
	}
	w.logger.Debug("Trigger response sent to server")

	return nil
}

func convertProtoToTrigger(pt *proto.Trigger) trigger.Trigger {
	return trigger.Trigger{
		Type:    trigger.TriggerType(pt.Type),
		HTTP:    convertProtoHttpTriggerToHttpTrigger(pt.Http),
		GRPC:    convertProtoGrpcTriggerToGrpcTrigger(pt.Grpc),
		TraceID: convertProtoTraceIDTriggerToTraceIDTrigger(pt.TraceID),
		Kafka:   convertProtoKafkaTriggerToKafkaTrigger(pt.Kafka),
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

func convertProtoKafkaTriggerToKafkaTrigger(kafkaRequest *proto.KafkaRequest) *trigger.KafkaRequest {
	if kafkaRequest == nil {
		return nil
	}

	headers := make([]trigger.KafkaMessageHeader, len(kafkaRequest.Headers))

	for i, h := range headers {
		headers[i] = trigger.KafkaMessageHeader{Key: h.Key, Value: h.Value}
	}

	return &trigger.KafkaRequest{
		BrokerURLs:      kafkaRequest.BrokerUrls,
		Topic:           kafkaRequest.Topic,
		Headers:         headers,
		Authentication:  convertProtoKafkaAuthToKafkaAuth(kafkaRequest.Authentication),
		MessageKey:      kafkaRequest.MessageKey,
		MessageValue:    kafkaRequest.MessageValue,
		SSLVerification: kafkaRequest.SslVerification,
	}
}

func convertProtoKafkaAuthToKafkaAuth(kafkaAuthentication *proto.KafkaAuthentication) *trigger.KafkaAuthenticator {
	if kafkaAuthentication == nil {
		return nil
	}

	return &trigger.KafkaAuthenticator{
		Type: kafkaAuthentication.Type,
		Plain: &trigger.KafkaPlainAuthenticator{
			Username: kafkaAuthentication.Plain.Username,
			Password: kafkaAuthentication.Plain.Password,
		},
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
			Kafka:   convertKafkaResponseToProto(response.Result.Kafka),
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
		Body:       []byte(http.Body),
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
		Body:       []byte(grpc.Body),
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

func convertKafkaResponseToProto(kafka *trigger.KafkaResponse) *proto.KafkaResponse {
	if kafka == nil || kafka.Offset == "" {
		return nil
	}

	return &proto.KafkaResponse{
		Partition: kafka.Partition,
		Offset:    kafka.Offset,
	}
}
