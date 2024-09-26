package workers_test

import (
	"context"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/client/mocks"
	"github.com/kubeshop/tracetest/agent/collector"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/workers"
	"github.com/kubeshop/tracetest/agent/workers/poller"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

func TestTracesWorkerFailure(t *testing.T) {
	ctx := ContextWithTracingEnabled()
	controlPlane := mocks.NewGrpcServer()

	client, err := client.Connect(ctx, controlPlane.Addr(), client.WithInsecure())
	require.NoError(t, err)

	cache := collector.NewTraceCache()
	inMemoryDs := poller.NewInMemoryDatastore(cache)

	worker := workers.NewTracesWorker(client, workers.WithTracesInMemoryDatastore(inMemoryDs))

	client.OnTraceModeRequest(func(ctx context.Context, pr *proto.TraceModeRequest) error {
		return worker.Handler(ctx, pr)
	})

	err = client.Start(ctx)
	require.NoError(t, err)

	request := &proto.TraceModeRequest{
		RequestID: "test",
	}

	controlPlane.SendTraceModeRequest(ctx, request)
	time.Sleep(1 * time.Second)

	resp := controlPlane.GetLastTraceModeResponse()
	require.NotNil(t, resp, "agent did not send graphql introspection response back to server")

	assert.Equal(t, resp.Data.Error.Message, "invalid request")
}

func TestTracesWorkerListTraces(t *testing.T) {
	ctx := ContextWithTracingEnabled()
	controlPlane := mocks.NewGrpcServer()

	client, err := client.Connect(ctx, controlPlane.Addr(), client.WithInsecure())
	require.NoError(t, err)

	cache := collector.NewTraceCache()
	inMemoryDs := poller.NewInMemoryDatastore(cache)

	cache.Append("my-test-1", []*v1.Span{})
	cache.Append("my-test-2", []*v1.Span{})
	cache.Append("my-test-3", []*v1.Span{})
	cache.Append("my-test-4", []*v1.Span{})
	cache.Append("my-test-5", []*v1.Span{})

	worker := workers.NewTracesWorker(client, workers.WithTracesInMemoryDatastore(inMemoryDs))

	client.OnTraceModeRequest(func(ctx context.Context, pr *proto.TraceModeRequest) error {
		return worker.Handler(ctx, pr)
	})

	err = client.Start(ctx)
	require.NoError(t, err)

	request := &proto.TraceModeRequest{
		RequestID: "test",
		ListTracesRequest: &proto.ListTracesRequest{
			Take: 10,
			Skip: 0,
			Datastore: &proto.DataStore{
				Type: "otlp",
			},
		},
	}

	controlPlane.SendTraceModeRequest(ctx, request)
	time.Sleep(1 * time.Second)

	resp := controlPlane.GetLastTraceModeResponse()
	require.NotNil(t, resp, "agent did not send graphql introspection response back to server")

	assert.Equal(t, len(resp.Data.ListTracesResponse.Traces), 5)
}

func TestTracesWorkerGetTrace(t *testing.T) {
	ctx := ContextWithTracingEnabled()
	controlPlane := mocks.NewGrpcServer()

	client, err := client.Connect(ctx, controlPlane.Addr(), client.WithInsecure())
	require.NoError(t, err)

	cache := collector.NewTraceCache()
	inMemoryDs := poller.NewInMemoryDatastore(cache)

	traceID := id.NewRandGenerator().TraceID().String()

	cache.Append(traceID, []*v1.Span{{
		SpanId:  []byte(id.NewRandGenerator().SpanID().String()),
		TraceId: []byte(traceID),
		Name:    "test-span",
	}})

	worker := workers.NewTracesWorker(client, workers.WithTracesInMemoryDatastore(inMemoryDs))

	client.OnTraceModeRequest(func(ctx context.Context, pr *proto.TraceModeRequest) error {
		return worker.Handler(ctx, pr)
	})

	err = client.Start(ctx)
	require.NoError(t, err)

	request := &proto.TraceModeRequest{
		RequestID: "test",
		GetTraceRequest: &proto.GetTraceRequest{
			TraceID: traceID,
			Datastore: &proto.DataStore{
				Type: "otlp",
			},
		},
	}

	controlPlane.SendTraceModeRequest(ctx, request)
	time.Sleep(1 * time.Second)

	resp := controlPlane.GetLastTraceModeResponse()
	require.NotNil(t, resp, "agent did not send graphql introspection response back to server")

	assert.Nil(t, resp.Data.Error)
	assert.Equal(t, len(resp.Data.GetTraceResponse.Spans), 1)
	assert.Equal(t, resp.Data.GetTraceResponse.Spans[0].Name, "test-span")
}

func TestTracesWorkerNotSupportedDatastore(t *testing.T) {
	ctx := ContextWithTracingEnabled()
	controlPlane := mocks.NewGrpcServer()

	client, err := client.Connect(ctx, controlPlane.Addr(), client.WithInsecure())
	require.NoError(t, err)

	cache := collector.NewTraceCache()
	inMemoryDs := poller.NewInMemoryDatastore(cache)

	traceID := id.NewRandGenerator().TraceID().String()

	cache.Append(traceID, []*v1.Span{{
		SpanId:  []byte(id.NewRandGenerator().SpanID().String()),
		TraceId: []byte(traceID),
		Name:    "test-span",
	}})

	worker := workers.NewTracesWorker(client, workers.WithTracesInMemoryDatastore(inMemoryDs))

	client.OnTraceModeRequest(func(ctx context.Context, pr *proto.TraceModeRequest) error {
		return worker.Handler(ctx, pr)
	})

	err = client.Start(ctx)
	require.NoError(t, err)

	request := &proto.TraceModeRequest{
		RequestID: "test",
		GetTraceRequest: &proto.GetTraceRequest{
			TraceID: traceID,
			Datastore: &proto.DataStore{
				Type: "jaeger",
				Jaeger: &proto.JaegerConfig{
					Grpc: &proto.GrpcClientSettings{
						Endpoint: "localhost:16685",
					},
				},
			},
		},
	}

	controlPlane.SendTraceModeRequest(ctx, request)
	time.Sleep(1 * time.Second)

	resp := controlPlane.GetLastTraceModeResponse()
	require.NotNil(t, resp, "agent did not send graphql introspection response back to server")

	assert.NotNil(t, resp.Data.Error)
	assert.Equal(t, resp.Data.Error.Message, workers.ErrNotSupportedDataStore.Error())
}
