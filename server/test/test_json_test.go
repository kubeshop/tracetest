package test_test

import (
	"encoding/json"
	"testing"

	"github.com/kubeshop/tracetest/server/pkg/maps"
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

	assert.Equal(t, test.SpanQuery("span[tracetest.span.type=\"general\" name=\"Tracetest trigger\"]"), testObject.Specs[0].Selector)
	assert.Equal(t, "my check", testObject.Specs[0].Name)
	assert.Len(t, testObject.Specs[0].Assertions, 1)
	assert.Equal(t, test.Assertion("attr:name = \"Tracetest trigger\""), testObject.Specs[0].Assertions[0])

	assert.Equal(t, test.SpanQuery("span[name=\"GET /api/tests\"]"), testObject.Specs[1].Selector)
	assert.Equal(t, "validate status", testObject.Specs[1].Name)
	assert.Len(t, testObject.Specs[1].Assertions, 1)
	assert.Equal(t, test.Assertion("attr:http.status = 200"), testObject.Specs[1].Assertions[0])
}

func TestV1WithEmptySelector(t *testing.T) {
	specsJSONWithEmptySelector := `[
		{
			"Key": "",
			"Value": {
				"Name": "DURATION_CHECK",
				"Assertions": ["attr:tracetest.span.duration < 2s"]
			}
		},
		{
			"Key": "span[tracetest.span.type=\"database\"]",
			"Value": {
				"Name": "All Database Spans: Processing time is less than 100ms",
				"Assertions": ["attr:tracetest.span.duration < 100ms"]
			}
		}
	]`
	testObject := test.Test{}
	err := json.Unmarshal([]byte(specsJSONWithEmptySelector), &testObject.Specs)

	require.NoError(t, err)
	assert.Equal(t, test.SpanQuery(""), testObject.Specs[0].Selector)
	assert.Equal(t, "DURATION_CHECK", testObject.Specs[0].Name)
	assert.Len(t, testObject.Specs[0].Assertions, 1)
	assert.Equal(t, test.Assertion("attr:tracetest.span.duration < 2s"), testObject.Specs[0].Assertions[0])

	assert.Equal(t, test.SpanQuery("span[tracetest.span.type=\"database\"]"), testObject.Specs[1].Selector)
	assert.Equal(t, "All Database Spans: Processing time is less than 100ms", testObject.Specs[1].Name)
	assert.Len(t, testObject.Specs[1].Assertions, 1)
	assert.Equal(t, test.Assertion("attr:tracetest.span.duration < 100ms"), testObject.Specs[1].Assertions[0])
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

	assert.Equal(t, test.SpanQuery("span[tracetest.span.type=\"general\" name=\"Tracetest trigger\"]"), testObject.Specs[0].Selector)
	assert.Equal(t, "my check", testObject.Specs[0].Name)
	assert.Len(t, testObject.Specs[0].Assertions, 1)
	assert.Equal(t, test.Assertion("attr:name = \"Tracetest trigger\""), testObject.Specs[0].Assertions[0])

	assert.Equal(t, test.SpanQuery("span[name=\"GET /api/tests\"]"), testObject.Specs[1].Selector)
	assert.Equal(t, "validate status", testObject.Specs[1].Name)
	assert.Len(t, testObject.Specs[1].Assertions, 1)
	assert.Equal(t, test.Assertion("attr:http.status = 200"), testObject.Specs[1].Assertions[0])
}

func TestOutputsV1(t *testing.T) {
	v1Format := maps.Ordered[string, test.Output]{}
	v1Format = v1Format.
		MustAdd("USER_ID", test.Output{
			Selector: test.SpanQuery(`span[name = "user creation"]`),
			Value:    `attr:user_id`,
		}).
		MustAdd("USER_NAME", test.Output{
			Selector: test.SpanQuery(`span[name = "user creation"]`),
			Value:    `attr:user_name`,
		})

	v1Json, err := json.Marshal(v1Format)
	require.NoError(t, err)

	testObject := test.Test{}
	err = json.Unmarshal([]byte(v1Json), &testObject.Outputs)

	require.NoError(t, err)
	require.Len(t, testObject.Outputs, 2)

	assert.Equal(t, "USER_ID", testObject.Outputs[0].Name)
	assert.Equal(t, test.SpanQuery(`span[name = "user creation"]`), testObject.Outputs[0].Selector)
	assert.Equal(t, `attr:user_id`, testObject.Outputs[0].Value)

	assert.Equal(t, "USER_NAME", testObject.Outputs[1].Name)
	assert.Equal(t, test.SpanQuery(`span[name = "user creation"]`), testObject.Outputs[1].Selector)
	assert.Equal(t, `attr:user_name`, testObject.Outputs[1].Value)
}

func TestOutputsV2(t *testing.T) {
	v2Format := make([]test.Output, 0)
	v2Format = append(v2Format, test.Output{
		Name:     "USER_ID",
		Selector: test.SpanQuery(`span[name = "user creation"]`),
		Value:    `attr:user_id`,
	})
	v2Format = append(v2Format, test.Output{
		Name:     "USER_NAME",
		Selector: test.SpanQuery(`span[name = "user creation"]`),
		Value:    `attr:user_name`,
	})

	v2Json, err := json.Marshal(v2Format)
	require.NoError(t, err)

	testObject := test.Test{}
	err = json.Unmarshal([]byte(v2Json), &testObject.Outputs)

	require.NoError(t, err)
	require.Len(t, testObject.Outputs, 2)
	assert.Equal(t, "USER_ID", testObject.Outputs[0].Name)
	assert.Equal(t, test.SpanQuery(`span[name = "user creation"]`), testObject.Outputs[0].Selector)
	assert.Equal(t, `attr:user_id`, testObject.Outputs[0].Value)

	assert.Equal(t, "USER_NAME", testObject.Outputs[1].Name)
	assert.Equal(t, test.SpanQuery(`span[name = "user creation"]`), testObject.Outputs[1].Selector)
	assert.Equal(t, `attr:user_name`, testObject.Outputs[1].Value)
}
