package workers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/client/mocks"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/workers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPollerWorker(t *testing.T) {
	ctx := context.Background()
	controlPlane := mocks.NewGrpcServer()

	client, err := client.Connect(ctx, controlPlane.Addr())
	require.NoError(t, err)

	pollerWorker := workers.NewPollerWorker(client)

	client.OnPollingRequest(func(ctx context.Context, pr *proto.PollingRequest) error {
		return pollerWorker.Poll(ctx, pr)
	})

	err = client.Start(ctx)
	require.NoError(t, err)

	tempoAPI := createTempoFakeApi()

	pollingRequest := proto.PollingRequest{
		TestID:  "test",
		RunID:   1,
		TraceID: "42a2c381da1a5b3a32bc4988bf2431b0",
		Datastore: &proto.DataStore{
			Type: "tempo",
			Tempo: &proto.TempoConfig{
				Type: "http",
				Http: &proto.HttpClientSettings{
					Url: tempoAPI.URL,
				},
			},
		},
	}

	controlPlane.SendPollingRequest(&pollingRequest)

	time.Sleep(1 * time.Second)

	// expect traces to be sent to endpoint
	pollingResponse := controlPlane.GetLastPollingResponse()
	require.NotNil(t, pollingResponse, "agent did not send polling response back to server")

	assert.Len(t, pollingResponse.Spans, 2)
	assert.Equal(t, "", pollingResponse.Spans[0].ParentId)
	assert.Equal(t, pollingResponse.Spans[0].Id, pollingResponse.Spans[1].ParentId)
}

func createTempoFakeApi() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
			"batches": [{
				"scopeSpans": [
					{
						"spans": [
							{
								"spanId": "42a2c381da1a5b3a32bc4988bf2431b0",
								"parentSpanId": "",
								"name": "root",
								"kind": "internal",
								"startTimeUnixNano": "0",
								"endTimeUnixNano": "100",
								"attributes": [],
								"events": [],
								"status": {"code": "ok"}
							},
							{
								"spanId": "99a2c381da1a5b3a32bc4988bf2431c3",
								"parentSpanId": "42a2c381da1a5b3a32bc4988bf2431b0",
								"name": "span 2",
								"kind": "internal",
								"startTimeUnixNano": "0",
								"endTimeUnixNano": "100",
								"attributes": [],
								"events": [],
								"status": {"code": "ok"}
							}
						]
					}
				]
			}]
		}`))
		w.WriteHeader(http.StatusOK)
	}))
}
