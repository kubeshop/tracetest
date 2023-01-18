package mappings

import (
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
)

func (m OpenAPI) TRACEIDRequest(in *model.TRACEIDRequest) openapi.TraceidRequest {
	if in == nil {
		return openapi.TraceidRequest{}
	}

	return openapi.TraceidRequest{
		Id: in.ID,
	}
}

func (m OpenAPI) TRACEIDResponse(in *model.TRACEIDResponse) openapi.TraceidResponse {
	if in == nil {
		return openapi.TraceidResponse{}
	}
	return openapi.TraceidResponse{
		Id: in.ID,
	}
}

func (m Model) TRACEIDRequest(in openapi.TraceidRequest) *model.TRACEIDRequest {
	if in.Id == "" {
		return nil
	}

	return &model.TRACEIDRequest{
		ID: in.Id,
	}
}

func (m Model) TRACEIDResponse(in openapi.TraceidResponse) *model.TRACEIDResponse {
	if in.Id == "" {
		return nil
	}

	return &model.TRACEIDResponse{
		ID: in.Id,
	}
}
