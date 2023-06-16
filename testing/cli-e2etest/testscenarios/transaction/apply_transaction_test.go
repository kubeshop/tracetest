package transaction

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/cli-e2etest/testscenarios/types"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestApplyTransaction(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	// setup isolated e2e environment
	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to set up a new transaction
	// Then it should be applied with success
	newTransactionPath := env.GetTestResourcePath(t, "new-transaction")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply transaction --file %s", newTransactionPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	transaction := helpers.UnmarshalYAML[types.TransactionResource](t, result.StdOut)

	require.Equal("Transaction", transaction.Type)
	require.Equal("Qti5R3_VR", transaction.Spec.ID)
	require.Equal("New Transaction", transaction.Spec.Name)
	require.Equal("a transaction", transaction.Spec.Description)
	require.Len(transaction.Spec.Steps, 2)
	require.Equal("9wtAH2_Vg", transaction.Spec.Steps[0])
	require.Equal("ajksdkasjbd", transaction.Spec.Steps[1])

	// When I try to get the transaction applied on the last step
	// Then it should return it
	result = tracetestcli.Exec(t, "get transaction --id Qti5R3_VR --output yaml", tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	require.Equal("Transaction", transaction.Type)
	require.Equal("Qti5R3_VR", transaction.Spec.ID)
	require.Equal("New Transaction", transaction.Spec.Name)
	require.Equal("a transaction", transaction.Spec.Description)
	require.Len(transaction.Spec.Steps, 2)
	require.Equal("9wtAH2_Vg", transaction.Spec.Steps[0])
	require.Equal("ajksdkasjbd", transaction.Spec.Steps[1])

	// When I try to update the last transaction
	// Then it should be applied with success
	updatedNewTransactionPath := env.GetTestResourcePath(t, "updated-new-transaction")

	result = tracetestcli.Exec(t, fmt.Sprintf("apply transaction --file %s", updatedNewTransactionPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	updatedTransaction := helpers.UnmarshalYAML[types.TransactionResource](t, result.StdOut)
	require.Equal("Transaction", updatedTransaction.Type)
	require.Equal("Qti5R3_VR", updatedTransaction.Spec.ID)
	require.Equal("Updated Transaction", updatedTransaction.Spec.Name)
	require.Equal("an updated transaction", updatedTransaction.Spec.Description)
	require.Len(updatedTransaction.Spec.Steps, 3)
	require.Equal("9wtAH2_Vg", updatedTransaction.Spec.Steps[0])
	require.Equal("ajksdkasjbd", updatedTransaction.Spec.Steps[1])
	require.Equal("ajksdkasjbd", updatedTransaction.Spec.Steps[2])

	// When I try to get the transaction applied on the last step
	// Then it should return it
	result = tracetestcli.Exec(t, "get transaction --id Qti5R3_VR --output yaml", tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	updatedTransaction = helpers.UnmarshalYAML[types.TransactionResource](t, result.StdOut)
	require.Equal("Transaction", updatedTransaction.Type)
	require.Equal("Qti5R3_VR", updatedTransaction.Spec.ID)
	require.Equal("Updated Transaction", updatedTransaction.Spec.Name)
	require.Equal("an updated transaction", updatedTransaction.Spec.Description)
	require.Len(updatedTransaction.Spec.Steps, 3)
	require.Equal("9wtAH2_Vg", updatedTransaction.Spec.Steps[0])
	require.Equal("ajksdkasjbd", updatedTransaction.Spec.Steps[1])
	require.Equal("ajksdkasjbd", updatedTransaction.Spec.Steps[2])
}
