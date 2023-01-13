package mappings

import (
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
)

func (m OpenAPI) TRACEIDRequest(in *model.TRACEIDRequest) openapi.TriggerTriggerSettingsTraceid {
	if in == nil {
		return openapi.TriggerTriggerSettingsTraceid{}
	}

	return openapi.TriggerTriggerSettingsTraceid{
		Id: in.ID,
	}
}

func (m Model) TRACEIDRequest(in openapi.TriggerTriggerSettingsTraceid) *model.TRACEIDRequest {
	if in.Id == "" {
		return nil
	}

	return &model.TRACEIDRequest{
		ID: in.Id,
	}
}
