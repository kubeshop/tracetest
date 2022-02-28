package executor_test

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	openapi "github.com/kubeshop/tracetest/server/go"
	"github.com/kubeshop/tracetest/server/go/executor"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func TestExecute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		t.Log(req.Header)

		tp, ok := req.Header["Traceparent"]
		if !ok {
			t.Fatalf("missing Traceparent header %#v", req.Header)
		}
		assert.Len(t, tp, 1)
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()

	ex, err := executor.New()
	assert.NoError(t, err)

	test := &openapi.Test{
		Name: "test",
		ServiceUnderTest: openapi.TestServiceUnderTest{
			Url: server.URL,
		},
	}

	rnd := rand.New(rand.NewSource(0))
	tid := trace.TraceID{}
	rnd.Read(tid[:])
	sid := trace.SpanID{}
	rnd.Read(sid[:])

	_, err = ex.Execute(test, tid, sid)
	assert.NoError(t, err)
}
