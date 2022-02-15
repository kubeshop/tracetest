package executor_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	openapi "github.com/GIT_USER_ID/GIT_REPO_ID/go"
	"github.com/GIT_USER_ID/GIT_REPO_ID/go/executor"
	"github.com/stretchr/testify/assert"
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

	resp, err := ex.Execute(test)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
