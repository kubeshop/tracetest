package trigger

import (
	"context"
	"fmt"
	"github.com/kubeshop/tracetest/server/expression"
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
			TRACEID: &model.TRACEIDResponse{ID: test.ServiceUnderTest.TRACEID.ID},
		},
	}

	return response, nil
}

func (t *traceidTriggerer) Type() model.TriggerType {
	return model.TriggerTypeTRACEID
}

func (t *traceidTriggerer) Resolve(ctx context.Context, test model.Test, opts *TriggerOptions) (model.Test, error) {
	traceid := test.ServiceUnderTest.TRACEID
	if traceid == nil {
		return test, fmt.Errorf("no settings provided for TRACEID triggerer")
	}

	id, err := opts.Executor.ResolveStatement(WrapInQuotes(traceid.ID, "\""))
	if err != nil {
		return test, err
	}

	traceid.ID = id
	test.ServiceUnderTest.TRACEID = traceid

	return test, nil
}

func (t *traceidTriggerer) Variables(ctx context.Context, test model.Test, executor expression.Executor) (expression.VariablesMap, error) {
	triggerVariables := expression.VariablesMap{}

	traceid := test.ServiceUnderTest.TRACEID
	if traceid == nil {
		return triggerVariables, fmt.Errorf("no settings provided for TRACEID triggerer")
	}

	idVariables, err := executor.StatementTermsByType(WrapInQuotes(traceid.ID, "\""), expression.EnvironmentType)
	if err != nil {
		return triggerVariables, err
	}

	triggerVariables = triggerVariables.MergeStringArray(idVariables)

	return triggerVariables, nil
}
