package config_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	configFromFlagsWithLogger := func(logger *logger, inputFlags []string) *config.Config {
		flags := pflag.NewFlagSet("fake", pflag.ExitOnError)
		config.SetupFlags(flags)
		err := flags.Parse(inputFlags)
		require.NoError(t, err)
		cfg, err := config.New(flags, logger)
		require.NoError(t, err)
		return cfg
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase := testCase
			logger := newLogger()
			configFromFlagsWithLogger(logger, testCase.flags)
			assert.Equal(t, testCase.expectedMessages, logger.messages)
		})
	}
}
