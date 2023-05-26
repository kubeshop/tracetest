package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/kubeshop/tracetest/cli-e2etest/command"
	"github.com/stretchr/testify/require"
)

func UnmarshalJSON[T any](t *testing.T, data string) T {
	var value T

	err := json.Unmarshal([]byte(data), &value)
	require.NoError(t, err)

	return value
}

func UnmarshalYAML[T any](t *testing.T, data string) T {
	var value T

	err := yaml.Unmarshal([]byte(data), &value)
	require.NoError(t, err)

	return value
}

func UnmarshalYAMLSequence[T any](t *testing.T, data string) []T {
	decoder := yaml.NewDecoder(bytes.NewBuffer([]byte(data)))

	result := []T{}

	for {
		var value T
		err := decoder.Decode(&value)

		if errors.Is(err, io.EOF) {
			break
		}

		require.NoError(t, err)

		result = append(result, value)
	}

	return result
}

func RequireExitCodeEqual(t *testing.T, result *command.ExecResult, expectedExitCode int) {
	require.Equal(
		t, expectedExitCode, result.ExitCode,
		"command %s finished with wrong exit code. Expected code: %d Obtained code: %d \nStdOut: %s \nStdErr: %s",
		result.CommandExecuted, expectedExitCode, result.ExitCode, result.StdOut, result.StdErr,
	)
}
