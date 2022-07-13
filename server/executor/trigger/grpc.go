package trigger

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/kubeshop/tracetest/server/model"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func GRPC() Triggerer {
	return &grpcTriggerer{
		traceProvider: traceProvider(),
	}
}

type grpcTriggerer struct {
	traceProvider *sdktrace.TracerProvider
}

func (te *grpcTriggerer) Trigger(_ context.Context, test model.Test, tid trace.TraceID, sid trace.SpanID) (Response, error) {

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
	var tf trace.TraceFlags
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: tf.WithSampled(true),
		TraceState: trace.TraceState{},
		Remote:     true,
	})
	ctx := trace.ContextWithSpanContext(context.Background(), sc)

	tReq := trigger.GRPC
	fmt.Printf("**** grpc %+v\n", tReq)

	body := strings.NewReader(tReq.Request)

	conn, err := dial(ctx, tReq.Address)
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

	rf, _, err := grpcurl.RequestParserAndFormatter(grpcurl.FormatJSON, desc, body, options)
	if err != nil {
		return response, fmt.Errorf("failed to construct request parser and formatter for %q", err)
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

	err = grpcurl.InvokeRPC(ctx, desc, conn, tReq.Method, []string{}, h, rf.Next)
	if err != nil {
		return response, err
	}

	// mapped := mapResp(resp)
	// response.Result.GRPC = &mapped
	// response.SpanAttributes = map[string]string{
	// 	"tracetest.run.trigger.grpc.response_status": strconv.Itoa(resp.StatusCode),
	// }

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

// func mapResp(resp *grpc.Response) model.GRPCResponse {
// 	var mappedHeaders []model.GRPCHeader
// 	for key, headers := range resp.Header {
// 		for _, val := range headers {
// 			val := model.GRPCHeader{
// 				Key:   key,
// 				Value: val,
// 			}
// 			mappedHeaders = append(mappedHeaders, val)
// 		}
// 	}
// 	var body string
// 	if b, err := io.ReadAll(resp.Body); err == nil {
// 		body = string(b)
// 	} else {
// 		fmt.Println(err)
// 	}

// 	return model.GRPCResponse{
// 		// Status:     resp.Status,
// 		StatusCode: resp.StatusCode,
// 		// Headers:    mappedHeaders,
// 		Body: body,
// 	}
// }

func dial(ctx context.Context, address string) (*grpc.ClientConn, error) {
	var creds credentials.TransportCredentials
	network := "tcp"

	return grpcurl.BlockingDial(ctx, network, address, creds)
}

var _ grpcurl.InvocationEventHandler = (*eventHandler)(nil)

type eventHandler struct {
	marshaller jsonpb.Marshaler
}

func (h *eventHandler) OnResolveMethod(md *desc.MethodDescriptor) {
}

func (h *eventHandler) OnSendHeaders(md metadata.MD) {

}

func (h *eventHandler) OnReceiveHeaders(md metadata.MD) {

}

func (h *eventHandler) OnReceiveResponse(resp proto.Message) {
	j, err := h.marshaller.MarshalToString(resp)
	if err != nil {
		panic(err)
	}

	fmt.Println(j)
}

func (h *eventHandler) OnReceiveTrailers(stat *status.Status, md metadata.MD) {

}
