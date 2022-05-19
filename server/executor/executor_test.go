package executor_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kubeshop/tracetest/executor"
	"github.com/kubeshop/tracetest/id"
	"github.com/kubeshop/tracetest/model"
	"github.com/stretchr/testify/assert"
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

	test := model.Test{
		Name: "test",
		ServiceUnderTest: model.ServiceUnderTest{
			Request: model.HTTPRequest{
				URL:    server.URL,
				Method: model.HTTPMethodGET,
				Headers: []model.HTTPHeader{
					{Key: "Key1", Value: "Value1"},
				},
				Body: "body",
			},
		},
	}

	resp, err := ex.Execute(test, id.NewRandGenerator().TraceID(), id.NewRandGenerator().SpanID())
	assert.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
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

	test := model.Test{
		Name: "test",
		ServiceUnderTest: model.ServiceUnderTest{
			Request: model.HTTPRequest{
				URL:    server.URL,
				Method: model.HTTPMethodPOST,
				Headers: []model.HTTPHeader{
					{Key: "Key1", Value: "Value1"},
				},
				Body: "body",
			},
		},
	}

	resp, err := ex.Execute(test, id.NewRandGenerator().TraceID(), id.NewRandGenerator().SpanID())
	assert.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
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

	test := model.Test{
		Name: "test",
		ServiceUnderTest: model.ServiceUnderTest{
			Request: model.HTTPRequest{
				URL:    server.URL,
				Method: model.HTTPMethodPOST,
				Headers: []model.HTTPHeader{
					{Key: "Key1", Value: "Value1"},
				},
				Auth: &model.HTTPAuthenticator{
					Type: "apiKey",
					Props: map[string]string{
						"key":   "key",
						"value": "value",
						"in":    string(model.APIKeyPositionHeader),
					},
				},
				Body: "body",
			},
		},
	}

	resp, err := ex.Execute(test, id.NewRandGenerator().TraceID(), id.NewRandGenerator().SpanID())
	assert.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
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

	test := model.Test{
		Name: "test",
		ServiceUnderTest: model.ServiceUnderTest{
			Request: model.HTTPRequest{
				URL:    server.URL,
				Method: model.HTTPMethodPOST,
				Headers: []model.HTTPHeader{
					{Key: "Key1", Value: "Value1"},
				},
				Auth: &model.HTTPAuthenticator{
					Type: "basic",
					Props: map[string]string{
						"username": "username",
						"password": "password",
					},
				},
				Body: "body",
			},
		},
	}

	resp, err := ex.Execute(test, id.NewRandGenerator().TraceID(), id.NewRandGenerator().SpanID())
	assert.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
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

	test := model.Test{
		Name: "test",
		ServiceUnderTest: model.ServiceUnderTest{
			Request: model.HTTPRequest{
				URL:    server.URL,
				Method: model.HTTPMethodPOST,
				Headers: []model.HTTPHeader{
					{Key: "Key1", Value: "Value1"},
				},
				Auth: &model.HTTPAuthenticator{
					Type: "bearer",
					Props: map[string]string{
						"token": "token",
					},
				},
				Body: "body",
			},
		},
	}

	resp, err := ex.Execute(test, id.NewRandGenerator().TraceID(), id.NewRandGenerator().SpanID())
	assert.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "OK", resp.Body)
}
