package config_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
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

	executeDeprecationTestCase(t, testCases)
}

func executeDeprecationTestCase(t *testing.T, testCases []deprecatedOptionTest) {
	for _, testCase := range testCases {
		testCase.fileContent = strings.ReplaceAll(testCase.fileContent, "\t", "")

		fileName := "tracetest.yaml"
		err := ioutil.WriteFile(fileName, []byte(testCase.fileContent), 0644)
		require.NoError(t, err)
		defer os.Remove(fileName)

		observedZapCore, observedLogs := observer.New(zap.InfoLevel)
		observedLogger := zap.New(observedZapCore)

		config.New(nil, observedLogger)

		warnings := make([]string, 0)
		for _, line := range observedLogs.All() {
			if line.Message == "" {
				continue
			}

			warnings = append(warnings, line.Message)
		}

		assert.Len(t, warnings, len(testCase.expectedDeprecatedKeys))
		for _, message := range testCase.expectedDeprecatedKeys {
			assertSlicesContainsString(t, warnings, message)
		}
	}
}

func assertSlicesContainsString(t *testing.T, slice []string, str string) {
	for _, item := range slice {
		if strings.Contains(item, str) {
			return
		}
	}

	assert.Fail(t, fmt.Sprintf("%s was not found", str))
}
