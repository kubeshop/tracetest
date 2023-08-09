package mappings

import (
	"fmt"
	"strconv"
	"time"

	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/assertions/selectors"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/test/trigger"
	"github.com/kubeshop/tracetest/server/testsuite"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/kubeshop/tracetest/server/variableset"
	"go.opentelemetry.io/otel/trace"
)

// out

type OpenAPI struct {
	traceConversionConfig traces.ConversionConfig
}

func optionalTime(in time.Time) *time.Time {
	if in.IsZero() {
		return nil
	}

	return &in
}

func (m OpenAPI) TestSuiteRun(in testsuite.TestSuiteRun) openapi.TestSuiteRun {
	steps := make([]openapi.TestRun, 0, len(in.Steps))

	for _, step := range in.Steps {
		steps = append(steps, m.Run(&step))
	}

	return openapi.TestSuiteRun{
		Id:                          strconv.Itoa(in.ID),
		Version:                     int32(in.TestSuiteVersion),
		CreatedAt:                   in.CreatedAt,
		CompletedAt:                 in.CompletedAt,
		State:                       string(in.State),
		Steps:                       steps,
		Metadata:                    in.Metadata,
		VariableSet:                 m.VariableSet(in.VariableSet),
		Pass:                        int32(in.Pass),
		Fail:                        int32(in.Fail),
		AllStepsRequiredGatesPassed: in.AllStepsRequiredGatesPassed,
	}
}

func (m OpenAPI) Test(in test.Test) openapi.Test {
	return openapi.Test{
		Id:          string(in.ID),
		Name:        in.Name,
		Description: in.Description,
		Trigger:     m.Trigger(in.Trigger),
		Specs:       m.Specs(in.Specs),
		Version:     int32(*in.Version),
		CreatedAt:   *in.CreatedAt,
		Outputs:     m.Outputs(in.Outputs),
		Summary: openapi.TestSummary{
			Runs: int32(in.Summary.Runs),
			LastRun: openapi.TestSummaryLastRun{
				Time:          optionalTime(in.Summary.LastRun.Time),
				Passes:        int32(in.Summary.LastRun.Passes),
				Fails:         int32(in.Summary.LastRun.Fails),
				AnalyzerScore: int32(in.Summary.LastRun.AnalyzerScore),
			},
		},
	}
}

// TODO: after migrating tests and transactions, we can remove this
func (m OpenAPI) VariableSet(in variableset.VariableSet) openapi.VariableSet {
	return openapi.VariableSet{
		Id:          in.ID.String(),
		Name:        in.Name,
		Description: in.Description,
		Values:      m.VariableSetValues(in.Values),
	}
}

func (m OpenAPI) VariableSetValues(in []variableset.VariableSetValue) []openapi.VariableSetValue {
	values := make([]openapi.VariableSetValue, len(in))
	for i, v := range in {
		values[i] = openapi.VariableSetValue{Key: v.Key, Value: v.Value}
	}

	return values
}

func (m OpenAPI) Outputs(in []test.Output) []openapi.TestOutput {
	res := make([]openapi.TestOutput, 0, len(in))
	for _, output := range in {
		res = append(res, openapi.TestOutput{
			Name:           output.Name,
			Selector:       string(output.Selector),
			SelectorParsed: m.Selector(output.Selector),
			Value:          output.Value,
		})
	}

	return res
}

func (m OpenAPI) Trigger(in trigger.Trigger) openapi.Trigger {
	return openapi.Trigger{
		Type:        string(in.Type),
		HttpRequest: m.HTTPRequest(in.HTTP),
		Grpc:        m.GRPCRequest(in.GRPC),
		Traceid:     m.TraceIDRequest(in.TraceID),
	}
}

func (m OpenAPI) TriggerResult(in trigger.TriggerResult) openapi.TriggerResult {

	return openapi.TriggerResult{
		Type: string(in.Type),
		TriggerResult: openapi.TriggerResultTriggerResult{
			Http:    m.HTTPResponse(in.HTTP),
			Grpc:    m.GRPCResponse(in.GRPC),
			Traceid: m.TraceIDResponse(in.TraceID),
		},
	}
}

func (m OpenAPI) Tests(in []test.Test) []openapi.Test {
	tests := make([]openapi.Test, len(in))
	for i, t := range in {
		tests[i] = m.Test(t)
	}

	return tests
}

func (m OpenAPI) VariableSets(in []variableset.VariableSet) []openapi.VariableSet {
	environments := make([]openapi.VariableSet, len(in))
	for i, t := range in {
		environments[i] = m.VariableSet(t)
	}

	return environments
}

func (m OpenAPI) Specs(in test.Specs) []openapi.TestSpec {
	specs := make([]openapi.TestSpec, len(in))

	for i, spec := range in {
		assertions := make([]string, len(spec.Assertions))
		for j, a := range spec.Assertions {
			assertions[j] = string(a)
		}

		specs[i] = openapi.TestSpec{
			Name:           spec.Name,
			Selector:       string(spec.Selector),
			SelectorParsed: m.Selector(test.SpanQuery(spec.Selector)),
			Assertions:     assertions,
		}
	}

	return specs
}

func (m OpenAPI) Selector(in test.SpanQuery) openapi.Selector {
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

func (m OpenAPI) Result(in *test.RunResults) openapi.AssertionResults {
	if in == nil {
		return openapi.AssertionResults{}
	}

	results := make([]openapi.AssertionResultsResultsInner, in.Results.Len())

	i := 0
	in.Results.ForEach(func(query test.SpanQuery, inRes []test.AssertionResult) error {
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

			res[j] = openapi.AssertionResult{
				AllPassed:   r.AllPassed,
				Assertion:   string(r.Assertion),
				SpanResults: sres,
			}
		}
		results[i] = openapi.AssertionResultsResultsInner{
			Selector: m.Selector(query),
			Results:  res,
		}
		i++
		return nil
	})

	return openapi.AssertionResults{
		AllPassed: in.AllPassed,
		Results:   results,
	}
}

func (m OpenAPI) Run(in *test.Run) openapi.TestRun {
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
		TriggerTime:               int32(in.TriggerTime()),
		CreatedAt:                 in.CreatedAt,
		ServiceTriggeredAt:        in.ServiceTriggeredAt,
		ServiceTriggerCompletedAt: in.ServiceTriggerCompletedAt,
		ObtainedTraceAt:           in.ObtainedTraceAt,
		CompletedAt:               in.CompletedAt,
		TriggerResult:             m.TriggerResult(in.TriggerResult),
		TestVersion:               int32(in.TestVersion),
		Trace:                     m.Trace(in.Trace),
		Result:                    m.Result(in.Results),
		Outputs:                   m.RunOutputs(in.Outputs),
		Metadata:                  in.Metadata,
		VariableSet:               m.VariableSet(in.VariableSet),
		TestSuiteId:               in.TestSuiteID,
		TestSuiteRunId:            in.TestSuiteRunID,
		Linter:                    m.LinterResult(in.Linter),
		RequiredGatesResult:       m.RequiredGatesResult(in.RequiredGatesResult),
	}
}

func (m OpenAPI) RunOutputs(in maps.Ordered[string, test.RunOutput]) []openapi.TestRunOutputsInner {
	res := make([]openapi.TestRunOutputsInner, 0, in.Len())

	in.ForEach(func(key string, val test.RunOutput) error {
		res = append(res, openapi.TestRunOutputsInner{
			Name:   key,
			Value:  val.Value,
			SpanId: val.SpanID,
			Error:  errToString(val.Error),
		})
		return nil
	})

	return res
}

func (m OpenAPI) Runs(in []test.Run) []openapi.TestRun {
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

func (m Model) Test(in openapi.Test) (test.Test, error) {
	definition := m.Definition(in.Specs)
	outputs := m.Outputs(in.Outputs)

	version := int(in.Version)
	return test.Test{
		ID:          id.ID(in.Id),
		Name:        in.Name,
		Description: in.Description,
		Trigger:     m.Trigger(in.Trigger),
		Specs:       definition,
		Outputs:     outputs,
		Version:     &version,
	}, nil
}

func (m Model) Outputs(in []openapi.TestOutput) test.Outputs {
	res := make(test.Outputs, 0, len(in))

	for _, output := range in {
		res = append(res, test.Output{
			Name:     output.Name,
			Selector: test.SpanQuery(output.SelectorParsed.Query),
			Value:    output.Value,
		})
	}

	return res
}

func (m Model) Tests(in []openapi.Test) ([]test.Test, error) {
	tests := make([]test.Test, len(in))
	for i, t := range in {
		testObject, err := m.Test(t)
		if err != nil {
			return []test.Test{}, fmt.Errorf("could not convert test: %w", err)
		}
		tests[i] = testObject
	}

	return tests, nil
}

func (m Model) Definition(in []openapi.TestSpec) test.Specs {
	specs := make(test.Specs, 0, len(in))
	for _, spec := range in {
		asserts := make([]test.Assertion, len(spec.Assertions))
		for i, a := range spec.Assertions {
			assertion := test.Assertion(a)
			asserts[i] = assertion
		}
		specs = append(specs, test.TestSpec{
			Selector:   test.SpanQuery(spec.SelectorParsed.Query),
			Name:       spec.Name,
			Assertions: asserts,
		})
	}

	return specs
}

func (m Model) Run(in openapi.TestRun) (*test.Run, error) {
	tid, _ := trace.TraceIDFromHex(in.TraceId)
	sid, _ := trace.SpanIDFromHex(in.SpanId)
	id, _ := strconv.Atoi(in.Id)
	result, err := m.Result(in.Result)

	if err != nil {
		return &test.Run{}, fmt.Errorf("could not convert result: %w", err)
	}

	return &test.Run{
		ID:                        id,
		TraceID:                   tid,
		SpanID:                    sid,
		State:                     test.RunState(in.State),
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
		Outputs:                   m.RunOutputs(in.Outputs),
		Metadata:                  in.Metadata,
		VariableSet:               m.VariableSet(in.VariableSet),
	}, nil
}

func (m Model) RunOutputs(in []openapi.TestRunOutputsInner) maps.Ordered[string, test.RunOutput] {
	res := maps.Ordered[string, test.RunOutput]{}

	for _, output := range in {
		res.Add(output.Name, test.RunOutput{
			Value:  output.Value,
			Name:   output.Name,
			SpanID: output.SpanId,
			Error:  fmt.Errorf(output.Error),
		})
	}

	return res
}

func (m Model) Trigger(in openapi.Trigger) trigger.Trigger {
	return trigger.Trigger{
		Type:    trigger.TriggerType(in.Type),
		HTTP:    m.HTTPRequest(in.HttpRequest),
		GRPC:    m.GRPCRequest(in.Grpc),
		TraceID: m.TraceIDRequest(in.Traceid),
	}
}

func (m Model) TriggerResult(in openapi.TriggerResult) trigger.TriggerResult {

	return trigger.TriggerResult{
		Type:    trigger.TriggerType(in.Type),
		HTTP:    m.HTTPResponse(in.TriggerResult.Http),
		GRPC:    m.GRPCResponse(in.TriggerResult.Grpc),
		TraceID: m.TraceIDResponse(in.TriggerResult.Traceid),
	}
}

func (m Model) Result(in openapi.AssertionResults) (*test.RunResults, error) {
	results := maps.Ordered[test.SpanQuery, []test.AssertionResult]{}

	for _, res := range in.Results {
		ars := make([]test.AssertionResult, len(res.Results))
		for i, r := range res.Results {
			sars := make([]test.SpanAssertionResult, len(r.SpanResults))
			for j, sar := range r.SpanResults {
				var sid *trace.SpanID
				if sar.SpanId != "" {
					s, _ := trace.SpanIDFromHex(sar.SpanId)
					sid = &s
				}
				sars[j] = test.SpanAssertionResult{
					SpanID:        sid,
					ObservedValue: sar.ObservedValue,
					CompareErr:    fmt.Errorf(sar.Error),
				}
			}

			assertion := test.Assertion(r.Assertion)

			ars[i] = test.AssertionResult{
				Assertion: assertion,
				AllPassed: r.AllPassed,
				Results:   sars,
			}
		}
		results, _ = results.Add(test.SpanQuery(res.Selector.Query), ars)
	}

	return &test.RunResults{
		AllPassed: in.AllPassed,
		Results:   results,
	}, nil
}

func (m Model) Runs(in []openapi.TestRun) ([]test.Run, error) {
	runs := make([]test.Run, len(in))
	for i, r := range in {
		run, err := m.Run(r)
		if err != nil {
			return []test.Run{}, fmt.Errorf("could not convert run: %w", err)
		}
		runs[i] = *run
	}

	return runs, nil
}

func (m Model) VariableSet(in openapi.VariableSet) variableset.VariableSet {
	return variableset.VariableSet{
		ID:          id.ID(in.Id),
		Name:        in.Name,
		Description: in.Description,
		Values:      m.VariableSetValue(in.Values),
	}
}

func (m Model) VariableSetValue(in []openapi.VariableSetValue) []variableset.VariableSetValue {
	values := make([]variableset.VariableSetValue, len(in))
	for i, h := range in {
		values[i] = variableset.VariableSetValue{Key: h.Key, Value: h.Value}
	}

	return values
}
