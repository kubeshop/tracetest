package mappings

import (
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/test/trigger"
)

// out

func (m OpenAPI) GRPCRequest(in *trigger.GRPCRequest) openapi.GrpcRequest {
	if in == nil {
		return openapi.GrpcRequest{}
	}

	return openapi.GrpcRequest{
		ProtobufFile: in.ProtobufFile,
		Address:      in.Address,
		Service:      in.Service,
		Method:       in.Method,
		Metadata:     m.GRPCMetadata(in.Metadata),
		Auth:         m.Auth(in.Auth),
		Request:      in.Request,
	}
}

func (m OpenAPI) GRPCResponse(in *trigger.GRPCResponse) openapi.GrpcResponse {
	if in == nil {
		return openapi.GrpcResponse{}
	}
	return openapi.GrpcResponse{
		StatusCode: int32(in.StatusCode),
		Metadata:   m.GRPCMetadata(in.Metadata),
		Body:       in.Body,
	}
}

func (m OpenAPI) GRPCMetadata(in []trigger.GRPCHeader) []openapi.GrpcHeader {
	headers := make([]openapi.GrpcHeader, len(in))
	for i, h := range in {
		headers[i] = openapi.GrpcHeader{Key: h.Key, Value: h.Value}
	}

	return headers
}

//in

func (m Model) GRPCHeaders(in []openapi.GrpcHeader) []trigger.GRPCHeader {
	headers := make([]trigger.GRPCHeader, len(in))
	for i, h := range in {
		headers[i] = trigger.GRPCHeader{Key: h.Key, Value: h.Value}
	}

	return headers
}

func (m Model) GRPCRequest(in openapi.GrpcRequest) *trigger.GRPCRequest {
	// ignore unset grpc requests
	if in.Address == "" {
		return nil
	}

	return &trigger.GRPCRequest{
		ProtobufFile: in.ProtobufFile,
		Address:      in.Address,
		Method:       in.Method,
		Metadata:     m.GRPCHeaders(in.Metadata),
		Auth:         m.Auth(in.Auth),
		Request:      in.Request,
	}
}

func (m Model) GRPCResponse(in openapi.GrpcResponse) *trigger.GRPCResponse {
	// ignore unset grcp responses
	if in.StatusCode == 0 {
		return nil
	}

	return &trigger.GRPCResponse{
		StatusCode: int(in.StatusCode),
		Metadata:   m.GRPCHeaders(in.Metadata),
		Body:       in.Body,
	}
}
