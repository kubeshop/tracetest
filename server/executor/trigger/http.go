package trigger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

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
		req.Header.Set(h.Key, h.Value)
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
		if isJson(b) {
			body = string(b)
		} else {
			body = fmt.Sprintf(`{"content": "%s"}`, string(b))
		}
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

func isJson(content []byte) bool {
	var ignored map[string]interface{}
	err := json.Unmarshal(content, &ignored)

	return err == nil
}
