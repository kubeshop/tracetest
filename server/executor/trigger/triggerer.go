package trigger

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/kubeshop/tracetest/agent/workers/trigger"
	"github.com/kubeshop/tracetest/server/expression"
	"github.com/kubeshop/tracetest/server/test"
	"go.opentelemetry.io/otel/trace"
)

type TriggerOptions struct {
	TraceID trace.TraceID
}

type ResolveOptions struct {
	Executor expression.Executor
}

type Triggerer interface {
	Trigger(context.Context, test.Test, *TriggerOptions) (Response, error)
	Type() trigger.TriggerType
	Resolve(context.Context, test.Test, *ResolveOptions) (test.Test, error)
}

type Response struct {
	SpanAttributes map[string]string
	Result         trigger.TriggerResult
	TraceID        trace.TraceID
	SpanID         trace.SpanID
}

func NewRegistry(tracer, triggerSpanTracer trace.Tracer) *Registry {
	return &Registry{
		tracer:            tracer,
		triggerSpanTracer: triggerSpanTracer,
		reg:               map[trigger.TriggerType]Triggerer{},
	}
}

type Registry struct {
	sync.Mutex
	tracer            trace.Tracer
	triggerSpanTracer trace.Tracer
	reg               map[trigger.TriggerType]Triggerer
}

func (r *Registry) Add(t Triggerer) {
	r.Lock()
	defer r.Unlock()

	r.reg[t.Type()] = t
}

var ErrTriggererTypeNotRegistered = errors.New("triggerer type not found")

func (r *Registry) Get(triggererType trigger.TriggerType) (Triggerer, error) {
	r.Lock()
	defer r.Unlock()

	if triggererType.IsTraceIDBased() {
		triggererType = trigger.TriggerTypeTraceID
	}

	t, found := r.reg[triggererType]
	if !found {
		return nil, fmt.Errorf(`cannot get trigger type "%s": %w`, triggererType, ErrTriggererTypeNotRegistered)
	}

	return Instrument(r.tracer, r.triggerSpanTracer, t), nil
}

func WrapInQuotes(input string, quoteChar string) string {
	input = strings.ReplaceAll(input, quoteChar, fmt.Sprintf("\\%s", quoteChar))
	return fmt.Sprintf("%s%s%s", quoteChar, input, quoteChar)
}
