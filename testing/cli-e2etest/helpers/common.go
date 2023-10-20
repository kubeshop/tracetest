package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/command"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/testscenarios/types"
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

func InjectIdIntoDemoFile(t *testing.T, filePath, id string) {
	// This is a workaround method used to deal with the upsert restriction for demo resources
	// for more details look into: https://github.com/kubeshop/tracetest/issues/2719

	fileContent, err := os.ReadFile(filePath)
	require.NoError(t, err)

	demo := UnmarshalYAML[types.DemoResource](t, string(fileContent))
	demo.Spec.Id = id

	newFileContent, err := yaml.Marshal(demo)
	require.NoError(t, err)

	err = os.WriteFile(filePath, newFileContent, os.ModeAppend)
	require.NoError(t, err)
}

func Copy(source, dst string) {
	os.Remove(dst)
	sourceFile, err := os.Open(source)
	if err != nil {
		panic(err)
	}
	defer sourceFile.Close()

	newFile, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, sourceFile)
	if err != nil {
		panic(err)
	}
}

func RemoveIDFromTestSuiteFile(t *testing.T, filePath string) {
	fileContent, err := os.ReadFile(filePath)
	require.NoError(t, err)

	suite := UnmarshalYAML[types.TestSuiteResource](t, string(fileContent))
	suite.Spec.ID = ""

	newFileContent, err := yaml.Marshal(suite)
	require.NoError(t, err)

	err = os.WriteFile(filePath, newFileContent, os.ModeAppend)
	require.NoError(t, err)
}
