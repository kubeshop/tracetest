package middleware

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	semconv "go.opentelemetry.io/collector/semconv/v1.5.0"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type httpMetricMiddleware struct {
	next                     http.Handler
	requestDurationHistogram metric.Int64Histogram
}

func (m *httpMetricMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw := NewStatusCodeCapturerWriter(w)
	initialTime := time.Now()
	m.next.ServeHTTP(rw, r)
	duration := time.Since(initialTime)

	route := mux.CurrentRoute(r)
	pathTemplate, _ := route.GetPathTemplate()
	if pathTemplate == "" {
		pathTemplate = "/"
	}

	metricAttributes := []attribute.KeyValue{
		attribute.String(semconv.AttributeHTTPRoute, pathTemplate),
		attribute.String(semconv.AttributeHTTPMethod, r.Method),
		attribute.Int(semconv.AttributeHTTPStatusCode, rw.statusCode),
	}

	if tenantID := TenantIDFromContext(r.Context()); tenantID != "" {
		metricAttributes = append(metricAttributes, attribute.String("tracetest.tenant_id", tenantID))
	}

	m.requestDurationHistogram.Record(r.Context(), duration.Milliseconds(), metric.WithAttributeSet(
		attribute.NewSet(metricAttributes...),
	))
}

var _ http.Handler = &httpMetricMiddleware{}

type responseWriter struct {
	http.ResponseWriter
	responseBody []byte
	statusCode   int
}

func NewStatusCodeCapturerWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		w,
		[]byte{},
		http.StatusOK,
	}
}

func (w *responseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(body []byte) (int, error) {
	w.responseBody = body
	return w.ResponseWriter.Write(body)
}

func (w *responseWriter) StatusCode() int {
	return w.statusCode
}

func (w *responseWriter) Body() []byte {
	return w.responseBody
}

func NewMetricMiddleware(meter metric.Meter) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		durationHistogram, _ := meter.Int64Histogram("http.server.latency", metric.WithUnit("ms"))

		return &httpMetricMiddleware{
			next:                     next,
			requestDurationHistogram: durationHistogram,
		}
	}
}
