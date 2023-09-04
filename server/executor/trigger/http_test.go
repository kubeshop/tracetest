package trigger_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	triggerer "github.com/kubeshop/tracetest/server/executor/trigger"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/test/trigger"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func createContext() context.Context {
	return trace.ContextWithSpanContext(context.TODO(), trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: id.NewRandGenerator().TraceID(),
		SpanID:  id.NewRandGenerator().SpanID(),
	}))
}

func TestTriggerGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		t.Log(req.Header)

		tp, ok := req.Header["Traceparent"]
		if !ok {
			t.Fatalf("missing Traceparent header %#v", req.Header)
		}
		assert.Len(t, tp, 1)

		assert.Equal(t, "GET", req.Method)
		assert.Equal(t, "Value1", req.Header.Get("Key1"))

		b, err := io.ReadAll(req.Body)
		assert.NoError(t, err)
		assert.Equal(t, "body", string(b))

		rw.WriteHeader(200)
		_, err = rw.Write([]byte(`OK`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	test := test.Test{
		Name: "test",
		Trigger: trigger.Trigger{
			Type: trigger.TriggerTypeHTTP,
			HTTP: &trigger.HTTPRequest{
				URL:    server.URL,
				Method: trigger.HTTPMethodGET,
				Headers: []trigger.HTTPHeader{
					{Key: "Key1", Value: "Value1"},
				},
				Body: "body",
			},
		},
	}

	ex := triggerer.HTTP()

	resp, err := ex.Trigger(createContext(), test)
	assert.NoError(t, err)

	assert.Equal(t, 200, resp.Result.HTTP.StatusCode)
	assert.Equal(t, "OK", resp.Result.HTTP.Body)
}

func TestTriggerPost(t *testing.T) {
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
		assert.Equal(t, "body", string(b))

		rw.WriteHeader(200)
		_, err = rw.Write([]byte(`OK`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	test := test.Test{
		Name: "test",
		Trigger: trigger.Trigger{
			Type: trigger.TriggerTypeHTTP,
			HTTP: &trigger.HTTPRequest{
				URL:    server.URL,
				Method: trigger.HTTPMethodPOST,
				Headers: []trigger.HTTPHeader{
					{Key: "Key1", Value: "Value1"},
				},
				Body: "body",
			},
		},
	}

	ex := triggerer.HTTP()

	resp, err := ex.Trigger(createContext(), test)
	assert.NoError(t, err)

	assert.Equal(t, 200, resp.Result.HTTP.StatusCode)
	assert.Equal(t, "OK", resp.Result.HTTP.Body)
}

func TestTriggerPostWithApiKeyAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		t.Log(req.Header)

		tp, ok := req.Header["Traceparent"]
		if !ok {
			t.Fatalf("missing Traceparent header %#v", req.Header)
		}
		assert.Len(t, tp, 1)

		tp, ok = req.Header["Key"]
		if !ok {
			t.Fatalf("missing key header %#v", req.Header)
		}
		assert.Len(t, tp, 1)
		assert.Equal(t, tp[0], "value")
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, "Value1", req.Header.Get("Key1"))

		b, err := io.ReadAll(req.Body)
		assert.NoError(t, err)
		assert.Equal(t, "body", string(b))

		rw.WriteHeader(200)
		_, err = rw.Write([]byte(`OK`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	test := test.Test{
		Name: "test",
		Trigger: trigger.Trigger{
			Type: trigger.TriggerTypeHTTP,
			HTTP: &trigger.HTTPRequest{
				URL:    server.URL,
				Method: trigger.HTTPMethodPOST,
				Headers: []trigger.HTTPHeader{
					{Key: "Key1", Value: "Value1"},
				},
				Auth: &trigger.HTTPAuthenticator{
					Type: "apiKey",
					APIKey: &trigger.APIKeyAuthenticator{
						Key:   "key",
						Value: "value",
						In:    trigger.APIKeyPositionHeader,
					},
				},
				Body: "body",
			},
		},
	}

	ex := triggerer.HTTP()

	resp, err := ex.Trigger(createContext(), test)
	assert.NoError(t, err)

	assert.Equal(t, 200, resp.Result.HTTP.StatusCode)
	assert.Equal(t, "OK", resp.Result.HTTP.Body)
}

func TestTriggerPostWithBasicAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		t.Log(req.Header)

		tp, ok := req.Header["Traceparent"]
		if !ok {
			t.Fatalf("missing Traceparent header %#v", req.Header)
		}
		assert.Len(t, tp, 1)

		tp, ok = req.Header["Authorization"]
		if !ok {
			t.Fatalf("missing Authorization header %#v", req.Header)
		}
		assert.Len(t, tp, 1)
		assert.Equal(t, tp[0], "Basic dXNlcm5hbWU6cGFzc3dvcmQ=")
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, "Value1", req.Header.Get("Key1"))

		b, err := io.ReadAll(req.Body)
		assert.NoError(t, err)
		assert.Equal(t, "body", string(b))

		rw.WriteHeader(200)
		_, err = rw.Write([]byte(`OK`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	test := test.Test{
		Name: "test",
		Trigger: trigger.Trigger{
			Type: trigger.TriggerTypeHTTP,
			HTTP: &trigger.HTTPRequest{
				URL:    server.URL,
				Method: trigger.HTTPMethodPOST,
				Headers: []trigger.HTTPHeader{
					{Key: "Key1", Value: "Value1"},
				},
				Auth: &trigger.HTTPAuthenticator{
					Type: "basic",
					Basic: &trigger.BasicAuthenticator{
						Username: "username",
						Password: "password",
					},
				},
				Body: "body",
			},
		},
	}

	ex := triggerer.HTTP()

	resp, err := ex.Trigger(createContext(), test)
	assert.NoError(t, err)

	assert.Equal(t, 200, resp.Result.HTTP.StatusCode)
	assert.Equal(t, "OK", resp.Result.HTTP.Body)
}

func TestTriggerPostWithBearerAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		t.Log(req.Header)

		tp, ok := req.Header["Traceparent"]
		if !ok {
			t.Fatalf("missing Traceparent header %#v", req.Header)
		}
		assert.Len(t, tp, 1)

		tp, ok = req.Header["Authorization"]
		if !ok {
			t.Fatalf("missing Authorization header %#v", req.Header)
		}
		assert.Len(t, tp, 1)
		assert.Equal(t, tp[0], "token")
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, "Value1", req.Header.Get("Key1"))

		b, err := io.ReadAll(req.Body)
		assert.NoError(t, err)
		assert.Equal(t, "body", string(b))

		rw.WriteHeader(200)
		_, err = rw.Write([]byte(`OK`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	test := test.Test{
		Name: "test",
		Trigger: trigger.Trigger{
			Type: trigger.TriggerTypeHTTP,
			HTTP: &trigger.HTTPRequest{
				URL:    server.URL,
				Method: trigger.HTTPMethodPOST,
				Headers: []trigger.HTTPHeader{
					{Key: "Key1", Value: "Value1"},
				},
				Auth: &trigger.HTTPAuthenticator{
					Type: "bearer",
					Bearer: &trigger.BearerAuthenticator{
						Bearer: "token",
					},
				},
				Body: "body",
			},
		},
	}

	ex := triggerer.HTTP()

	resp, err := ex.Trigger(createContext(), test)
	assert.NoError(t, err)

	assert.Equal(t, 200, resp.Result.HTTP.StatusCode)
	assert.Equal(t, "OK", resp.Result.HTTP.Body)
}
