package trigger

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/model"
)

func TRACEID() Triggerer {
	return &traceidTriggerer{}
}

type traceidTriggerer struct{}

func (t *traceidTriggerer) Trigger(ctx context.Context, test model.Test, opts *TriggerOptions) (Response, error) {
	response := Response{
		Result: model.TriggerResult{
			Type:    t.Type(),
			TRACEID: &model.TRACEIDResponse{ID: test.ServiceUnderTest.TraceID.ID},
		},
	}

	return response, nil
}

func (t *traceidTriggerer) Type() model.TriggerType {
	return model.TriggerTypeTRACEID
}

func (t *traceidTriggerer) Resolve(ctx context.Context, test model.Test, opts *TriggerOptions) (model.Test, error) {
	traceid := test.ServiceUnderTest.TraceID
	if traceid == nil {
		return test, fmt.Errorf("no settings provided for TRACEID triggerer")
	}

	id, err := opts.Executor.ResolveStatement(WrapInQuotes(traceid.ID, "\""))
	if err != nil {
		return test, err
	}

	traceid.ID = id
	test.ServiceUnderTest.TraceID = traceid

	return test, nil
}
