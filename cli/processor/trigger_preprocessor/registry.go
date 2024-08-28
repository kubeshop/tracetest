package trigger_preprocessor

import (
	"fmt"

	"github.com/kubeshop/tracetest/agent/workers/trigger"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"go.uber.org/zap"
)

type TriggerPreprocessor interface {
	Preprocess(input fileutil.File, test openapi.TestResource) (openapi.TestResource, error)
	Type() trigger.TriggerType
}

type Registry struct {
	processors map[trigger.TriggerType]TriggerPreprocessor
}

func NewRegistry(logger *zap.Logger) Registry {
	return Registry{
		processors: map[trigger.TriggerType]TriggerPreprocessor{},
	}
}

func (r Registry) Register(processor TriggerPreprocessor) Registry {
	r.processors[processor.Type()] = processor
	return r
}

var ErrNotFound = fmt.Errorf("preprocessor not found")

func (r Registry) Preprocess(input fileutil.File, test openapi.TestResource) (openapi.TestResource, error) {
	triggerType := test.Spec.Trigger.GetType()

	processor, ok := r.processors[trigger.TriggerType(triggerType)]
	if ok {
		return processor.Preprocess(input, test)
	}

	return test, nil
}
