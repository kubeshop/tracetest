package config_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/containerd/containerd/log"
	"github.com/kubeshop/tracetest/server/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type deprecatedOptionTest struct {
	name                   string
	fileContent            string
	expectedDeprecatedKeys []string
}

func TestDeprecatedOptions(t *testing.T) {
	testCases := []deprecatedOptionTest{
		{
			name:                   "should_not_detect_any_deprecated_config",
			expectedDeprecatedKeys: []string{},
			fileContent: `
				postgres:
				  host: localhost
				  port: 5432
			`,
		},
		{
			name:                   "should_detect_any_deprecated_config",
			expectedDeprecatedKeys: []string{"postgresConnString"},
			fileContent: `
				postgresConnString: "this is deprecated"
				postgres:
				  host: localhost
				  port: 5432
			`,
		},
	}

	newFunction(t, testCases)
}

func newFunction(t *testing.T, testCases []deprecatedOptionTest) {
	for _, testCase := range testCases {
		testCase.fileContent = strings.ReplaceAll(testCase.fileContent, "\t", "")

		fileName := "tracetest.yaml"
		err := ioutil.WriteFile(fileName, []byte(testCase.fileContent), 0644)
		require.NoError(t, err)
		defer os.Remove(fileName)

		r, w, err := os.Pipe()
		require.NoError(t, err)

		old := log.L.Logger.Out
		log.L.Logger.Out = w
		defer func() {
			log.L.Logger.Out = old
		}()

		outC := make(chan string)

		go func() {
			var buf bytes.Buffer
			io.Copy(&buf, r)
			outC <- buf.String()
		}()

		config.New(nil)

		w.Close()
		output := <-outC

		warnings := make([]string, 0)
		for _, line := range strings.Split(output, "\n") {
			if line == "" {
				continue
			}

			warnings = append(warnings, line)
		}

		assert.Len(t, warnings, len(testCase.expectedDeprecatedKeys))
		for _, message := range testCase.expectedDeprecatedKeys {
			assert.Contains(t, output, message)
		}
	}
}
