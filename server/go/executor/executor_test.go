package executor_test

import (
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	openapi "github.com/kubeshop/tracetest/server/go"
	"github.com/kubeshop/tracetest/server/go/executor"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func TestExecuteGet(t *testing.T) {
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

	ex, err := executor.New()
	assert.NoError(t, err)

	test := &openapi.Test{
		Name: "test",
		ServiceUnderTest: openapi.TestServiceUnderTest{
			Request: openapi.HttpRequest{
				Url:    server.URL,
				Method: "GET",
				Headers: []openapi.HttpResponseHeaders{
					{Key: "Key1", Value: "Value1"},
				},
				Body: "body",
			},
		},
	}

	rnd := rand.New(rand.NewSource(0))
	tid := trace.TraceID{}
	rnd.Read(tid[:])
	sid := trace.SpanID{}
	rnd.Read(sid[:])

	resp, err := ex.Execute(test, tid, sid)
	assert.NoError(t, err)

	assert.Equal(t, int32(200), resp.Response.StatusCode)
	assert.Equal(t, "OK", resp.Response.Body)
}

func TestExecutePost(t *testing.T) {
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

	ex, err := executor.New()
	assert.NoError(t, err)

	test := &openapi.Test{
		Name: "test",
		ServiceUnderTest: openapi.TestServiceUnderTest{
			Request: openapi.HttpRequest{
				Url:    server.URL,
				Method: "POST",
				Headers: []openapi.HttpResponseHeaders{
					{Key: "Key1", Value: "Value1"},
				},
				Body: "body",
			},
		},
	}

	rnd := rand.New(rand.NewSource(0))
	tid := trace.TraceID{}
	rnd.Read(tid[:])
	sid := trace.SpanID{}
	rnd.Read(sid[:])

	resp, err := ex.Execute(test, tid, sid)
	assert.NoError(t, err)

	assert.Equal(t, int32(200), resp.Response.StatusCode)
	assert.Equal(t, "OK", resp.Response.Body)
}
