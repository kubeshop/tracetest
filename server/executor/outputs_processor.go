package executor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kubeshop/tracetest/server/assertions/selectors"
	"github.com/kubeshop/tracetest/server/expression"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type OutputsProcessorFn func(context.Context, model.OrderedMap[string, model.Output], traces.Trace) (model.OrderedMap[string, string], error)

func InstrumentedOutputProcessor(tracer trace.Tracer) OutputsProcessorFn {
	op := instrumentedOutputProcessor{tracer}
	return op.process
}

type instrumentedOutputProcessor struct {
	tracer trace.Tracer
}

func (op instrumentedOutputProcessor) process(ctx context.Context, outputs model.OrderedMap[string, model.Output], t traces.Trace) (model.OrderedMap[string, string], error) {
	ctx, span := op.tracer.Start(ctx, "Process outputs")
	defer span.End()

	result, err := outputProcessor(ctx, outputs, t)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return result, err
	}

	encoded, err := json.Marshal(result.Unordered())
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return result, err
	}

	span.SetAttributes(
		attribute.String("tracetest.run.output_processor.outputs", string(encoded)),
	)

	return result, err
}

func outputProcessor(ctx context.Context, outputs model.OrderedMap[string, model.Output], tr traces.Trace) (model.OrderedMap[string, string], error) {
	res := model.OrderedMap[string, string]{}

	parsed, err := parseOutputs(outputs)
	if err != nil {
		return model.OrderedMap[string, string]{}, err
	}

	err = parsed.ForEach(func(key string, out parsedOutput) error {
		spans := out.selector.Filter(tr)

		mads := expression.MetaAttributesDataStore{SelectedSpans: spans}
		value := ""
		spans.
			ForEach(func(_ int, span traces.Span) bool {
				value = extractAttr(traces.Span{}, mads, out.expr)
				// take only the first value
				return false
			}).
			OrEmpty(func() {
				value = extractAttr(traces.Span{}, mads, out.expr)
			})

		res, err = res.Add(key, value)
		if err != nil {
			return fmt.Errorf(`cannot process output "%s": %w`, key, err)
		}

		return nil
	})

	if err != nil {
		return model.OrderedMap[string, string]{}, err
	}

	return res, nil

}

func extractAttr(span traces.Span, mads expression.MetaAttributesDataStore, expr expression.Expr) string {
	attributeDataStore := expression.AttributeDataStore{Span: span}
	expressionExecutor := expression.NewExecutor(attributeDataStore, mads)

	actualValue, _ := expressionExecutor.Expression(expr)

	return actualValue.String()
}

type parsedOutput struct {
	selector selectors.Selector
	expr     expression.Expr
}

func parseOutputs(outputs model.OrderedMap[string, model.Output]) (model.OrderedMap[string, parsedOutput], error) {
	var parsed model.OrderedMap[string, parsedOutput]

	parseErr := outputs.ForEach(func(key string, out model.Output) error {
		expr, err := expression.Parse(out.Value)
		if err != nil {
			return fmt.Errorf(`cannot parse output "%s" value "%s": %w`, key, out.Value, err)
		}

		selector, err := selectors.New(string(out.Selector))
		if err != nil {
			return fmt.Errorf(`cannot parse output "%s" selector "%s": %w`, key, string(out.Selector), err)
		}

		parsed, _ = parsed.Add(key, parsedOutput{
			selector: selector,
			expr:     expr,
		})
		return nil
	})

	if parseErr != nil {
		return model.OrderedMap[string, parsedOutput]{}, parseErr
	}

	return parsed, nil
}
