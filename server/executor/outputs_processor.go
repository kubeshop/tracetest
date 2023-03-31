package executor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kubeshop/tracetest/server/assertions/selectors"
	"github.com/kubeshop/tracetest/server/expression"
	"github.com/kubeshop/tracetest/server/model"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type OutputsProcessorFn func(context.Context, model.OrderedMap[string, model.Output], model.Trace, []expression.DataStore) (model.OrderedMap[string, model.RunOutput], error)

func InstrumentedOutputProcessor(tracer trace.Tracer) OutputsProcessorFn {
	op := instrumentedOutputProcessor{tracer}
	return op.process
}

type instrumentedOutputProcessor struct {
	tracer trace.Tracer
}

func (op instrumentedOutputProcessor) process(ctx context.Context, outputs model.OrderedMap[string, model.Output], t model.Trace, ds []expression.DataStore) (model.OrderedMap[string, model.RunOutput], error) {
	ctx, span := op.tracer.Start(ctx, "Process outputs")
	defer span.End()

	result, err := outputProcessor(ctx, outputs, t, ds)
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

func outputProcessor(ctx context.Context, outputs model.OrderedMap[string, model.Output], tr model.Trace, ds []expression.DataStore) (model.OrderedMap[string, model.RunOutput], error) {
	res := model.OrderedMap[string, model.RunOutput]{}

	parsed, err := parseOutputs(outputs)
	if err != nil {
		return res, err
	}

	err = parsed.ForEach(func(key string, out parsedOutput) error {
		spans := out.selector.Filter(tr)

		stores := append([]expression.DataStore{expression.MetaAttributesDataStore{SelectedSpans: spans}}, ds...)

		value := ""
		spanId := ""
		resolved := false
		spans.
			ForEach(func(_ int, span model.Span) bool {
				value = extractAttr(span, stores, out.expr)
				spanId = span.ID.String()
				resolved = true
				// take only the first value
				return false
			}).
			OrEmpty(func() {
				value = extractAttr(model.Span{}, stores, out.expr)
				resolved = false
			})

		res, err = res.Add(key, model.RunOutput{
			Value:    value,
			SpanID:   spanId,
			Name:     key,
			Resolved: resolved,
		})
		if err != nil {
			return fmt.Errorf(`cannot process output "%s": %w`, key, err)
		}

		return nil
	})

	if err != nil {
		return model.OrderedMap[string, model.RunOutput]{}, err
	}

	return res, nil
}

func extractAttr(span model.Span, ds []expression.DataStore, expr expression.Expr) string {
	ds = append([]expression.DataStore{expression.AttributeDataStore{Span: span}}, ds...)

	expressionExecutor := expression.NewExecutor(ds...)

	actualValue, _ := expressionExecutor.ResolveExpression(&expr)

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
