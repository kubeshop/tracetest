package trigger_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kubeshop/tracetest/agent/workers/trigger"
	triggerer "github.com/kubeshop/tracetest/agent/workers/trigger"
	"github.com/stretchr/testify/assert"
)

func TestGraphqlTrigger(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		t.Log(req.Header)

		tp, ok := req.Header["Traceparent"]
		if !ok {
			t.Fatalf("missing Traceparent header %#v", req.Header)
		}
		assert.Len(t, tp, 1)

		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, "Value1", req.Header.Get("Key1"))

		b, err := io.ReadAll(req.Body)
		assert.NoError(t, err)
		assert.Equal(t, `query { films { name } }`, string(b))

		rw.Header().Add("Content-Type", "application/json")
		rw.WriteHeader(200)
		_, err = rw.Write([]byte(`{ "data": { "films": [{ "name": "A New Hope" }] } }`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	triggerConfig := trigger.Trigger{
		Type: trigger.TriggerTypeGraphql,
		Graphql: &trigger.GraphqlRequest{
			URL: server.URL,
			Headers: []trigger.HTTPHeader{
				{Key: "Key1", Value: "Value1"},
			},
			Body: `query { films { name } }`,
		},
	}

	httpTriggerer := triggerer.HTTP()

	ex := triggerer.GRAPHQL(httpTriggerer)

	resp, err := ex.Trigger(createContext(), triggerConfig, nil)
	assert.NoError(t, err)

	assert.Equal(t, 200, resp.Result.Graphql.StatusCode)
	assert.Equal(t, `{ "data": { "films": [{ "name": "A New Hope" }] } }`, resp.Result.Graphql.Body)
}
