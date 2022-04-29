package executor_test

import (
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kubeshop/tracetest/executor"
	"github.com/kubeshop/tracetest/openapi"
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

	assert.Equal(t, int32(200), resp.StatusCode)
	assert.Equal(t, "OK", resp.Body)
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

	assert.Equal(t, int32(200), resp.StatusCode)
	assert.Equal(t, "OK", resp.Body)
}

func TestExecutePostWithApiKeyAuth(t *testing.T) {
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
				Auth: openapi.HttpAuth{
					Type: "apiKey",
					ApiKey: openapi.HttpAuthApiKey{
						Key:   "key",
						Value: "value",
						In:    "header",
					},
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

	assert.Equal(t, int32(200), resp.StatusCode)
	assert.Equal(t, "OK", resp.Body)
}

func TestExecutePostWithBasicAuth(t *testing.T) {
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
				Auth: openapi.HttpAuth{
					Type: "basic",
					Basic: openapi.HttpAuthBasic{
						Username: "username",
						Password: "password",
					},
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

	assert.Equal(t, int32(200), resp.StatusCode)
	assert.Equal(t, "OK", resp.Body)
}

func TestExecutePostWithBearerAuth(t *testing.T) {
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
				Auth: openapi.HttpAuth{
					Type: "bearer",
					Bearer: openapi.HttpAuthBearer{
						Token: "token",
					},
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

	assert.Equal(t, int32(200), resp.StatusCode)
	assert.Equal(t, "OK", resp.Body)
}
