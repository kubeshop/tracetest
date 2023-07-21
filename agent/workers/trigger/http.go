package trigger

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/goware/urlx"
	"github.com/kubeshop/tracetest/server/test/trigger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func HTTP() Triggerer {
	return &httpTriggerer{}
}

type httpTriggerer struct{}

func httpClient(sslVerification bool) http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: time.Second,
		}).DialContext,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !sslVerification},
	}

	return http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

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

func (te *httpTriggerer) Trigger(ctx context.Context, triggerConfig trigger.Trigger, opts *Options) (Response, error) {
	response := Response{
		Result: trigger.TriggerResult{
			Type: te.Type(),
		},
	}

	if triggerConfig.Type != trigger.TriggerTypeHTTP {
		return response, fmt.Errorf(`trigger type "%s" not supported by HTTP triggerer`, triggerConfig.Type)
	}

	client := httpClient(triggerConfig.HTTP.SSLVerification)

	ctx = trace.ContextWithSpanContext(ctx, newSpanContext(ctx))
	ctx, cncl := context.WithTimeout(ctx, 30*time.Second)
	defer cncl()

	tReq := triggerConfig.HTTP
	var body io.Reader
	if tReq.Body != "" {
		body = bytes.NewBufferString(tReq.Body)
	}

	parsedUrl, err := urlx.Parse(tReq.URL)
	if err != nil {
		return response, err
	}

	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(string(tReq.Method)), parsedUrl.String(), body)
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

func (t *httpTriggerer) Type() trigger.TriggerType {
	return trigger.TriggerTypeHTTP
}

func mapResp(resp *http.Response) trigger.HTTPResponse {
	var mappedHeaders []trigger.HTTPHeader
	for key, headers := range resp.Header {
		for _, val := range headers {
			val := trigger.HTTPHeader{
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

	return trigger.HTTPResponse{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Headers:    mappedHeaders,
		Body:       body,
	}
}
