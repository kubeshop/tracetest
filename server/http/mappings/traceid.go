package mappings

import (
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/test/trigger"
)

func (m OpenAPI) TraceIDRequest(in *trigger.TraceIDRequest) openapi.TraceidRequest {
	if in == nil {
		return openapi.TraceidRequest{}
	}

	return openapi.TraceidRequest{
		Id: in.ID,
	}
}

func (m OpenAPI) TraceIDResponse(in *trigger.TraceIDResponse) openapi.TraceidResponse {
	if in == nil {
		return openapi.TraceidResponse{}
	}
	return openapi.TraceidResponse{
		Id: in.ID,
	}
}

func (m Model) TraceIDRequest(in openapi.TraceidRequest) *trigger.TraceIDRequest {
	if in.Id == "" {
		return nil
	}

	return &trigger.TraceIDRequest{
		ID: in.Id,
	}
}

func (m Model) TraceIDResponse(in openapi.TraceidResponse) *trigger.TraceIDResponse {
	if in.Id == "" {
		return nil
	}

	return &trigger.TraceIDResponse{
		ID: in.Id,
	}
}
