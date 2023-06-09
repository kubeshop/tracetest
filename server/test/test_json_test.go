package test_test

import (
	"encoding/json"
	"testing"

	"github.com/kubeshop/tracetest/server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpecV1(t *testing.T) {
	oldSpecFormat := `
	[
		{
			"Key": "span[tracetest.span.type=\"general\" name=\"Tracetest trigger\"]",
			"Value": {
				"Name": "my check",
				"Assertions": [
					"attr:name = \"Tracetest trigger\""
				]
			}
		},
		{
			"Key": "span[name=\"GET /api/tests\"]",
			"Value": {
				"Name": "validate status",
				"Assertions": [
					"attr:http.status = 200"
				]
			}
		}
	]
	`

	testObject := test.Test{}
	err := json.Unmarshal([]byte(oldSpecFormat), &testObject.Specs)

	require.NoError(t, err)
	require.Len(t, testObject.Specs, 2)

	assert.Equal(t, test.SpanQuery("span[tracetest.span.type=\"general\" name=\"Tracetest trigger\"]"), testObject.Specs[0].Selector.Query)
	assert.Equal(t, "my check", testObject.Specs[0].Name)
	assert.Len(t, testObject.Specs[0].Assertions, 1)
	assert.Equal(t, test.Assertion("attr:name = \"Tracetest trigger\""), testObject.Specs[0].Assertions[0])

	assert.Equal(t, test.SpanQuery("span[name=\"GET /api/tests\"]"), testObject.Specs[1].Selector.Query)
	assert.Equal(t, "validate status", testObject.Specs[1].Name)
	assert.Len(t, testObject.Specs[1].Assertions, 1)
	assert.Equal(t, test.Assertion("attr:http.status = 200"), testObject.Specs[1].Assertions[0])
}

func TestSpecV2(t *testing.T) {
	specFormat := `
	[
		{
			"selector": "span[tracetest.span.type=\"general\" name=\"Tracetest trigger\"]",
			"name": "my check",
			"assertions": [
				"attr:name = \"Tracetest trigger\""
			]
		},
		{
			"selector": "span[name=\"GET /api/tests\"]",
			"name": "validate status",
			"assertions": [
				"attr:http.status = 200"
			]
		}
	]
	`

	testObject := test.Test{}
	err := json.Unmarshal([]byte(specFormat), &testObject.Specs)

	require.NoError(t, err)
	require.Len(t, testObject.Specs, 2)

	assert.Equal(t, test.SpanQuery("span[tracetest.span.type=\"general\" name=\"Tracetest trigger\"]"), testObject.Specs[0].Selector.Query)
	assert.Equal(t, "my check", testObject.Specs[0].Name)
	assert.Len(t, testObject.Specs[0].Assertions, 1)
	assert.Equal(t, test.Assertion("attr:name = \"Tracetest trigger\""), testObject.Specs[0].Assertions[0])

	assert.Equal(t, test.SpanQuery("span[name=\"GET /api/tests\"]"), testObject.Specs[1].Selector.Query)
	assert.Equal(t, "validate status", testObject.Specs[1].Name)
	assert.Len(t, testObject.Specs[1].Assertions, 1)
	assert.Equal(t, test.Assertion("attr:http.status = 200"), testObject.Specs[1].Assertions[0])
}
