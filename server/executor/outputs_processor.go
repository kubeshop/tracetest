package executor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kubeshop/tracetest/server/assertions/selectors"
	"github.com/kubeshop/tracetest/server/expression"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/kubeshop/tracetest/server/test"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type OutputsProcessorFn func(context.Context, test.Outputs, model.Trace, []expression.DataStore) (maps.Ordered[string, test.RunOutput], error)

func InstrumentedOutputProcessor(tracer trace.Tracer) OutputsProcessorFn {
	op := instrumentedOutputProcessor{tracer}
	return op.process
}

type instrumentedOutputProcessor struct {
	tracer trace.Tracer
}

func (op instrumentedOutputProcessor) process(ctx context.Context, outputs test.Outputs, t model.Trace, ds []expression.DataStore) (maps.Ordered[string, test.RunOutput], error) {
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

func outputProcessor(ctx context.Context, outputs test.Outputs, tr model.Trace, ds []expression.DataStore) (maps.Ordered[string, test.RunOutput], error) {
	res := maps.Ordered[string, test.RunOutput]{}

	parsed, err := parseOutputs(outputs)
	if err != nil {
		return res, err
	}

	err = parsed.ForEach(func(key string, out parsedOutput) error {
		if out.err != nil {
			res, err = res.Add(key, test.RunOutput{
				Value:    "",
				SpanID:   "",
				Name:     key,
				Resolved: false,
				Error:    out.err,
			})
			if err != nil {
				return fmt.Errorf(`cannot process output "%s": %w`, key, err)
			}

			return nil
		}

		spans := out.selector.Filter(tr)

		stores := append([]expression.DataStore{expression.MetaAttributesDataStore{SelectedSpans: spans}}, ds...)

		value := ""
		spanId := ""
		resolved := false
		var outputError error = nil
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
				outputError = fmt.Errorf(`cannot find matching spans for output "%s"`, key)
			})

		res, err = res.Add(key, test.RunOutput{
			Value:    value,
			SpanID:   spanId,
			Name:     key,
			Resolved: resolved,
			Error:    outputError,
		})
		if err != nil {
			return fmt.Errorf(`cannot process output "%s": %w`, key, err)
		}

		return nil
	})

	if err != nil {
		return maps.Ordered[string, test.RunOutput]{}, err
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
	err      error
}

func parseOutputs(outputs test.Outputs) (maps.Ordered[string, parsedOutput], error) {
	var parsed maps.Ordered[string, parsedOutput]

	for _, output := range outputs {
		key := output.Name
		out := output
		var selector selectors.Selector
		var expr expression.Expr
		var outputErr error

		expr, err := expression.Parse(out.Value)
		if err != nil {
			outputErr = fmt.Errorf(`cannot parse output "%s" value "%s": %w`, key, out.Value, err)
		}

		selector, err = selectors.New(string(out.Selector))
		if err != nil {
			outputErr = fmt.Errorf(`cannot parse output "%s" selector "%s": %w`, key, string(out.Selector), err)
		}

		parsed, _ = parsed.Add(key, parsedOutput{
			selector: selector,
			expr:     expr,
			err:      outputErr,
		})
	}

	return parsed, nil
}
