package trigger

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/model"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func GRPC(config config.Config) (Triggerer, error) {
	tracerProvider, err := getTracerProvider(config)
	if err != nil {
		return nil, fmt.Errorf("could not create HTTP triggerer: %w", err)
	}

	return &grpcTriggerer{
		traceProvider: tracerProvider,
		tracer:        tracerProvider.Tracer("tracetest"),
	}, nil
}

type grpcTriggerer struct {
	traceProvider *sdktrace.TracerProvider
	tracer        trace.Tracer
}

func (te *grpcTriggerer) Trigger(_ context.Context, test model.Test) (Response, error) {

	response := Response{
		Result: model.TriggerResult{
			Type: te.Type(),
		},
	}

	trigger := test.ServiceUnderTest
	if trigger.Type != model.TriggerTypeGRPC {
		return response, fmt.Errorf(`trigger type "%s" not supported by GRPC triggerer`, trigger.Type)
	}

	if trigger.GRPC == nil {
		return response, fmt.Errorf("no settings provided for GRPC triggerer")
	}

	ctx, span := te.tracer.Start(context.Background(), "Tracetest Trigger")
	defer span.End()

	var tf trace.TraceFlags
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    span.SpanContext().TraceID(),
		SpanID:     span.SpanContext().SpanID(),
		TraceFlags: tf.WithSampled(true),
		TraceState: trace.TraceState{},
		Remote:     true,
	})

	ctx = trace.ContextWithSpanContext(context.Background(), sc)

	defer span.End()

	tReq := trigger.GRPC

	conn, err := te.dial(ctx, tReq.Address)
	if err != nil {
		return response, fmt.Errorf("cannot dial service: %w", err)
	}

	desc, err := protoDescription(tReq.ProtobufFile)
	if err != nil {
		return response, fmt.Errorf("cannot read descriptors: %w", err)
	}

	options := grpcurl.FormatOptions{
		EmitJSONDefaultFields: true,
		IncludeTextSeparator:  true,
		AllowUnknownFields:    true,
	}

	rf, _, err := grpcurl.RequestParserAndFormatter(grpcurl.FormatJSON, desc, strings.NewReader(tReq.Request), options)
	if err != nil {
		return response, fmt.Errorf("failed to construct request parser and formatter for %w", err)
	}

	anyResolver := grpcurl.AnyResolverFromDescriptorSource(desc)
	marshaler := jsonpb.Marshaler{
		EmitDefaults: true,
		Indent:       "  ",
		AnyResolver:  anyResolver,
	}

	h := &eventHandler{
		marshaller: marshaler,
	}

	err = grpcurl.InvokeRPC(ctx, desc, conn, tReq.Method, tReq.Headers(), h, rf.Next)
	if err != nil {
		return response, err
	}

	response.Result.GRPC = &model.GRPCResponse{
		Metadata:   mapHeaders(h.respMD),
		StatusCode: int(h.respCode),
		Status:     h.respCode.String(),
		Body:       h.respBody,
	}

	response.SpanAttributes = map[string]string{
		"tracetest.run.trigger.grpc.response_status_code": strconv.Itoa(int(h.respCode)),
		"tracetest.run.trigger.grpc.response_status":      h.respCode.String(),
	}

	return response, nil
}

func (t *grpcTriggerer) Type() model.TriggerType {
	return model.TriggerTypeGRPC
}

func protoDescription(content string) (grpcurl.DescriptorSource, error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "protofile")
	if err != nil {
		return nil, fmt.Errorf("cannot create tmp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err = tmpFile.Write([]byte(content)); err != nil {
		return nil, fmt.Errorf("cannot write tmp file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return nil, fmt.Errorf("cannot close tmp file: %w", err)
	}

	desc, err := grpcurl.DescriptorSourceFromProtoFiles([]string{os.TempDir()}, tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("cannot parse proto file: %w", err)
	}

	return desc, nil

}

func mapHeaders(md metadata.MD) []model.GRPCHeader {
	var mappedHeaders []model.GRPCHeader
	for key, headers := range md {
		for _, val := range headers {
			val := model.GRPCHeader{
				Key:   key,
				Value: val,
			}
			mappedHeaders = append(mappedHeaders, val)
		}
	}

	return mappedHeaders
}

func (t *grpcTriggerer) dial(ctx context.Context, address string) (*grpc.ClientConn, error) {
	var creds credentials.TransportCredentials
	network := "tcp"

	return grpcurl.BlockingDial(
		ctx, network, address, creds,
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor(
			otelgrpc.WithTracerProvider(noopTracerProvider()),
			otelgrpc.WithPropagators(propagators()),
		)),
	)
}

var _ grpcurl.InvocationEventHandler = (*eventHandler)(nil)

type eventHandler struct {
	marshaller jsonpb.Marshaler
	respBody   string
	respCode   codes.Code
	respMD     metadata.MD
}

func (h *eventHandler) OnResolveMethod(md *desc.MethodDescriptor) {}

func (h *eventHandler) OnSendHeaders(md metadata.MD) {
}

func (h *eventHandler) OnReceiveHeaders(md metadata.MD) {
	h.respMD = md
}

func (h *eventHandler) OnReceiveResponse(resp proto.Message) {
	j, err := h.marshaller.MarshalToString(resp)
	if err != nil {
		panic(err)
	}

	h.respBody = j
}

func (h *eventHandler) OnReceiveTrailers(stat *status.Status, md metadata.MD) {
	h.respCode = stat.Code()
}
