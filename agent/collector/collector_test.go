package collector_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/agent/collector"
	"github.com/kubeshop/tracetest/agent/collector/mocks"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

func TestCollector(t *testing.T) {
	targetServer, err := mocks.NewOTLPIngestionServer()
	require.NoError(t, err)

	noopTracer := trace.NewNoopTracerProvider().Tracer("noop_tracer")

	c, err := collector.Start(
		context.Background(),
		collector.Config{
			HTTPPort:        4318,
			GRPCPort:        4317,
			BatchTimeout:    2 * time.Second,
			RemoteServerURL: targetServer.Addr(),
		},
		noopTracer,
		collector.WithStartRemoteServer(true),
	)
	require.NoError(t, err)

	defer c.Stop()

	tracer, err := mocks.NewTracer(context.Background(), "localhost:4317")
	require.NoError(t, err)

	func(ctx context.Context) {
		for i := 0; i < 10; i++ {
			spanCtx, span := tracer.Start(ctx, fmt.Sprintf("span %d", i))
			ctx = spanCtx

			defer span.End()
		}
	}(context.Background())

	time.Sleep(500 * time.Millisecond)
	// Should not have any spans yet, because batch timeout is 2 seconds
	assert.Len(t, targetServer.ReceivedSpans(), 0)

	// Now after waiting the timeout, it should contain all spans
	time.Sleep(4 * time.Second)
	assert.Len(t, targetServer.ReceivedSpans(), 10)
}

func TestCollectorWatchingSpansFromTest(t *testing.T) {
	targetServer, err := mocks.NewOTLPIngestionServer()
	require.NoError(t, err)

	cache := collector.NewTraceCache()
	noopTracer := trace.NewNoopTracerProvider().Tracer("noop_tracer")

	c, err := collector.Start(
		context.Background(),
		collector.Config{
			HTTPPort:        4318,
			GRPCPort:        4317,
			BatchTimeout:    2 * time.Second,
			RemoteServerURL: targetServer.Addr(),
		},
		noopTracer,
		collector.WithTraceCache(cache),
		collector.WithStartRemoteServer(true),
	)
	require.NoError(t, err)

	defer c.Stop()

	tracer, err := mocks.NewTracer(context.Background(), "localhost:4317")
	require.NoError(t, err)

	watchedTraceID := id.NewRandGenerator().TraceID()
	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: watchedTraceID,
	})

	ctx := trace.ContextWithSpanContext(context.Background(), spanContext)

	cache.Append(watchedTraceID.String(), []*v1.Span{})

	// 10 spans will be watched and stored in the cache
	func(ctx context.Context) {
		for i := 0; i < 10; i++ {
			spanCtx, span := tracer.Start(ctx, fmt.Sprintf("watched span %d", i))
			ctx = spanCtx

			defer span.End()
		}
	}(ctx)

	// 10 spans will not be watched neither stored in the cache
	func(ctx context.Context) {
		for i := 0; i < 10; i++ {
			spanCtx, span := tracer.Start(ctx, fmt.Sprintf("span %d", i))
			ctx = spanCtx

			defer span.End()
		}
	}(context.Background())

	time.Sleep(500 * time.Millisecond)
	// Should not have any spans yet, because batch timeout is 2 seconds
	assert.Len(t, targetServer.ReceivedSpans(), 0)

	// Now after waiting the timeout, it should contain all spans
	time.Sleep(4 * time.Second)
	assert.Len(t, targetServer.ReceivedSpans(), 20)

	cachedSpans, ok := cache.Get(watchedTraceID.String())
	assert.True(t, ok)
	assert.Len(t, cachedSpans, 10)
	for _, span := range cachedSpans {
		isWatched := strings.Contains(span.Name, "watched")
		assert.True(t, isWatched)
	}
}
