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
	requestCounter           metric.Int64Counter
}

func (m *httpMetricMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw := NewStatusCodeCapturerWriter(w)
	initialTime := time.Now()
	m.next.ServeHTTP(rw, r)
	duration := time.Since(initialTime)

	route := mux.CurrentRoute(r)
	pathTemplate, _ := route.GetPathTemplate()

	metricAttributes := []attribute.KeyValue{
		attribute.String(semconv.AttributeHTTPRoute, pathTemplate),
		attribute.String(semconv.AttributeHTTPMethod, r.Method),
		attribute.Int(semconv.AttributeHTTPStatusCode, rw.statusCode),
	}

	if tenantID := TenantIDFromContext(r.Context()); tenantID != "" {
		metricAttributes = append(metricAttributes, attribute.String("tracetest.tenant_id", tenantID))
	}

	m.requestDurationHistogram.Record(r.Context(), duration.Milliseconds(), metric.WithAttributes(metricAttributes...))
	m.requestCounter.Add(r.Context(), 1, metric.WithAttributes(metricAttributes...))
}

var _ http.Handler = &httpMetricMiddleware{}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewStatusCodeCapturerWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (lrw *responseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func NewMetricMiddleware(meter metric.Meter) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		durationHistogram, _ := meter.Int64Histogram("http.server.duration", metric.WithUnit("ms"))
		requestCounter, _ := meter.Int64Counter("http.server.request.count")

		return &httpMetricMiddleware{
			next:                     next,
			requestDurationHistogram: durationHistogram,
			requestCounter:           requestCounter,
		}
	}
}
