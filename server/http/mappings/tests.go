package mappings

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/assertions/selectors"
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

func optionalTime(in time.Time) *time.Time {
	if in.IsZero() {
		return nil
	}

	return &in
}

func (m OpenAPI) Transaction(in model.Transaction) openapi.Transaction {
	steps := make([]openapi.Test, len(in.Steps))
	for i, step := range in.Steps {
		steps[i] = m.Test(step)
	}

	return openapi.Transaction{
		Id:          in.ID.String(),
		Name:        in.Name,
		Description: in.Description,
		Version:     int32(in.Version),
		Steps:       steps,
		CreatedAt:   in.CreatedAt,
		Summary: openapi.TestSummary{
			Runs: int32(in.Summary.Runs),
			LastRun: openapi.TestSummaryLastRun{
				Time:   optionalTime(in.Summary.LastRun.Time),
				Passes: int32(in.Summary.LastRun.Fails),
				Fails:  int32(in.Summary.LastRun.Passes),
			},
		},
	}
}

func (m OpenAPI) TransactionRun(in model.TransactionRun) openapi.TransactionRun {
	steps := make([]openapi.TestRun, 0, len(in.Steps))

	for _, step := range in.Steps {
		steps = append(steps, m.Run(&step))
	}

	return openapi.TransactionRun{
		Id:          strconv.Itoa(in.ID),
		Version:     int32(in.TransactionVersion),
		CreatedAt:   in.CreatedAt,
		CompletedAt: in.CompletedAt,
		State:       string(in.State),
		Steps:       steps,
		Metadata:    in.Metadata,
		Environment: m.Environment(in.Environment),
		Pass:        int32(in.Pass),
		Fail:        int32(in.Fail),
	}
}

func (m OpenAPI) Test(in model.Test) openapi.Test {
	return openapi.Test{
		Id:               string(in.ID),
		Name:             in.Name,
		Description:      in.Description,
		ServiceUnderTest: m.Trigger(in.ServiceUnderTest),
		Specs:            m.Specs(in.Specs),
		Version:          int32(in.Version),
		CreatedAt:        in.CreatedAt,
		Outputs:          m.Outputs(in.Outputs),
		Summary: openapi.TestSummary{
			Runs: int32(in.Summary.Runs),
			LastRun: openapi.TestSummaryLastRun{
				Time:   optionalTime(in.Summary.LastRun.Time),
				Passes: int32(in.Summary.LastRun.Passes),
				Fails:  int32(in.Summary.LastRun.Fails),
			},
		},
	}
}

func (m OpenAPI) Environment(in model.Environment) openapi.Environment {
	return openapi.Environment{
		Id:          in.ID,
		Name:        in.Name,
		Description: in.Description,
		Values:      m.EnvironmentValues(in.Values),
	}
}

func (m OpenAPI) EnvironmentValues(in []model.EnvironmentValue) []openapi.EnvironmentValue {
	values := make([]openapi.EnvironmentValue, len(in))
	for i, v := range in {
		values[i] = openapi.EnvironmentValue{Key: v.Key, Value: v.Value}
	}

	return values
}

func (m OpenAPI) Outputs(in model.OrderedMap[string, model.Output]) []openapi.TestOutput {
	res := make([]openapi.TestOutput, 0, in.Len())
	in.ForEach(func(key string, val model.Output) error {
		res = append(res, openapi.TestOutput{
			Name:     key,
			Selector: m.Selector(val.Selector),
			Value:    val.Value,
		})
		return nil
	})

	return res
}

func (m OpenAPI) Trigger(in model.Trigger) openapi.Trigger {
	return openapi.Trigger{
		TriggerType: string(in.Type),
		TriggerSettings: openapi.TriggerTriggerSettings{
			Http:    m.HTTPRequest(in.HTTP),
			Grpc:    m.GRPCRequest(in.GRPC),
			Traceid: m.TRACEIDRequest(in.TRACEID),
		},
	}
}

func (m OpenAPI) TriggerResult(in model.TriggerResult) openapi.TriggerResult {

	return openapi.TriggerResult{
		TriggerType: string(in.Type),
		TriggerResult: openapi.TriggerResultTriggerResult{
			Http:    m.HTTPResponse(in.HTTP),
			Grpc:    m.GRPCResponse(in.GRPC),
			Traceid: m.TRACEIDResponse(in.TRACEID),
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

func (m OpenAPI) Environments(in []model.Environment) []openapi.Environment {
	environments := make([]openapi.Environment, len(in))
	for i, t := range in {
		environments[i] = m.Environment(t)
	}

	return environments
}

func (m OpenAPI) Specs(in model.OrderedMap[model.SpanQuery, model.NamedAssertions]) openapi.TestSpecs {

	specs := make([]openapi.TestSpecsSpecs, in.Len())

	i := 0
	in.ForEach(func(spanQuery model.SpanQuery, namedAssertions model.NamedAssertions) error {
		assertions := make([]string, len(namedAssertions.Assertions))
		for j, a := range namedAssertions.Assertions {
			assertions[j] = string(a)
		}

		specs[i] = openapi.TestSpecsSpecs{
			Name:       &namedAssertions.Name,
			Selector:   m.Selector(spanQuery),
			Assertions: assertions,
		}
		i++
		return nil
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
	in.Results.ForEach(func(query model.SpanQuery, inRes []model.AssertionResult) error {
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
		results[i] = openapi.AssertionResultsResults{
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
		Environment:               m.Environment(in.Environment),
	}
}

func (m OpenAPI) RunOutputs(in model.OrderedMap[string, model.RunOutput]) []openapi.TestRunOutputs {
	res := make([]openapi.TestRunOutputs, 0, in.Len())

	in.ForEach(func(key string, val model.RunOutput) error {
		res = append(res, openapi.TestRunOutputs{
			Name:   key,
			Value:  val.Value,
			SpanId: val.SpanID,
		})
		return nil
	})

	return res
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
	testRepository        model.TestRepository
}

func (m Model) Transaction(ctx context.Context, in openapi.Transaction) (model.Transaction, error) {
	tests := make([]model.Test, len(in.Steps))
	for i, test := range in.Steps {
		test, err := m.testRepository.GetLatestTestVersion(ctx, id.ID(test.Id))
		if err != nil {
			return model.Transaction{}, err
		}

		tests[i] = test
	}

	return model.Transaction{
		ID:          id.ID(in.Id),
		Name:        in.Name,
		Description: in.Description,
		Version:     int(in.Version),
		Steps:       tests,
	}, nil
}

func (m Model) Test(in openapi.Test) (model.Test, error) {
	definition, err := m.Definition(in.Specs)
	if err != nil {
		return model.Test{}, fmt.Errorf("could not convert definition: %w", err)
	}

	outputs, err := m.Outputs(in.Outputs)
	if err != nil {
		return model.Test{}, fmt.Errorf("could not convert outputs: %w", err)
	}

	return model.Test{
		ID:               id.ID(in.Id),
		Name:             in.Name,
		Description:      in.Description,
		ServiceUnderTest: m.Trigger(in.ServiceUnderTest),
		Specs:            definition,
		Outputs:          outputs,
		Version:          int(in.Version),
	}, nil
}

func (m Model) Outputs(in []openapi.TestOutput) (model.OrderedMap[string, model.Output], error) {
	res := model.OrderedMap[string, model.Output]{}

	var err error
	for _, output := range in {
		res, err = res.Add(output.Name, model.Output{
			Selector: model.SpanQuery(output.Selector.Query),
			Value:    output.Value,
		})

		if err != nil {
			return res, fmt.Errorf("cannot parse outputs: %w", err)
		}
	}

	return res, nil
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
			assertion := model.Assertion(a)
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
		Outputs:                   m.RunOutputs(in.Outputs),
		Metadata:                  in.Metadata,
		Environment:               m.Environment(in.Environment),
	}, nil
}

func (m Model) RunOutputs(in []openapi.TestRunOutputs) model.OrderedMap[string, model.RunOutput] {
	res := model.OrderedMap[string, model.RunOutput]{}

	for _, output := range in {
		res.Add(output.Name, model.RunOutput{
			Value:  output.Value,
			Name:   output.Name,
			SpanID: output.SpanId,
		})
	}

	return res
}

func (m Model) Trigger(in openapi.Trigger) model.Trigger {
	return model.Trigger{
		Type:    model.TriggerType(in.TriggerType),
		HTTP:    m.HTTPRequest(in.TriggerSettings.Http),
		GRPC:    m.GRPCRequest(in.TriggerSettings.Grpc),
		TRACEID: m.TRACEIDRequest(in.TriggerSettings.Traceid),
	}
}

func (m Model) TriggerResult(in openapi.TriggerResult) model.TriggerResult {

	return model.TriggerResult{
		Type:    model.TriggerType(in.TriggerType),
		HTTP:    m.HTTPResponse(in.TriggerResult.Http),
		GRPC:    m.GRPCResponse(in.TriggerResult.Grpc),
		TRACEID: m.TRACEIDResponse(in.TriggerResult.Traceid),
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

			assertion := model.Assertion(r.Assertion)

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

func (m Model) EnvironmentValue(in []openapi.EnvironmentValue) []model.EnvironmentValue {
	values := make([]model.EnvironmentValue, len(in))
	for i, h := range in {
		values[i] = model.EnvironmentValue{Key: h.Key, Value: h.Value}
	}

	return values
}

func (m Model) Environment(in openapi.Environment) model.Environment {
	return model.Environment{
		ID:          in.Id,
		Name:        in.Name,
		Description: in.Description,
		Values:      m.EnvironmentValue(in.Values),
	}
}
