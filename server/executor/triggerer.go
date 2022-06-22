package executor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kubeshop/tracetest/server/model"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/contrib/propagators/ot"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type Triggerer interface {
	Trigger(context.Context, model.Test, trace.TraceID, trace.SpanID) (model.HTTPResponse, error)
}

func NewTriggerer(tracer trace.Tracer) (Triggerer, error) {
	tp, err := initTracing()
	if err != nil {
		return nil, err
	}
	return &instrumentedTriggerer{
		tracer: tracer,
		triggerer: &httpTriggerer{
			traceProvider: tp,
		},
	}, nil
}

type instrumentedTriggerer struct {
	tracer    trace.Tracer
	triggerer Triggerer
}

func (te *instrumentedTriggerer) Trigger(ctx context.Context, test model.Test, tid trace.TraceID, sid trace.SpanID) (model.HTTPResponse, error) {
	ctx, span := te.tracer.Start(ctx, "Trigger test")
	defer span.End()

	resp, err := te.triggerer.Trigger(ctx, test, tid, sid)

	span.SetAttributes(
		attribute.String("tracetest.run.trigger.test_id", test.ID.String()),
		attribute.String("tracetest.run.trigger.type", "http"),
		attribute.Int("tracetest.run.trigger.http.response_code", resp.StatusCode),
	)

	return resp, err
}

type httpTriggerer struct {
	traceProvider *sdktrace.TracerProvider
}

func (te *httpTriggerer) Trigger(_ context.Context, test model.Test, tid trace.TraceID, sid trace.SpanID) (model.HTTPResponse, error) {
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport,
			otelhttp.WithTracerProvider(te.traceProvider),
			otelhttp.WithPropagators(propagation.NewCompositeTextMapPropagator(propagation.Baggage{},
				b3.New(),
				jaeger.Jaeger{},
				ot.OT{},
				xray.Propagator{},
				propagation.TraceContext{})),
		),
	}

	var tf trace.TraceFlags
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: tf.WithSampled(true),
		TraceState: trace.TraceState{},
		Remote:     true,
	})

	var req *http.Request
	tReq := test.ServiceUnderTest.Request
	var body io.Reader
	if tReq.Body != "" {
		body = bytes.NewBufferString(tReq.Body)
	}
	req, err := http.NewRequest(strings.ToUpper(string(tReq.Method)), tReq.URL, body)
	if err != nil {
		return model.HTTPResponse{}, err
	}
	for _, h := range tReq.Headers {
		req.Header.Set(h.Key, h.Value)
	}

	test.ServiceUnderTest.Request.Authenticate(req)

	resp, err := client.Do(req.WithContext(trace.ContextWithSpanContext(context.Background(), sc)))
	if err != nil {
		return model.HTTPResponse{}, err
	}

	return mapResp(resp), nil
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

func initTracing() (*sdktrace.TracerProvider, error) {
	// Set standard attributes per semantic conventions
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("tracetest"),
	)

	spanExporter, err := stdouttrace.New(stdouttrace.WithWriter(io.Discard))
	if err != nil {
		return nil, err
	}
	// Create and set the TraceProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(spanExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample())),
	)

	return tp, nil
}
