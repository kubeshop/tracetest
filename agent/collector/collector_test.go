package collector_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/agent/collector"
	"github.com/kubeshop/tracetest/agent/collector/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollector(t *testing.T) {
	targetServer, err := mocks.NewOTLPIngestionServer()
	require.NoError(t, err)

	collector.Start(context.Background(), collector.Config{
		HTTPPort:        4318,
		GRPCPort:        4317,
		BatchTimeout:    2 * time.Second,
		RemoteServerURL: targetServer.Addr(),
	})

	tracer, err := targetServer.CreateTracer("localhost:4317")
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
