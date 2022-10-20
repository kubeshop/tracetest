package formatters_test

import (
	"strings"
	"testing"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/stretchr/testify/assert"
)

func TestListOutput(t *testing.T) {
	cases := []struct {
		name     string
		tests    []openapi.Test
		expected string
	}{
		{
			name:     "NoTests",
			tests:    []openapi.Test{},
			expected: `No tests`,
		},
		{
			name: "HaveTests",
			tests: []openapi.Test{
				{
					Id:   openapi.PtrString("123456"),
					Name: openapi.PtrString("Test One"),
				},
				{
					Id:   openapi.PtrString("456789"),
					Name: openapi.PtrString("Test Two"),
				},
			},
			expected: "" + // vs code trims the last whitespace on save. this awful method avoids that\
				" ID       NAME       URL                                \n" +
				"-------- ---------- ------------------------------------\n" +
				" 123456   Test One   http://localhost:11633/test/123456 \n" +
				" 456789   Test Two   http://localhost:11633/test/456789 \n",
		},
	}

	formatter := formatters.TestsList(config.Config{
		Scheme:   "http",
		Endpoint: "localhost:11633",
	})
	for _, c := range cases {
		output := formatter.Format(c.tests)
		assert.Equal(t, strings.Trim(c.expected, "\n"), output)
	}
}
