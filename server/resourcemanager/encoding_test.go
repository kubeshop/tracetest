package resourcemanager_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kubeshop/tracetest/server/resourcemanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncoder(t *testing.T) {
	buildRequest := func(t *testing.T, body string) *http.Request {
		req, err := http.NewRequest("GET", "/", strings.NewReader(body))
		require.NoError(t, err)
		return req
	}

	testEncoder := func(
		t *testing.T,
		r *http.Request,
		inputContentType string, expectedRequestDecoded any,
		outputContentType string, sampleResponseDecoded any, expectedResponseEncoded string,
	) {
		// new encoder
		enc := resourcemanager.EncoderFromRequest(r)

		// process request body
		actualRequestDecoded := map[string]string{}
		err := enc.DecodeRequestBody(&actualRequestDecoded)
		require.NoError(t, err)

		// process response
		rec := httptest.NewRecorder()
		err = enc.WriteEncodedResponse(rec, sampleResponseDecoded)
		require.NoError(t, err)
		resp := rec.Result()
		actualResponseDecoded, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, inputContentType, enc.RequestContentType())
		assert.Equal(t, outputContentType, enc.ResponseContentType())
		assert.Equal(t, outputContentType, resp.Header.Get("Content-Type"))

		assert.Equal(t, expectedRequestDecoded, actualRequestDecoded)
		assert.Equal(t, expectedResponseEncoded, string(actualResponseDecoded))
	}

	t.Run("Default", func(t *testing.T) {
		sampleRequestEncoded := `{"name": "example"}`
		expectedRequestDecoded := map[string]string{"name": "example"}

		sampleResponseDecoded := map[string]string{"greeting": "hi example"}
		expectedResponseEncoded := `{"greeting":"hi example"}`

		// example request/response
		req := buildRequest(t, sampleRequestEncoded)

		testEncoder(
			t,
			req,

			"application/json",
			expectedRequestDecoded,

			"application/json",
			sampleResponseDecoded,
			expectedResponseEncoded,
		)
	})

	t.Run("JSON", func(t *testing.T) {
		sampleRequestEncoded := `{"name": "example"}`
		expectedRequestDecoded := map[string]string{"name": "example"}

		sampleResponseDecoded := map[string]string{"greeting": "hi example"}
		expectedResponseEncoded := `{"greeting":"hi example"}`

		// example request/response
		req := buildRequest(t, sampleRequestEncoded)
		req.Header.Set("Content-Type", "application/json")

		testEncoder(
			t,
			req,

			"application/json",
			expectedRequestDecoded,

			"application/json",
			sampleResponseDecoded,
			expectedResponseEncoded,
		)
	})

	t.Run("Yaml", func(t *testing.T) {
		sampleRequestEncoded := `
name: example
`
		expectedRequestDecoded := map[string]string{"name": "example"}

		sampleResponseDecoded := map[string]string{"greeting": "hi example"}
		expectedResponseEncoded := `greeting: hi example
`

		// example request/response
		req := buildRequest(t, sampleRequestEncoded)
		req.Header.Set("Content-Type", "text/yaml")

		testEncoder(
			t,
			req,
			"text/yaml",
			expectedRequestDecoded,

			"text/yaml",
			sampleResponseDecoded,
			expectedResponseEncoded,
		)
	})

	t.Run("Mixed", func(t *testing.T) {
		sampleRequestEncoded := `{"name": "example"}`
		expectedRequestDecoded := map[string]string{"name": "example"}

		sampleResponseDecoded := map[string]string{"greeting": "hi example"}
		expectedResponseEncoded := `greeting: hi example
`

		// example request/response
		req := buildRequest(t, sampleRequestEncoded)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "text/yaml")

		testEncoder(
			t,
			req,
			"application/json",
			expectedRequestDecoded,

			"text/yaml",
			sampleResponseDecoded,
			expectedResponseEncoded,
		)
	})

}
