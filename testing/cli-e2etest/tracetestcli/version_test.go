package tracetestcli_test

import (
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestVersionCommand(t *testing.T) {
	result, err := tracetestcli.Exec("version")
	require.NoError(t, err)
	require.Equal(t, 0, result.ExitCode)
}
