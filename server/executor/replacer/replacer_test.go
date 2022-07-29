package replacer_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/executor/replacer"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPReplacer(t *testing.T) {
	test := model.Test{
		Name: "A test",
		ServiceUnderTest: model.Trigger{
			Type: model.TriggerTypeHTTP,
			HTTP: &model.HTTPRequest{
				Method: model.HTTPMethodPOST,
				URL:    "http://my-api.com/api/users",
				Headers: []model.HTTPHeader{
					{Key: "X-api-key", Value: "{{ uuid() }}"},
					{Key: "my-key", Value: "my-value"},
				},
				Body: `{ "id": "{{ uuid() }}", "name": "{{ fullName() }}", "age": {{ randomInt(18, 99) }} }`,
			},
		},
	}

	newTest, err := replacer.ReplaceTestPlaceholders(test)
	require.NoError(t, err)
	for _, header := range newTest.ServiceUnderTest.HTTP.Headers {
		assert.NotContains(t, header.Value, "{{")
	}

	assert.NotContains(t, newTest.ServiceUnderTest.HTTP.Body, "{{")
}
