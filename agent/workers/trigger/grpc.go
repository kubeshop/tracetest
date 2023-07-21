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
	"github.com/kubeshop/tracetest/server/test/trigger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func GRPC() Triggerer {
	return &grpcTriggerer{}
}

type grpcTriggerer struct{}

func (te *grpcTriggerer) Trigger(ctx context.Context, triggerConfig trigger.Trigger, opts *Options) (Response, error) {
	response := Response{
		Result: trigger.TriggerResult{
			Type: te.Type(),
		},
	}

	if triggerConfig.Type != trigger.TriggerTypeGRPC {
		return response, fmt.Errorf(`trigger type "%s" not supported by GRPC triggerer`, triggerConfig.Type)
	}

	if triggerConfig.GRPC == nil {
		return response, fmt.Errorf("no settings provided for GRPC triggerer")
	}

	tReq := triggerConfig.GRPC

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

	md := tReq.MD()
	otelgrpc.Inject(ctx, md, otelgrpc.WithPropagators(propagators()))

	err = grpcurl.InvokeRPC(ctx, desc, conn, tReq.Method, mdToHeaders(md), h, rf.Next)
	if err != nil {
		return response, err
	}

	response.Result.GRPC = &trigger.GRPCResponse{
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

func (t *grpcTriggerer) Type() trigger.TriggerType {
	return trigger.TriggerTypeGRPC
}

func mdToHeaders(md *metadata.MD) []string {
	h := []string{}

	for k, vs := range *md {
		h = append(h, k+": "+strings.Join(vs, " "))
	}

	return h
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

func mapHeaders(md metadata.MD) []trigger.GRPCHeader {
	var mappedHeaders []trigger.GRPCHeader
	for key, headers := range md {
		for _, val := range headers {
			val := trigger.GRPCHeader{
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
