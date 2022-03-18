package executor

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	openapi "github.com/kubeshop/tracetest/server/go"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type TestExecutor struct {
	traceProvider *sdktrace.TracerProvider
}

func New() (*TestExecutor, error) {
	tp, err := initTracing()
	if err != nil {
		return nil, err
	}
	return &TestExecutor{
		traceProvider: tp,
	}, nil
}

func (te *TestExecutor) Execute(test *openapi.Test, tid trace.TraceID, sid trace.SpanID) (*openapi.TestRunResult, error) {
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport,
			otelhttp.WithTracerProvider(te.traceProvider),
			otelhttp.WithPropagators(propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})),
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

	req, err := http.NewRequest(http.MethodGet, test.ServiceUnderTest.Url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req.WithContext(trace.ContextWithSpanContext(context.Background(), sc)))
	if err != nil {
		return nil, err
	}

	return &openapi.TestRunResult{
		ResultId: uuid.New().String(),
		Response: mapResp(resp),
	}, nil
}

func mapResp(resp *http.Response) openapi.HttpResponse {
	var mappedHeaders []openapi.HttpResponseHeaders
	for key, headers := range resp.Header {
		for _, val := range headers {
			val := openapi.HttpResponseHeaders{
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

	return openapi.HttpResponse{
		Status:     resp.Status,
		StatusCode: int32(resp.StatusCode),
		Headers:    mappedHeaders,
		Body:       body,
	}
}

func initTracing() (*sdktrace.TracerProvider, error) {
	// Set standard attributes per semantic conventions
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("projectx"),
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
