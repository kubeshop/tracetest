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
)

type deprecatedOptionTest struct {
	name                   string
	fileContent            string
	expectedDeprecatedKeys []string
}

type logger struct {
	messages []string
}

func (l *logger) Println(in ...any) {
	message := fmt.Sprintf("%s", in...)
	l.messages = append(l.messages, message)
}

func (l *logger) GetMessages() []string {
	return l.messages
}

func TestDeprecatedOptions(t *testing.T) {
	testCases := []deprecatedOptionTest{
		{
			name:                   "shouldNotDetectAnyDeprecatedConfig",
			expectedDeprecatedKeys: []string{},
			fileContent: `
				postgres:
				  host: localhost
				  port: 5432
			`,
		},
		{
			name:                   "shouldDetectAnyDeprecatedConfig",
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

		logger := logger{
			messages: []string{},
		}

		config.New(nil, &logger)

		warnings := make([]string, 0)
		for _, line := range logger.GetMessages() {
			if line == "" {
				continue
			}

			warnings = append(warnings, line)
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
