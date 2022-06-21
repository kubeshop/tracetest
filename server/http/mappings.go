package http

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/encoding/yaml/conversion"
	"github.com/kubeshop/tracetest/server/encoding/yaml/definition"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/trace"
)

type OpenAPIMapper struct{}

func (m OpenAPIMapper) TestDefinitionFile(in model.Test) definition.Test {
	testDefinition, _ := conversion.ConvertOpenAPITestIntoDefinitionObject(m.Test(in))
	return testDefinition
}

func (m OpenAPIMapper) Test(in model.Test) openapi.Test {
	return openapi.Test{
		Id:          in.ID.String(),
		Name:        in.Name,
		Description: in.Description,
		ServiceUnderTest: openapi.TestServiceUnderTest{
			Request: m.HTTPRequest(in.ServiceUnderTest.Request),
		},
		Definition: m.Definition(in.Definition),
		Version:    int32(in.Version),
	}
}

func (m OpenAPIMapper) HTTPHeaders(in []model.HTTPHeader) []openapi.HttpHeader {
	headers := make([]openapi.HttpHeader, len(in))
	for i, h := range in {
		headers[i] = openapi.HttpHeader{Key: h.Key, Value: h.Value}
	}

	return headers
}

func (m OpenAPIMapper) HTTPRequest(in model.HTTPRequest) openapi.HttpRequest {
	return openapi.HttpRequest{
		Url:     in.URL,
		Method:  string(in.Method),
		Headers: m.HTTPHeaders(in.Headers),
		Body:    in.Body,
		Auth:    m.Auth(in.Auth),
	}
}

func (m OpenAPIMapper) HTTPResponse(in model.HTTPResponse) openapi.HttpResponse {

	return openapi.HttpResponse{
		Status:     in.Status,
		StatusCode: int32(in.StatusCode),
		Headers:    m.HTTPHeaders(in.Headers),
		Body:       in.Body,
	}
}

func (m OpenAPIMapper) Auth(in *model.HTTPAuthenticator) openapi.HttpAuth {
	if in == nil {
		return openapi.HttpAuth{}
	}

	auth := openapi.HttpAuth{
		Type: in.Type,
	}
	switch in.Type {
	case "apiKey":
		auth.ApiKey = openapi.HttpAuthApiKey{
			Key:   in.Props["key"],
			Value: in.Props["value"],
			In:    in.Props["in"],
		}
	case "basic":
		auth.Basic = openapi.HttpAuthBasic{
			Username: in.Props["username"],
			Password: in.Props["password"],
		}
	case "bearer":
		auth.Bearer = openapi.HttpAuthBearer{
			Token: in.Props["bearer"],
		}
	}

	return auth
}

func (m OpenAPIMapper) Tests(in []model.Test) []openapi.Test {
	tests := make([]openapi.Test, len(in))
	for i, t := range in {
		tests[i] = m.Test(t)
	}

	return tests
}

func (m OpenAPIMapper) Definition(in model.OrderedMap[model.SpanQuery, []model.Assertion]) openapi.TestDefinition {

	defs := make([]openapi.TestDefinitionDefinitions, in.Len())

	i := 0
	in.Map(func(spanQuery model.SpanQuery, asserts []model.Assertion) {
		assertions := make([]openapi.Assertion, len(asserts))
		for j, a := range asserts {
			assertions[j] = m.Assertion(a)
		}

		defs[i] = openapi.TestDefinitionDefinitions{
			Selector:   string(spanQuery),
			Assertions: assertions,
		}
		i++
	})

	return openapi.TestDefinition{
		Definitions: defs,
	}
}

func (m OpenAPIMapper) Trace(in *traces.Trace) openapi.Trace {
	if in == nil {
		return openapi.Trace{}
	}

	flat := map[string]openapi.Span{}
	for id, span := range in.Flat {
		flat[id.String()] = m.Span(*span)
	}

	return openapi.Trace{
		TraceId: in.ID.String(),
		Tree:    m.Span(in.RootSpan),
		Flat:    flat,
	}
}

func (m OpenAPIMapper) Span(in traces.Span) openapi.Span {
	parentID := ""
	if in.Parent != nil {
		parentID = in.Parent.ID.String()
	}

	return openapi.Span{
		Id:         in.ID.String(),
		ParentId:   parentID,
		StartTime:  int32(in.StartTime.UnixMilli()),
		EndTime:    int32(in.EndTime.UnixMilli()),
		Attributes: map[string]string(in.Attributes),
		Children:   m.Spans(in.Children),
	}
}

func (m OpenAPIMapper) Spans(in []*traces.Span) []openapi.Span {
	spans := make([]openapi.Span, len(in))
	for i, s := range in {
		spans[i] = m.Span(*s)
	}

	return spans
}

func (m OpenAPIMapper) Result(in *model.RunResults) openapi.AssertionResults {
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
				sres[k] = openapi.AssertionSpanResult{
					SpanId:        asr.SpanID.String(),
					ObservedValue: asr.ObservedValue,
					Passed:        asr.CompareErr == nil,
					Error:         errToString(asr.CompareErr),
				}
			}
			res[j] = openapi.AssertionResult{
				AllPassed:   r.AllPassed,
				Assertion:   m.Assertion(r.Assertion),
				SpanResults: sres,
			}
		}
		results[i] = openapi.AssertionResultsResults{
			Selector: string(query),
			Results:  res,
		}
		i++
	})

	return openapi.AssertionResults{
		AllPassed: in.AllPassed,
		Results:   results,
	}
}
func (m OpenAPIMapper) Assertion(in model.Assertion) openapi.Assertion {
	return openapi.Assertion{
		Attribute:  in.Attribute,
		Comparator: in.Comparator.String(),
		Expected:   in.Value,
	}
}

func (m OpenAPIMapper) Run(in *model.Run) openapi.TestRun {
	if in == nil {
		return openapi.TestRun{}
	}
	return openapi.TestRun{
		Id:                        in.ID.String(),
		TraceId:                   in.TraceID.String(),
		SpanId:                    in.SpanID.String(),
		State:                     string(in.State),
		LastErrorState:            errToString(in.LastError),
		ExectutionTime:            int32(in.ExecutionTime()),
		CreatedAt:                 in.CreatedAt,
		ServiceTriggeredAt:        in.ServiceTriggeredAt,
		ServiceTriggerCompletedAt: in.ServiceTriggerCompletedAt,
		ObtainedTraceAt:           in.ObtainedTraceAt,
		CompletedAt:               in.CompletedAt,
		Request:                   m.HTTPRequest(in.Request),
		Response:                  m.HTTPResponse(in.Response),
		TestVersion:               int32(in.TestVersion),
		Trace:                     m.Trace(in.Trace),
		Result:                    m.Result(in.Results),
	}
}

func (m OpenAPIMapper) Runs(in []model.Run) []openapi.TestRun {
	runs := make([]openapi.TestRun, len(in))
	for i, t := range in {
		runs[i] = m.Run(&t)
	}

	return runs
}

type ModelMapper struct {
	Comparators comparator.Registry
}

func (m ModelMapper) Test(in openapi.Test) model.Test {
	id, _ := uuid.Parse(in.Id)
	return model.Test{
		ID:          id,
		Name:        in.Name,
		Description: in.Description,
		ServiceUnderTest: model.ServiceUnderTest{
			Request: m.HTTPRequest(in.ServiceUnderTest.Request),
		},
		Definition: m.Definition(in.Definition),
		Version:    int(in.Version),
	}
}

func (m ModelMapper) HTTPHeaders(in []openapi.HttpHeader) []model.HTTPHeader {
	headers := make([]model.HTTPHeader, len(in))
	for i, h := range in {
		headers[i] = model.HTTPHeader{Key: h.Key, Value: h.Value}
	}

	return headers
}

func (m ModelMapper) HTTPRequest(in openapi.HttpRequest) model.HTTPRequest {
	return model.HTTPRequest{
		URL:     in.Url,
		Method:  model.HTTPMethod(in.Method),
		Headers: m.HTTPHeaders(in.Headers),
		Body:    in.Body,
		Auth:    m.Auth(in.Auth),
	}
}

func (m ModelMapper) HTTPResponse(in openapi.HttpResponse) model.HTTPResponse {
	return model.HTTPResponse{
		Status:     in.Status,
		StatusCode: int(in.StatusCode),
		Headers:    m.HTTPHeaders(in.Headers),
		Body:       in.Body,
	}
}

func (m ModelMapper) Auth(in openapi.HttpAuth) *model.HTTPAuthenticator {
	var props map[string]string
	switch in.Type {
	case "apiKey":
		props = map[string]string{
			"key":   in.ApiKey.Key,
			"value": in.ApiKey.Value,
			"in":    in.ApiKey.In,
		}
	case "basic":
		props = map[string]string{
			"username": in.Basic.Username,
			"password": in.Basic.Password,
		}
	case "bearer":
		props = map[string]string{
			"token": in.Bearer.Token,
		}
	}

	return &model.HTTPAuthenticator{
		Type:  in.Type,
		Props: props,
	}
}

func (m ModelMapper) Tests(in []openapi.Test) []model.Test {
	tests := make([]model.Test, len(in))
	for i, t := range in {
		tests[i] = m.Test(t)
	}

	return tests
}

func (m ModelMapper) ValidateDefinition(in openapi.TestDefinition) error {
	selectors := map[string]bool{}
	for _, d := range in.Definitions {
		if _, exists := selectors[d.Selector]; exists {
			return fmt.Errorf("duplicated selector %s", d.Selector)
		}

		selectors[d.Selector] = true
	}

	return nil
}

func (m ModelMapper) Definition(in openapi.TestDefinition) model.OrderedMap[model.SpanQuery, []model.Assertion] {
	defs := model.OrderedMap[model.SpanQuery, []model.Assertion]{}
	for _, d := range in.Definitions {
		asserts := make([]model.Assertion, len(d.Assertions))
		for i, a := range d.Assertions {
			asserts[i] = m.Assertion(a)
		}
		defs, _ = defs.Add(model.SpanQuery(d.Selector), asserts)
	}

	return defs
}

func (m ModelMapper) Run(in openapi.TestRun) *model.Run {
	id, _ := uuid.Parse(in.Id)
	tid, _ := trace.TraceIDFromHex(in.TraceId)
	sid, _ := trace.SpanIDFromHex(in.SpanId)
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
		Request:                   m.HTTPRequest(in.Request),
		Response:                  m.HTTPResponse(in.Response),
		Trace:                     m.Trace(in.Trace),
		Results:                   m.Result(in.Result),
	}
}

func (m ModelMapper) Result(in openapi.AssertionResults) *model.RunResults {
	results := model.OrderedMap[model.SpanQuery, []model.AssertionResult]{}

	for _, res := range in.Results {
		ars := make([]model.AssertionResult, len(res.Results))
		for i, r := range res.Results {
			sars := make([]model.SpanAssertionResult, len(r.SpanResults))
			for j, sar := range r.SpanResults {
				sid, _ := trace.SpanIDFromHex(sar.SpanId)
				sars[j] = model.SpanAssertionResult{
					SpanID:        &sid,
					ObservedValue: sar.ObservedValue,
					CompareErr:    fmt.Errorf(sar.Error),
				}
			}

			ars[i] = model.AssertionResult{
				Assertion: m.Assertion(r.Assertion),
				AllPassed: r.AllPassed,
				Results:   sars,
			}
		}
		results, _ = results.Add(model.SpanQuery(res.Selector), ars)
	}

	return &model.RunResults{
		AllPassed: in.AllPassed,
		Results:   results,
	}
}

func (m ModelMapper) Assertion(in openapi.Assertion) model.Assertion {
	comp, _ := m.Comparators.Get(in.Comparator)
	return model.Assertion{
		Attribute:  in.Attribute,
		Comparator: comp,
		Value:      in.Expected,
	}
}

func (m ModelMapper) Trace(in openapi.Trace) *traces.Trace {
	tid, _ := trace.TraceIDFromHex(in.TraceId)
	return &traces.Trace{
		ID:       tid,
		RootSpan: m.Span(in.Tree, nil),
	}
}

func (m ModelMapper) Span(in openapi.Span, parent *traces.Span) traces.Span {
	sid, _ := trace.SpanIDFromHex(in.Id)
	span := traces.Span{
		ID:         sid,
		Attributes: in.Attributes,
		Name:       in.Name,
		StartTime:  time.UnixMilli(int64(in.StartTime)),
		EndTime:    time.UnixMilli(int64(in.EndTime)),
		Parent:     parent,
	}
	span.Children = m.Spans(in.Children, &span)

	return span
}

func (m ModelMapper) Spans(in []openapi.Span, parent *traces.Span) []*traces.Span {
	spans := make([]*traces.Span, len(in))
	for i, s := range in {
		span := m.Span(s, parent)
		spans[i] = &span
	}

	return spans
}

func (m ModelMapper) Runs(in []openapi.TestRun) []model.Run {
	runs := make([]model.Run, len(in))
	for i, r := range in {
		runs[i] = *m.Run(r)
	}

	return runs
}

func errToString(err error) string {
	if err != nil {
		return err.Error()
	}

	return ""
}

func stringToErr(s string) error {
	if s != "" {
		return errors.New(s)
	}

	return nil
}
