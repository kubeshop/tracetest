package transaction

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestDeleteTransaction(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	// setup isolated e2e environment
	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to delete a transaction that don't exist
	// Then it should return an error and say that this resource does not exist
	result := tracetestcli.Exec(t, "delete transaction --id dont-exist", tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 1)
	require.Contains(result.StdErr, "Resource transaction with ID dont-exist not found") // TODO: update this message to singular

	// When I try to set up a new transaction
	// Then it should be applied with success
	newTransactionPath := env.GetTestResourcePath(t, "new-transaction")

	result = tracetestcli.Exec(t, fmt.Sprintf("apply transaction --file %s", newTransactionPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	// When I try to delete the transaction
	// Then it should delete with success
	result = tracetestcli.Exec(t, "delete transaction --id Qti5R3_VR", tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)
	require.Contains(result.StdOut, "✔ Transaction successfully deleted")

	// When I try to get a transaction again
	// Then it should return a message saying that the transaction was not found
	result = tracetestcli.Exec(t, "delete transaction --id Qti5R3_VR", tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 1)
	require.Contains(result.StdErr, "Resource transaction with ID Qti5R3_VR not found")
}
