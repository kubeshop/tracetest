package config_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/stretchr/testify/assert"
)

type logger struct {
	messages []string
}

func newLogger() *logger {
	return &logger{
		messages: []string{},
	}
}

func (l *logger) Println(in ...any) {
	str := fmt.Sprint(in...)

	// ignore config file messages
	if strings.HasPrefix(str, "Config file used:") {
		return
	}
	l.messages = append(l.messages, str)
}

func TestDeprecatedOptions(t *testing.T) {
	testCases := []struct {
		name             string
		flags            []string
		expectedMessages []string
	}{
		{
			name:             "NoDeprecations",
			flags:            []string{"--config", "./testdata/basic.yaml"},
			expectedMessages: []string{},
		},
		{
			name:             "Deprecated/postgresConnString",
			flags:            []string{"--postgresConnString", "some conn string"},
			expectedMessages: []string{`config "postgresConnString" is deprecated. Use the new postgres config structure instead.`},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase := testCase
			logger := newLogger()
			configWithFlags(t, testCase.flags, config.WithLogger(logger))
			assert.Equal(t, testCase.expectedMessages, logger.messages)
		})
	}
}
