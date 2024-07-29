package workers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/client/mocks"
	"github.com/kubeshop/tracetest/agent/collector"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/workers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTriggerWorker(t *testing.T) (*mocks.GrpcServerMock, collector.TraceCache) {
	controlPlane := mocks.NewGrpcServer()
	cache := collector.NewTraceCache()

	client, err := client.Connect(context.Background(), controlPlane.Addr(), client.WithInsecure())
	require.NoError(t, err)

	triggerWorker := workers.NewTriggerWorker(
		client,
		workers.WithTraceCache(cache),
		workers.WithTriggerStoppableProcessRunner(workers.NewProcessStopper().RunStoppableProcess),
	)

	client.OnTriggerRequest(func(ctx context.Context, tr *proto.TriggerRequest) error {
		err := triggerWorker.Trigger(ctx, tr)
		assert.NoError(t, err, "trigger failed")
		return err
	})

	client.Start(context.Background())

	return controlPlane, cache
}

func TestTrigger(t *testing.T) {
	controlPlane, cache := setupTriggerWorker(t)

	targetServer := createHelloWorldApi()
	traceID := "42a2c381da1a5b3a32bc4988bf2431b0"

	req := &proto.TriggerRequest{
		Metadata: map[string]string{
			"traceparent": "00-e42e7689e67c64b65ddd3a023a2f8f9d-afad25d95241afce-01",
		},
		TestID:  "my test",
		RunID:   1,
		TraceID: traceID,
		Trigger: &proto.Trigger{
			Type: "http",
			Http: &proto.HttpRequest{
				Method: "GET",
				Url:    targetServer.URL,
				Headers: []*proto.HttpHeader{
					{Key: "Content-Type", Value: "application/json"},
				},
			},
		},
	}

	ctx := ContextWithTracingEnabled()

	// make the control plane send a trigger request to the agent
	controlPlane.SendTriggerRequest(ctx, req)
	time.Sleep(1 * time.Second)

	resp := controlPlane.GetLastTriggerResponse()

	require.NotNil(t, resp)
	assert.Equal(t, "http", resp.Data.TriggerResult.Type)
	assert.Equal(t, int32(http.StatusOK), resp.Data.TriggerResult.Http.StatusCode)
	assert.JSONEq(t, `{"hello": "world"}`, string(resp.Data.TriggerResult.Http.Body))

	_, traceIdIsWatched := cache.Get(traceID)
	assert.True(t, traceIdIsWatched)

	assert.Equal(t, req.Metadata["traceparent"], resp.Data.Metadata["traceparent"])
}

func TestTriggerAgainstGoogle(t *testing.T) {
	ctx := context.Background()
	controlPlane, _ := setupTriggerWorker(t)

	traceID := "42a2c381da1a5b3a32bc4988bf2431b0"

	triggerRequest := &proto.TriggerRequest{
		TestID:  "my test",
		RunID:   1,
		TraceID: traceID,
		Trigger: &proto.Trigger{
			Type: "http",
			Http: &proto.HttpRequest{
				Method: "GET",
				Url:    "https://google.com",
				Headers: []*proto.HttpHeader{
					{Key: "Content-Type", Value: "application/json"},
				},
			},
		},
	}

	// make the control plane send a trigger request to the agent
	controlPlane.SendTriggerRequest(ctx, triggerRequest)
	time.Sleep(1 * time.Second)

	response := controlPlane.GetLastTriggerResponse()

	require.NotNil(t, response)
	assert.Equal(t, "http", response.Data.TriggerResult.Type)
	assert.Equal(t, int32(http.StatusOK), response.Data.TriggerResult.Http.StatusCode)
}

func TestTriggerInexistentAPI(t *testing.T) {
	ctx := context.Background()
	controlPlane, _ := setupTriggerWorker(t)

	traceID := "42a2c381da1a5b3a32bc4988bf2431b0"

	triggerRequest := &proto.TriggerRequest{
		TestID:  "my test",
		RunID:   1,
		TraceID: traceID,
		Trigger: &proto.Trigger{
			Type: "http",
			Http: &proto.HttpRequest{
				Method: "GET",
				Url:    "https://localhost:32148", // hopefully no one uses this port
				Headers: []*proto.HttpHeader{
					{Key: "Content-Type", Value: "application/json"},
				},
			},
		},
	}

	// make the control plane send a trigger request to the agent
	controlPlane.SendTriggerRequest(ctx, triggerRequest)
	time.Sleep(1 * time.Second)

	response := controlPlane.GetLastTriggerResponse()

	require.NotNil(t, response)
	assert.NotNil(t, response.Data.TriggerResult.Error)
	assert.Contains(t, response.Data.TriggerResult.Error.Message, "connection refused")
}

func createHelloWorldApi() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"hello": "world"}`))
		w.WriteHeader(http.StatusOK)
	}))
}
