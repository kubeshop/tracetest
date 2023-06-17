package transaction

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/cli-e2etest/testscenarios/types"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func addGetTransactionPreReqs(t *testing.T, env environment.Manager) {
	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to set up a new transaction
	// Then it should be applied with success
	newTransactionPath := env.GetTestResourcePath(t, "new-transaction")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply transaction --file %s", newTransactionPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)
}

func TestGetTransaction(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	t.Run("get with no transaction initialized", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And no transaction registered

		// When I try to get a transaction on yaml mode
		// Then it should return a error message
		result := tracetestcli.Exec(t, "get transaction --id no-id --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "Resource transaction with ID no-id not found")
	})

	addGetTransactionPreReqs(t, env)

	t.Run("get with YAML format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a transaction already set

		// When I try to get a transaction on yaml mode
		// Then it should print a YAML
		result := tracetestcli.Exec(t, "get transaction --id Qti5R3_VR --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		transaction := helpers.UnmarshalYAML[types.TransactionResource](t, result.StdOut)

		require.Equal("Transaction", transaction.Type)
		require.Equal("Qti5R3_VR", transaction.Spec.ID)
		require.Equal("New Transaction", transaction.Spec.Name)
		require.Equal("a transaction", transaction.Spec.Description)
		require.Len(transaction.Spec.Steps, 2)
		require.Equal("9wtAH2_Vg", transaction.Spec.Steps[0])
		require.Equal("ajksdkasjbd", transaction.Spec.Steps[1])
	})

	t.Run("get with JSON format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a transaction already set

		// When I try to get a transaction on json mode
		// Then it should print a json
		result := tracetestcli.Exec(t, "get transaction --id Qti5R3_VR --output json", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		transaction := helpers.UnmarshalJSON[types.TransactionResource](t, result.StdOut)

		require.Equal("Transaction", transaction.Type)
		require.Equal("Qti5R3_VR", transaction.Spec.ID)
		require.Equal("New Transaction", transaction.Spec.Name)
		require.Equal("a transaction", transaction.Spec.Description)
		require.Len(transaction.Spec.Steps, 2)
		require.Equal("9wtAH2_Vg", transaction.Spec.Steps[0])
		require.Equal("ajksdkasjbd", transaction.Spec.Steps[1])
	})

	t.Run("get with pretty format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a transaction already set

		// When I try to get a transaction on pretty mode
		// Then it should print a table with 4 lines printed: header, separator, transaction item and empty line
		result := tracetestcli.Exec(t, "get transaction --id Qti5R3_VR --output pretty", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "New Transaction")

		lines := strings.Split(result.StdOut, "\n")
		require.Len(lines, 4)
	})
}
