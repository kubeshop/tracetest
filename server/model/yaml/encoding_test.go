package yaml_test

import (
	"os"
	"testing"

	"github.com/kubeshop/tracetest/server/config/configresource"
	"github.com/kubeshop/tracetest/server/model/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func readFile(path string) []byte {
	f, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return f
}

func TestDecode(t *testing.T) {
	cases := []struct {
		name       string
		yaml       []byte
		file       yaml.File
		testSpecFn func(t *testing.T, expected, actual yaml.File)
	}{
		{
			name: "TestHTTPTrigger",
			yaml: readFile("./testdata/test_http.yaml"),
			file: yaml.File{
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
			testSpecFn: func(t *testing.T, expected, actual yaml.File) {
				test, err := actual.Test()
				require.NoError(t, err)
				assert.Equal(t, expected.Spec.(yaml.Test), test)
			},
		},

		{
			name: "Transaction",
			yaml: readFile("./testdata/transaction.yaml"),
			file: yaml.File{
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
			testSpecFn: func(t *testing.T, expected, actual yaml.File) {
				test, err := actual.Transaction()
				require.NoError(t, err)
				assert.Equal(t, expected.Spec.(yaml.Transaction), test)
			},
		},
		{
			name: "FileTypeConfig",
			yaml: readFile("./testdata/config.yaml"),
			file: yaml.File{
				Type: yaml.FileTypeConfig,
				Spec: configresource.Config{
					ID:               "current",
					Name:             "config",
					AnalyticsEnabled: true,
				},
			},
			testSpecFn: func(t *testing.T, expected, actual yaml.File) {
				config, err := actual.Config()
				require.NoError(t, err)
				assert.Equal(t, expected.Spec.(configresource.Config), config)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cl := c
			t.Parallel()

			t.Run("Decode", func(t *testing.T) {
				actual, err := yaml.Decode(cl.yaml)
				require.NoError(t, err)

				assert.Equal(t, cl.file.Type, actual.Type)
				assert.Equal(t, cl.file.Spec, actual.Spec)

				cl.testSpecFn(t, cl.file, actual)
			})

			t.Run("Encode", func(t *testing.T) {
				actual, err := yaml.Encode(cl.file)
				require.NoError(t, err)

				assert.Equal(t, string(cl.yaml), string(actual))
			})
		})
	}
}
