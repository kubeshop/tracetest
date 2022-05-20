package http

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kubeshop/tracetest/assertions/comparator"
	"github.com/kubeshop/tracetest/model"
	"github.com/kubeshop/tracetest/openapi"
	"github.com/kubeshop/tracetest/traces"
	"go.opentelemetry.io/otel/trace"
)

type openapiMapper struct{}

func (m openapiMapper) Test(in model.Test) openapi.Test {
	return openapi.Test{
		Id:          in.ID.String(),
		Name:        in.Name,
		Description: in.Description,
		ServiceUnderTest: openapi.TestServiceUnderTest{
			Request: m.HTTPRequest(in.ServiceUnderTest.Request),
		},
		Definition:       m.Definition(in.Definition),
		ReferenceTestRun: m.Run(in.ReferenceRun),
	}
}

func (m openapiMapper) HTTPHeaders(in []model.HTTPHeader) []openapi.HttpHeader {
	headers := make([]openapi.HttpHeader, len(in))
	for i, h := range in {
		headers[i] = openapi.HttpHeader{Key: h.Key, Value: h.Value}
	}

	return headers
}

func (m openapiMapper) HTTPRequest(in model.HTTPRequest) openapi.HttpRequest {
	return openapi.HttpRequest{
		Url:     in.URL,
		Method:  string(in.Method),
		Headers: m.HTTPHeaders(in.Headers),
		Body:    in.Body,
		Auth:    m.Auth(in.Auth),
	}
}

func (m openapiMapper) HTTPResponse(in model.HTTPResponse) openapi.HttpResponse {

	return openapi.HttpResponse{
		Status:     in.Status,
		StatusCode: int32(in.StatusCode),
		Headers:    m.HTTPHeaders(in.Headers),
		Body:       in.Body,
	}
}

func (m openapiMapper) Auth(in *model.HTTPAuthenticator) openapi.HttpAuth {
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

func (m openapiMapper) Tests(in []model.Test) []openapi.Test {
	tests := make([]openapi.Test, len(in))
	for i, t := range in {
		tests[i] = m.Test(t)
	}

	return tests
}

func (m openapiMapper) Definition(in model.Definition) openapi.TestDefinition {

	defs := make([]openapi.TestDefinitionDefinitions, len(in))

	i := 0
	for sel, def := range in {
		assertions := make([]openapi.Assertion, len(def))
		for j, a := range def {
			assertions[j] = m.Assertion(a)
		}

		defs[i] = openapi.TestDefinitionDefinitions{
			Selector:   string(sel),
			Assertions: assertions,
		}
		i++
	}

	return openapi.TestDefinition{
		Definitions: defs,
	}
}

func (m openapiMapper) Trace(in *traces.Trace) openapi.Trace {
	if in == nil {
		return openapi.Trace{}
	}

	flat := map[string]openapi.Span{}
	for id, span := range in.Flat {
		flat[id.String()] = openapi.Span{
			Id:         span.ID.String(),
			Attributes: map[string]string(span.Attributes),
		}
	}

	return openapi.Trace{
		TraceId: in.ID.String(),
		Tree:    m.Span(in.RootSpan),
		Flat:    flat,
	}
}

func (m openapiMapper) Span(in traces.Span) openapi.Span {
	return openapi.Span{
		Id:         in.ID.String(),
		StartTime:  in.StartTime,
		EndTime:    in.EndTime,
		Attributes: map[string]string(in.Attributes),
		Children:   m.Spans(in.Children),
	}
}

func (m openapiMapper) Spans(in []*traces.Span) []openapi.Span {
	spans := make([]openapi.Span, len(in))
	for i, s := range in {
		spans[i] = m.Span(*s)
	}

	return spans
}

func (m openapiMapper) Result(in *model.RunResults) openapi.AssertionResults {
	if in == nil {
		return openapi.AssertionResults{}
	}

	results := make([]openapi.AssertionResultsResults, len(in.Results))
	i := 0
	for query, inRes := range in.Results {
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
				Assertion:   m.Assertion(r.Assertion),
				SpanResults: sres,
			}
		}
		results[i] = openapi.AssertionResultsResults{
			Selector: string(query),
			Results:  res,
		}
		i++
	}
	return openapi.AssertionResults{
		AllPassed: in.AllPassed,
		Results:   results,
	}
}
func (m openapiMapper) Assertion(in model.Assertion) openapi.Assertion {
	return openapi.Assertion{
		Attribute:  in.Attribute,
		Comparator: in.Comparator.String(),
		Expected:   in.Value,
	}
}

func (m openapiMapper) Run(in *model.Run) openapi.TestRun {
	if in == nil {
		return openapi.TestRun{}
	}
	return openapi.TestRun{
		Id:             in.ID.String(),
		TraceId:        in.TraceID.String(),
		SpanId:         in.SpanID.String(),
		State:          string(in.State),
		LastErrorState: errToString(in.LastError),
		CreatedAt:      in.CreatedAt,
		CompletedAt:    in.CompletedAt,
		Request:        m.HTTPRequest(in.Request),
		Response:       m.HTTPResponse(in.Response),
		Trace:          m.Trace(in.Trace),
		Result:         m.Result(in.Results),
	}
}

func (m openapiMapper) Runs(in []model.Run) []openapi.TestRun {
	runs := make([]openapi.TestRun, len(in))
	for i, t := range in {
		runs[i] = m.Run(&t)
	}

	return runs
}

type modelMapper struct {
	Comparators comparator.Registry
}

func (m modelMapper) Test(in openapi.Test) model.Test {
	id, _ := uuid.Parse(in.Id)
	return model.Test{
		ID:          id,
		Name:        in.Name,
		Description: in.Description,
		ServiceUnderTest: model.ServiceUnderTest{
			Request: m.HTTPRequest(in.ServiceUnderTest.Request),
		},
		ReferenceRun: m.Run(in.ReferenceTestRun),
		Definition:   m.Definition(in.Definition),
	}
}

func (m modelMapper) HTTPHeaders(in []openapi.HttpHeader) []model.HTTPHeader {
	headers := make([]model.HTTPHeader, len(in))
	for i, h := range in {
		headers[i] = model.HTTPHeader{Key: h.Key, Value: h.Value}
	}

	return headers
}

func (m modelMapper) HTTPRequest(in openapi.HttpRequest) model.HTTPRequest {
	return model.HTTPRequest{
		URL:     in.Url,
		Method:  model.HTTPMethod(in.Method),
		Headers: m.HTTPHeaders(in.Headers),
		Body:    in.Body,
		Auth:    m.Auth(in.Auth),
	}
}

func (m modelMapper) HTTPResponse(in openapi.HttpResponse) model.HTTPResponse {
	return model.HTTPResponse{
		Status:     in.Status,
		StatusCode: int(in.StatusCode),
		Headers:    m.HTTPHeaders(in.Headers),
		Body:       in.Body,
	}
}

func (m modelMapper) Auth(in openapi.HttpAuth) *model.HTTPAuthenticator {
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

func (m modelMapper) Tests(in []openapi.Test) []model.Test {
	tests := make([]model.Test, len(in))
	for i, t := range in {
		tests[i] = m.Test(t)
	}

	return tests
}

func (m modelMapper) Definition(in openapi.TestDefinition) model.Definition {
	defs := model.Definition{}
	for _, d := range in.Definitions {
		asserts := make([]model.Assertion, len(d.Assertions))
		for i, a := range d.Assertions {
			asserts[i] = m.Assertion(a)
		}
		defs[model.SpanQuery(d.Selector)] = asserts
	}

	return defs
}

func (m modelMapper) Run(in openapi.TestRun) *model.Run {
	id, _ := uuid.Parse(in.Id)
	tid, _ := trace.TraceIDFromHex(in.TraceId)
	sid, _ := trace.SpanIDFromHex(in.SpanId)
	return &model.Run{
		ID:          id,
		TraceID:     tid,
		SpanID:      sid,
		State:       model.RunState(in.State),
		LastError:   stringToErr(in.LastErrorState),
		CreatedAt:   in.CreatedAt,
		CompletedAt: in.CompletedAt,
		Request:     m.HTTPRequest(in.Request),
		Response:    m.HTTPResponse(in.Response),
		Trace:       m.Trace(in.Trace),
		Results:     m.Result(in.Result),
	}
}

func (m modelMapper) Result(in openapi.AssertionResults) *model.RunResults {
	results := model.Results{}

	for _, res := range in.Results {
		results[model.SpanQuery(res.Selector)] = make([]model.AssertionResult, len(res.Results))
		for i, r := range res.Results {
			sars := make([]model.SpanAssertionResult, len(r.SpanResults))
			for j, sar := range r.SpanResults {
				sid, _ := trace.SpanIDFromHex(sar.SpanId)
				sars[j] = model.SpanAssertionResult{
					SpanID:        sid,
					ObservedValue: sar.ObservedValue,
					CompareErr:    fmt.Errorf(sar.Error),
				}
			}

			results[model.SpanQuery(res.Selector)][i] = model.AssertionResult{
				Assertion: m.Assertion(r.Assertion),
				AllPassed: r.AllPassed,
				Results:   sars,
			}
		}
	}

	return &model.RunResults{
		AllPassed: in.AllPassed,
		Results:   results,
	}
}

func (m modelMapper) Assertion(in openapi.Assertion) model.Assertion {
	comp, _ := m.Comparators.Get(in.Comparator)
	return model.Assertion{
		Attribute:  in.Attribute,
		Comparator: comp,
		Value:      in.Expected,
	}
}

func (m modelMapper) Trace(in openapi.Trace) *traces.Trace {
	tid, _ := trace.TraceIDFromHex(in.TraceId)
	return &traces.Trace{
		ID:       tid,
		RootSpan: m.Span(in.Tree, nil),
	}
}

func (m modelMapper) Span(in openapi.Span, parent *traces.Span) traces.Span {
	sid, _ := trace.SpanIDFromHex(in.Id)
	span := traces.Span{
		ID:         sid,
		Attributes: in.Attributes,
		Name:       in.Name,
		StartTime:  in.StartTime,
		EndTime:    in.EndTime,
		Parent:     parent,
	}
	span.Children = m.Spans(in.Children, &span)

	return span
}

func (m modelMapper) Spans(in []openapi.Span, parent *traces.Span) []*traces.Span {
	spans := make([]*traces.Span, len(in))
	for i, s := range in {
		span := m.Span(s, parent)
		spans[i] = &span
	}

	return spans
}

func (m modelMapper) Runs(in []openapi.TestRun) []model.Run {
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
