package trigger

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/kubeshop/tracetest/server/expression"
	"github.com/kubeshop/tracetest/server/model"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func HTTP() Triggerer {
	return &httpTriggerer{}
}

type httpTriggerer struct{}

func newSpanContext(ctx context.Context) trace.SpanContext {
	spanCtx := trace.SpanContextFromContext(ctx)
	var (
		tid trace.TraceID
		sid trace.SpanID
	)
	if spanCtx.IsValid() {
		tid = spanCtx.TraceID()
		sid = spanCtx.SpanID()
	}

	tracestate, _ := trace.ParseTraceState("tracetest=true")
	var tf trace.TraceFlags
	return trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: tf.WithSampled(true),
		TraceState: tracestate,
		Remote:     true,
	})
}

func (te *httpTriggerer) Trigger(ctx context.Context, test model.Test, opts *TriggerOptions) (Response, error) {
	response := Response{
		Result: model.TriggerResult{
			Type: te.Type(),
		},
	}

	trigger := test.ServiceUnderTest
	if trigger.Type != model.TriggerTypeHTTP {
		return response, fmt.Errorf(`trigger type "%s" not supported by HTTP triggerer`, trigger.Type)
	}

	client := http.DefaultClient

	ctx = trace.ContextWithSpanContext(ctx, newSpanContext(ctx))

	var req *http.Request
	tReq := trigger.HTTP
	var body io.Reader
	if tReq.Body != "" {
		body = bytes.NewBufferString(tReq.Body)
	}
	req, err := http.NewRequest(strings.ToUpper(string(tReq.Method)), tReq.URL, body)
	if err != nil {
		return response, err
	}
	for _, h := range tReq.Headers {
		req.Header.Add(h.Key, h.Value)
	}

	tReq.Authenticate(req)
	propagators().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return response, err
	}

	mapped := mapResp(resp)
	response.Result.HTTP = &mapped
	response.SpanAttributes = map[string]string{
		"tracetest.run.trigger.http.response_code": strconv.Itoa(resp.StatusCode),
	}

	return response, nil
}

func (t *httpTriggerer) Type() model.TriggerType {
	return model.TriggerTypeHTTP
}

func (t *httpTriggerer) Resolve(ctx context.Context, test model.Test, opts *TriggerOptions) (model.Test, error) {
	http := test.ServiceUnderTest.HTTP

	// add quotes before resolving the statements
	// require users to add explicit interpolation to fields

	if http == nil {
		return test, fmt.Errorf("no settings provided for HTTP triggerer")
	}

	url, err := opts.Executor.ResolveStatement(fmt.Sprintf("\"%s\"", http.URL))

	log.Println("resolved url", url, http.URL)

	if err != nil {
		return test, err
	}

	http.URL = url

	headers := []model.HTTPHeader{}
	for _, h := range http.Headers {
		h.Key, err = opts.Executor.ResolveStatement(fmt.Sprintf("\"%s\"", h.Key))
		log.Println("resolved header key", h.Key)
		if err != nil {
			return test, err
		}

		h.Value, err = opts.Executor.ResolveStatement(fmt.Sprintf("\"%s\"", h.Value))
		log.Println("resolved header value", h.Value)
		if err != nil {
			return test, err
		}

		headers = append(headers, h)
	}
	http.Headers = headers

	http.Body, err = opts.Executor.ResolveStatement(fmt.Sprintf("'%s'", http.Body))
	log.Println("resolved body", http.Body)
	if err != nil {
		return test, err
	}

	http.Auth, err = resolveAuth(http.Auth, opts.Executor)
	if err != nil {
		return test, err
	}

	test.ServiceUnderTest.HTTP = http

	log.Println("resolved HTTP triggerer", http)

	return test, nil
}

func resolveAuth(auth *model.HTTPAuthenticator, executor expression.Executor) (*model.HTTPAuthenticator, error) {
	if auth == nil {
		return nil, nil
	}

	for k, v := range auth.Props {
		resolved, err := executor.ResolveStatement(fmt.Sprintf("\"%s\"", v))
		if err != nil {
			return auth, err
		}

		auth.Props[k] = resolved
	}

	return auth, nil
}

func mapResp(resp *http.Response) model.HTTPResponse {
	var mappedHeaders []model.HTTPHeader
	for key, headers := range resp.Header {
		for _, val := range headers {
			val := model.HTTPHeader{
				Key:   key,
				Value: val,
			}
			mappedHeaders = append(mappedHeaders, val)
		}
	}
	var body string
	if b, err := io.ReadAll(resp.Body); err == nil {
		body = string(b)
	} else {
		fmt.Println(err)
	}

	return model.HTTPResponse{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Headers:    mappedHeaders,
		Body:       body,
	}
}
