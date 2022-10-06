package mappings

import (
	"fmt"
	"strconv"

	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/assertions/selectors"
	"github.com/kubeshop/tracetest/server/encoding/yaml/conversion"
	"github.com/kubeshop/tracetest/server/encoding/yaml/conversion/parser"
	"github.com/kubeshop/tracetest/server/encoding/yaml/definition"
	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/trace"
)

// out

type OpenAPI struct {
	traceConversionConfig traces.ConversionConfig
}

func (m OpenAPI) TestDefinitionFile(in model.Test) definition.Test {
	testDefinition, _ := conversion.ConvertOpenAPITestIntoDefinitionObject(m.Test(in))
	return testDefinition
}

func (m OpenAPI) Test(in model.Test) openapi.Test {
	return openapi.Test{
		Id:               string(in.ID),
		Name:             in.Name,
		Description:      in.Description,
		ServiceUnderTest: m.Trigger(in.ServiceUnderTest),
		Specs:            m.Specs(in.Specs),
		Version:          int32(in.Version),
		Summary: openapi.TestSummary{
			Runs: int32(in.Summary.Runs),
			LastRun: openapi.TestSummaryLastRun{
				Time:   in.Summary.LastRun.Time,
				Passes: int32(in.Summary.LastRun.Passes),
				Fails:  int32(in.Summary.LastRun.Fails),
			},
		},
	}
}

func (m OpenAPI) Trigger(in model.Trigger) openapi.Trigger {
	return openapi.Trigger{
		TriggerType: string(in.Type),
		TriggerSettings: openapi.TriggerTriggerSettings{
			Http: m.HTTPRequest(in.HTTP),
			Grpc: m.GRPCRequest(in.GRPC),
		},
	}
}

func (m OpenAPI) TriggerResult(in model.TriggerResult) openapi.TriggerResult {

	return openapi.TriggerResult{
		TriggerType: string(in.Type),
		TriggerResult: openapi.TriggerResultTriggerResult{
			Http: m.HTTPResponse(in.HTTP),
			Grpc: m.GRPCResponse(in.GRPC),
		},
	}
}

func (m OpenAPI) Tests(in []model.Test) []openapi.Test {
	tests := make([]openapi.Test, len(in))
	for i, t := range in {
		tests[i] = m.Test(t)
	}

	return tests
}

func (m OpenAPI) Specs(in model.OrderedMap[model.SpanQuery, model.NamedAssertions]) openapi.TestSpecs {

	specs := make([]openapi.TestSpecsSpecs, in.Len())

	i := 0
	in.Map(func(spanQuery model.SpanQuery, namedAssertions model.NamedAssertions) {
		assertions := make([]openapi.Assertion, len(namedAssertions.Assertions))
		for j, a := range namedAssertions.Assertions {
			assertions[j] = m.Assertion(a)
		}

		specs[i] = openapi.TestSpecsSpecs{
			Name:       &namedAssertions.Name,
			Selector:   m.Selector(spanQuery),
			Assertions: assertions,
		}
		i++
	})

	return openapi.TestSpecs{
		Specs: specs,
	}
}

func (m OpenAPI) Selector(in model.SpanQuery) openapi.Selector {
	structuredSelector := selectors.FromSpanQuery(in)

	spanSelectors := make([]openapi.SpanSelector, 0)
	for _, spanSelector := range structuredSelector.SpanSelectors {
		spanSelectors = append(spanSelectors, m.SpanSelector(spanSelector))
	}

	return openapi.Selector{
		Query:     string(in),
		Structure: spanSelectors,
	}
}

func (m OpenAPI) SpanSelector(in selectors.SpanSelector) openapi.SpanSelector {
	filters := make([]openapi.SelectorFilter, 0)
	for _, filter := range in.Filters {
		filters = append(filters, openapi.SelectorFilter{
			Property: filter.Property,
			Operator: filter.Operation.Name,
			Value:    filter.Value.AsString(),
		})
	}

	var pseudoClass *openapi.SelectorPseudoClass
	if in.PseudoClass != nil {
		pseudoClass = &openapi.SelectorPseudoClass{
			Name: in.PseudoClass.Name(),
		}

		if nthChildPseudoClass, ok := in.PseudoClass.(*selectors.NthChildPseudoClass); ok {
			pseudoClass.Argument = int32(nthChildPseudoClass.N)
		}
	}

	var child *openapi.SpanSelector
	if in.ChildSelector != nil {
		childSelector := m.SpanSelector(*in.ChildSelector)
		child = &childSelector
	}

	return openapi.SpanSelector{
		Filters:       filters,
		PseudoClass:   pseudoClass,
		ChildSelector: child,
	}
}

func (m OpenAPI) Result(in *model.RunResults) openapi.AssertionResults {
	if in == nil {
		return openapi.AssertionResults{}
	}

	results := make([]openapi.AssertionResultsResults, in.Results.Len())

	i := 0
	in.Results.Map(func(query model.SpanQuery, inRes []model.AssertionResult) {
		res := make([]openapi.AssertionResult, len(inRes))
		for j, r := range inRes {
			sres := make([]openapi.AssertionSpanResult, len(r.Results))
			for k, asr := range r.Results {
				sid := ""
				if asr.SpanID != nil {
					sid = asr.SpanID.String()
				}
				sres[k] = openapi.AssertionSpanResult{
					SpanId:        sid,
					ObservedValue: asr.ObservedValue,
					Passed:        asr.CompareErr == nil,
					Error:         errToString(asr.CompareErr),
				}
			}

			if m.traceConversionConfig.IsTimeField(r.Assertion.Attribute.String()) {
				for i, result := range sres {
					intValue, _ := strconv.Atoi(result.ObservedValue)
					result.ObservedValue = traces.ConvertNanoSecondsIntoProperTimeUnit(intValue)
					sres[i] = result
				}
			}

			res[j] = openapi.AssertionResult{
				AllPassed:   r.AllPassed,
				Assertion:   m.Assertion(r.Assertion),
				SpanResults: sres,
			}
		}
		results[i] = openapi.AssertionResultsResults{
			Selector: m.Selector(query),
			Results:  res,
		}
		i++
	})

	return openapi.AssertionResults{
		AllPassed: in.AllPassed,
		Results:   results,
	}
}

func (m OpenAPI) Assertion(in model.Assertion) openapi.Assertion {
	return openapi.Assertion{
		Attribute:  in.Attribute.String(),
		Comparator: in.Comparator.String(),
		Expected:   in.Value.String(),
	}
}

func (m OpenAPI) AssertionExpression(in *parser.Expression) *model.AssertionExpression {
	if in == nil {
		return nil
	}

	return &model.AssertionExpression{
		LiteralValue: model.LiteralValue{
			Value: in.LiteralValue.String(false),
			Type:  in.LiteralValue.Type(),
		},
		Operation:  in.Operation,
		Expression: m.AssertionExpression(in.Expression),
	}
}

func (m OpenAPI) Run(in *model.Run) openapi.TestRun {
	if in == nil {
		return openapi.TestRun{}
	}

	return openapi.TestRun{
		Id:                        strconv.Itoa(in.ID),
		TraceId:                   in.TraceID.String(),
		SpanId:                    in.SpanID.String(),
		State:                     string(in.State),
		LastErrorState:            errToString(in.LastError),
		ExecutionTime:             int32(in.ExecutionTime()),
		CreatedAt:                 in.CreatedAt,
		ServiceTriggeredAt:        in.ServiceTriggeredAt,
		ServiceTriggerCompletedAt: in.ServiceTriggerCompletedAt,
		ObtainedTraceAt:           in.ObtainedTraceAt,
		CompletedAt:               in.CompletedAt,
		TriggerResult:             m.TriggerResult(in.TriggerResult),
		TestVersion:               int32(in.TestVersion),
		Trace:                     m.Trace(in.Trace),
		Result:                    m.Result(in.Results),
		Metadata:                  in.Metadata,
	}
}

func (m OpenAPI) Runs(in []model.Run) []openapi.TestRun {
	runs := make([]openapi.TestRun, len(in))
	for i, t := range in {
		runs[i] = m.Run(&t)
	}

	return runs
}

// in
type Model struct {
	comparators           comparator.Registry
	traceConversionConfig traces.ConversionConfig
}

func (m Model) Test(in openapi.Test) (model.Test, error) {
	definition, err := m.Definition(in.Specs)
	if err != nil {
		return model.Test{}, fmt.Errorf("could not convert definition: %w", err)
	}
	return model.Test{
		ID:               id.ID(in.Id),
		Name:             in.Name,
		Description:      in.Description,
		ServiceUnderTest: m.Trigger(in.ServiceUnderTest),
		Specs:            definition,
		Version:          int(in.Version),
	}, nil
}

func (m Model) Tests(in []openapi.Test) ([]model.Test, error) {
	tests := make([]model.Test, len(in))
	for i, t := range in {
		test, err := m.Test(t)
		if err != nil {
			return []model.Test{}, fmt.Errorf("could not convert test: %w", err)
		}
		tests[i] = test
	}

	return tests, nil
}

func (m Model) ValidateDefinition(in openapi.TestSpecs) error {
	selectors := map[string]bool{}
	for _, d := range in.Specs {
		if _, exists := selectors[d.Selector.Query]; exists {
			return fmt.Errorf("duplicated selector %s", d.Selector.Query)
		}

		selectors[d.Selector.Query] = true
	}

	return nil
}

func (m Model) Definition(in openapi.TestSpecs) (model.OrderedMap[model.SpanQuery, model.NamedAssertions], error) {
	specs := model.OrderedMap[model.SpanQuery, model.NamedAssertions]{}
	for _, spec := range in.Specs {
		asserts := make([]model.Assertion, len(spec.Assertions))
		for i, a := range spec.Assertions {
			assertion, err := m.Assertion(a)
			if err != nil {
				return model.OrderedMap[model.SpanQuery, model.NamedAssertions]{}, fmt.Errorf("could not convert assertion: %w", err)
			}
			asserts[i] = assertion
		}
		name := ""
		if spec.Name != nil {
			name = *spec.Name
		}

		namedAssertions := model.NamedAssertions{
			Name:       name,
			Assertions: asserts,
		}
		specs, _ = specs.Add(model.SpanQuery(spec.Selector.Query), namedAssertions)
	}

	return specs, nil
}

func (m Model) Run(in openapi.TestRun) (*model.Run, error) {
	tid, _ := trace.TraceIDFromHex(in.TraceId)
	sid, _ := trace.SpanIDFromHex(in.SpanId)
	id, _ := strconv.Atoi(in.Id)
	result, err := m.Result(in.Result)

	if err != nil {
		return &model.Run{}, fmt.Errorf("could not convert result: %w", err)
	}

	return &model.Run{
		ID:                        id,
		TraceID:                   tid,
		SpanID:                    sid,
		State:                     model.RunState(in.State),
		LastError:                 stringToErr(in.LastErrorState),
		CreatedAt:                 in.CreatedAt,
		ServiceTriggeredAt:        in.ServiceTriggeredAt,
		ServiceTriggerCompletedAt: in.ServiceTriggerCompletedAt,
		ObtainedTraceAt:           in.ObtainedTraceAt,
		CompletedAt:               in.CompletedAt,
		TestVersion:               int(in.TestVersion),
		TriggerResult:             m.TriggerResult(in.TriggerResult),
		Trace:                     m.Trace(in.Trace),
		Results:                   result,
		Metadata:                  in.Metadata,
	}, nil
}

func (m Model) Trigger(in openapi.Trigger) model.Trigger {
	return model.Trigger{
		Type: model.TriggerType(in.TriggerType),
		HTTP: m.HTTPRequest(in.TriggerSettings.Http),
		GRPC: m.GRPCRequest(in.TriggerSettings.Grpc),
	}
}

func (m Model) TriggerResult(in openapi.TriggerResult) model.TriggerResult {

	return model.TriggerResult{
		Type: model.TriggerType(in.TriggerType),
		HTTP: m.HTTPResponse(in.TriggerResult.Http),
		GRPC: m.GRPCResponse(in.TriggerResult.Grpc),
	}
}

func (m Model) Result(in openapi.AssertionResults) (*model.RunResults, error) {
	results := model.OrderedMap[model.SpanQuery, []model.AssertionResult]{}

	for _, res := range in.Results {
		ars := make([]model.AssertionResult, len(res.Results))
		for i, r := range res.Results {
			sars := make([]model.SpanAssertionResult, len(r.SpanResults))
			for j, sar := range r.SpanResults {
				var sid *trace.SpanID
				if sar.SpanId != "" {
					s, _ := trace.SpanIDFromHex(sar.SpanId)
					sid = &s
				}
				sars[j] = model.SpanAssertionResult{
					SpanID:        sid,
					ObservedValue: sar.ObservedValue,
					CompareErr:    fmt.Errorf(sar.Error),
				}
			}

			assertion, err := m.Assertion(r.Assertion)
			if err != nil {
				return &model.RunResults{}, fmt.Errorf("could not convert assertion: %w", err)
			}

			ars[i] = model.AssertionResult{
				Assertion: assertion,
				AllPassed: r.AllPassed,
				Results:   sars,
			}
		}
		results, _ = results.Add(model.SpanQuery(res.Selector.Query), ars)
	}

	return &model.RunResults{
		AllPassed: in.AllPassed,
		Results:   results,
	}, nil
}

func (m Model) Assertion(in openapi.Assertion) (model.Assertion, error) {
	expression, err := parser.ParseAssertionExpression(in.Expected)
	if err != nil {
		return model.Assertion{}, err
	}

	comp, _ := m.comparators.Get(in.Comparator)
	return model.Assertion{
		Attribute:  model.Attribute(in.Attribute),
		Comparator: comp,
		Value:      m.AssertionExpression(expression),
	}, nil
}

func (m Model) AssertionExpression(in *parser.Expression) *model.AssertionExpression {
	if in == nil {
		return nil
	}

	return &model.AssertionExpression{
		LiteralValue: model.LiteralValue{
			Value: in.LiteralValue.String(false),
			Type:  in.LiteralValue.Type(),
		},
		Operation:  in.Operation,
		Expression: m.AssertionExpression(in.Expression),
	}
}

func (m Model) Runs(in []openapi.TestRun) ([]model.Run, error) {
	runs := make([]model.Run, len(in))
	for i, r := range in {
		run, err := m.Run(r)
		if err != nil {
			return []model.Run{}, fmt.Errorf("could not convert run: %w", err)
		}
		runs[i] = *run
	}

	return runs, nil
}
