package yaml_test

import (
	"os"
	"testing"

	"github.com/kubeshop/tracetest/server/model/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func readFile(path string) string {
	f, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(f)
}

func TestDecode(t *testing.T) {
	cases := []struct {
		name     string
		yaml     string
		expected yaml.File
		fn       func(t *testing.T, expected, actual yaml.File)
	}{
		{
			name: "TestHTTPTrigger",
			yaml: readFile("./testdata/test_http.yaml"),
			expected: yaml.File{
				Type: yaml.FileTypeTest,
				Spec: yaml.Test{
					ID:   "ZsoMdf44R",
					Name: "Get example",
					Trigger: yaml.TestTrigger{
						Type: "http",
						HTTPRequest: yaml.HTTPRequest{
							URL:    "http://test.com/list",
							Method: "GET",
						},
					},
					Specs: []yaml.TestSpec{
						{
							Selector: `span[name = "Tracetest trigger"]`,
							Assertions: []string{
								"tracetest.response.status = 200",
							},
						},
					},
				},
			},
			fn: func(t *testing.T, expected, actual yaml.File) {
				test, err := actual.Test()
				require.NoError(t, err)
				assert.Equal(t, expected.Spec.(yaml.Test), test)
			},
		},

		{
			name: "Transaction",
			yaml: readFile("./testdata/transaction.yaml"),
			expected: yaml.File{
				Type: yaml.FileTypeTransaction,
				Spec: yaml.Transaction{
					ID:   "ZsoMdf44R",
					Name: "Get example",
					Steps: []string{
						"step 1",
						"step 2",
					},
				},
			},
			fn: func(t *testing.T, expected, actual yaml.File) {
				test, err := actual.Transaction()
				require.NoError(t, err)
				assert.Equal(t, expected.Spec.(yaml.Transaction), test)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cl := c
			t.Parallel()

			actual, err := yaml.Decode(cl.yaml)
			require.NoError(t, err)

			assert.Equal(t, cl.expected.Type, actual.Type)
			assert.Equal(t, cl.expected.Spec, actual.Spec)

			cl.fn(t, cl.expected, actual)
		})
	}
}
