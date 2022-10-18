package definition_test

import (
	"os"
	"testing"

	"github.com/kubeshop/tracetest/server/encoding/yaml/definition"
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
		expected definition.File
		fn       func(t *testing.T, expected, actual definition.File)
	}{
		{
			name: "TestHTTPTrigger",
			yaml: readFile("./testdata/test_http.yaml"),
			expected: definition.File{
				Type: definition.FileTypeTest,
				Spec: definition.Test{
					ID:   "ZsoMdf44R",
					Name: "Get example",
					Trigger: definition.TestTrigger{
						Type: "http",
						HTTPRequest: definition.HTTPRequest{
							URL:    "http://test.com/list",
							Method: "GET",
						},
					},
					Specs: []definition.TestSpec{
						{
							Selector: `span[name = "Tracetest trigger"]`,
							Assertions: []string{
								"tracetest.response.status = 200",
							},
						},
					},
				},
			},
			fn: func(t *testing.T, expected, actual definition.File) {
				test, err := actual.Test()
				require.NoError(t, err)
				assert.Equal(t, expected.Spec.(definition.Test), test)
			},
		},

		{
			name: "Transaction",
			yaml: readFile("./testdata/transaction.yaml"),
			expected: definition.File{
				Type: definition.FileTypeTransaction,
				Spec: definition.Transaction{
					ID:   "ZsoMdf44R",
					Name: "Get example",
					Steps: []string{
						"step 1",
						"step 2",
					},
				},
			},
			fn: func(t *testing.T, expected, actual definition.File) {
				test, err := actual.Transaction()
				require.NoError(t, err)
				assert.Equal(t, expected.Spec.(definition.Transaction), test)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cl := c
			t.Parallel()

			actual, err := definition.Decode(cl.yaml)
			require.NoError(t, err)

			assert.Equal(t, cl.expected.Type, actual.Type)
			assert.Equal(t, cl.expected.Spec, actual.Spec)

			cl.fn(t, cl.expected, actual)
		})
	}
}
