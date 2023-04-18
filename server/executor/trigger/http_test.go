package trigger_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kubeshop/tracetest/server/executor/trigger"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
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

	test := model.Test{
		Name: "test",
		ServiceUnderTest: model.Trigger{
			Type: model.TriggerTypeHTTP,
			HTTP: &model.HTTPRequest{
				URL:    server.URL,
				Method: model.HTTPMethodGET,
				Headers: []model.HTTPHeader{
					{Key: "Key1", Value: "Value1"},
				},
				Body: "body",
			},
		},
	}

	ex := trigger.HTTP()

	resp, err := ex.Trigger(createContext(), test, nil)
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

	test := model.Test{
		Name: "test",
		ServiceUnderTest: model.Trigger{
			Type: model.TriggerTypeHTTP,
			HTTP: &model.HTTPRequest{
				URL:    server.URL,
				Method: model.HTTPMethodPOST,
				Headers: []model.HTTPHeader{
					{Key: "Key1", Value: "Value1"},
				},
				Body: "body",
			},
		},
	}

	ex := trigger.HTTP()

	resp, err := ex.Trigger(createContext(), test, nil)
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

	test := model.Test{
		Name: "test",
		ServiceUnderTest: model.Trigger{
			Type: model.TriggerTypeHTTP,
			HTTP: &model.HTTPRequest{
				URL:    server.URL,
				Method: model.HTTPMethodPOST,
				Headers: []model.HTTPHeader{
					{Key: "Key1", Value: "Value1"},
				},
				Auth: &model.HTTPAuthenticator{
					Type: "apiKey",
					APIKey: model.APIKeyAuthenticator{
						Key:   "key",
						Value: "value",
						In:    model.APIKeyPositionHeader,
					},
				},
				Body: "body",
			},
		},
	}

	ex := trigger.HTTP()

	resp, err := ex.Trigger(createContext(), test, nil)
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

	test := model.Test{
		Name: "test",
		ServiceUnderTest: model.Trigger{
			Type: model.TriggerTypeHTTP,
			HTTP: &model.HTTPRequest{
				URL:    server.URL,
				Method: model.HTTPMethodPOST,
				Headers: []model.HTTPHeader{
					{Key: "Key1", Value: "Value1"},
				},
				Auth: &model.HTTPAuthenticator{
					Type: "basic",
					Basic: model.BasicAuthenticator{
						Username: "username",
						Password: "password",
					},
				},
				Body: "body",
			},
		},
	}

	ex := trigger.HTTP()

	resp, err := ex.Trigger(createContext(), test, nil)
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

	test := model.Test{
		Name: "test",
		ServiceUnderTest: model.Trigger{
			Type: model.TriggerTypeHTTP,
			HTTP: &model.HTTPRequest{
				URL:    server.URL,
				Method: model.HTTPMethodPOST,
				Headers: []model.HTTPHeader{
					{Key: "Key1", Value: "Value1"},
				},
				Auth: &model.HTTPAuthenticator{
					Type: "bearer",
					Bearer: model.BearerAuthenticator{
						Bearer: "token",
					},
				},
				Body: "body",
			},
		},
	}

	ex := trigger.HTTP()

	resp, err := ex.Trigger(createContext(), test, nil)
	assert.NoError(t, err)

	assert.Equal(t, 200, resp.Result.HTTP.StatusCode)
	assert.Equal(t, "OK", resp.Result.HTTP.Body)
}
