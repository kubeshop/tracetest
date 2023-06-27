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

func addListTransactionPreReqs(t *testing.T, env environment.Manager) {
	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to set up a new environment
	// Then it should be applied with success
	newTransactionPath := env.GetTestResourcePath(t, "new-transaction")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply transaction --file %s", newTransactionPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	// When I try to set up a another environment
	// Then it should be applied with success
	anotherTransactionPath := env.GetTestResourcePath(t, "another-transaction")

	result = tracetestcli.Exec(t, fmt.Sprintf("apply transaction --file %s", anotherTransactionPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	// When I try to set up a third environment
	// Then it should be applied with success
	oneMoreTransactionPath := env.GetTestResourcePath(t, "one-more-transaction")

	result = tracetestcli.Exec(t, fmt.Sprintf("apply transaction --file %s", oneMoreTransactionPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)
}

func TestListTransactions(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	// setup isolated e2e environment
	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	t.Run("list no transactions", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And there is no envs
		result := tracetestcli.Exec(t, "list transaction --sortBy name --sortDirection asc --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		transactions := helpers.UnmarshalYAMLSequence[types.AugmentedTransactionResource](t, result.StdOut)
		require.Len(transactions, 0)
	})

	addListTransactionPreReqs(t, env)

	t.Run("list with invalid sortBy field", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list these transactions by an invalid field
		// Then I should receive an error
		result := tracetestcli.Exec(t, "list transaction --sortBy id --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 1)
		require.Contains(result.StdErr, "invalid sort field: id") // TODO: think on how to improve this error handling
	})

	t.Run("list with YAML format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list these transactions by a valid field and in YAML format
		// Then I should receive three transactions
		result := tracetestcli.Exec(t, "list transaction --sortBy name --sortDirection asc --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		transactions := helpers.UnmarshalYAMLSequence[types.AugmentedTransactionResource](t, result.StdOut)
		require.Len(transactions, 3)

		anotherTransaction := transactions[0]
		require.Equal("Transaction", anotherTransaction.Type)
		require.Equal("asuhfdkj", anotherTransaction.Spec.ID)
		require.Equal("Another Transaction", anotherTransaction.Spec.Name)
		require.Equal("another transaction", anotherTransaction.Spec.Description)
		require.Equal(0, anotherTransaction.Spec.Summary.Runs)
		require.Equal(0, anotherTransaction.Spec.Summary.LastRun.Fails)
		require.Equal(0, anotherTransaction.Spec.Summary.LastRun.Passes)
		require.Len(anotherTransaction.Spec.Steps, 4)
		require.Equal("9wtAH2_Vg", anotherTransaction.Spec.Steps[0])
		require.Equal("9wtAH2_Vg", anotherTransaction.Spec.Steps[1])
		require.Equal("ajksdkasjbd", anotherTransaction.Spec.Steps[2])
		require.Equal("ajksdkasjbd", anotherTransaction.Spec.Steps[3])

		newTransaction := transactions[1]
		require.Equal("Transaction", newTransaction.Type)
		require.Equal("Qti5R3_VR", newTransaction.Spec.ID)
		require.Equal("New Transaction", newTransaction.Spec.Name)
		require.Equal("a transaction", newTransaction.Spec.Description)
		require.Equal(0, newTransaction.Spec.Summary.Runs)
		require.Equal(0, newTransaction.Spec.Summary.LastRun.Fails)
		require.Equal(0, newTransaction.Spec.Summary.LastRun.Passes)
		require.Len(newTransaction.Spec.Steps, 2)
		require.Equal("9wtAH2_Vg", newTransaction.Spec.Steps[0])
		require.Equal("ajksdkasjbd", newTransaction.Spec.Steps[1])

		oneMoreTransaction := transactions[2]
		require.Equal("Transaction", oneMoreTransaction.Type)
		require.Equal("i2ug34j", oneMoreTransaction.Spec.ID)
		require.Equal("One More Transaction", oneMoreTransaction.Spec.Name)
		require.Equal("one more transaction", oneMoreTransaction.Spec.Description)
		require.Equal(0, oneMoreTransaction.Spec.Summary.Runs)
		require.Equal(0, oneMoreTransaction.Spec.Summary.LastRun.Fails)
		require.Equal(0, oneMoreTransaction.Spec.Summary.LastRun.Passes)
		require.Len(oneMoreTransaction.Spec.Steps, 3)
		require.Equal("9wtAH2_Vg", oneMoreTransaction.Spec.Steps[0])
		require.Equal("9wtAH2_Vg", oneMoreTransaction.Spec.Steps[1])
		require.Equal("ajksdkasjbd", oneMoreTransaction.Spec.Steps[2])
	})

	t.Run("list with JSON format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list these transactions by a valid field and in JSON format
		// Then I should receive three transactions
		result := tracetestcli.Exec(t, "list transaction --sortBy name --sortDirection asc --output json", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		transactions := helpers.UnmarshalJSON[types.ResourceList[types.AugmentedTransactionResource]](t, result.StdOut)
		require.Len(transactions.Items, 3)
		require.Equal(len(transactions.Items), transactions.Count)

		anotherTransaction := transactions.Items[0]
		require.Equal("Transaction", anotherTransaction.Type)
		require.Equal("asuhfdkj", anotherTransaction.Spec.ID)
		require.Equal("Another Transaction", anotherTransaction.Spec.Name)
		require.Equal("another transaction", anotherTransaction.Spec.Description)
		require.Equal(0, anotherTransaction.Spec.Summary.Runs)
		require.Equal(0, anotherTransaction.Spec.Summary.LastRun.Fails)
		require.Equal(0, anotherTransaction.Spec.Summary.LastRun.Passes)
		require.Len(anotherTransaction.Spec.Steps, 4)
		require.Equal("9wtAH2_Vg", anotherTransaction.Spec.Steps[0])
		require.Equal("9wtAH2_Vg", anotherTransaction.Spec.Steps[1])
		require.Equal("ajksdkasjbd", anotherTransaction.Spec.Steps[2])
		require.Equal("ajksdkasjbd", anotherTransaction.Spec.Steps[3])

		newTransaction := transactions.Items[1]
		require.Equal("Transaction", newTransaction.Type)
		require.Equal("Qti5R3_VR", newTransaction.Spec.ID)
		require.Equal("New Transaction", newTransaction.Spec.Name)
		require.Equal("a transaction", newTransaction.Spec.Description)
		require.Equal(0, newTransaction.Spec.Summary.Runs)
		require.Equal(0, newTransaction.Spec.Summary.LastRun.Fails)
		require.Equal(0, newTransaction.Spec.Summary.LastRun.Passes)
		require.Len(newTransaction.Spec.Steps, 2)
		require.Equal("9wtAH2_Vg", newTransaction.Spec.Steps[0])
		require.Equal("ajksdkasjbd", newTransaction.Spec.Steps[1])

		oneMoreTransaction := transactions.Items[2]
		require.Equal("Transaction", oneMoreTransaction.Type)
		require.Equal("i2ug34j", oneMoreTransaction.Spec.ID)
		require.Equal("One More Transaction", oneMoreTransaction.Spec.Name)
		require.Equal("one more transaction", oneMoreTransaction.Spec.Description)
		require.Equal(0, oneMoreTransaction.Spec.Summary.Runs)
		require.Equal(0, oneMoreTransaction.Spec.Summary.LastRun.Fails)
		require.Equal(0, oneMoreTransaction.Spec.Summary.LastRun.Passes)
		require.Len(oneMoreTransaction.Spec.Steps, 3)
		require.Equal("9wtAH2_Vg", oneMoreTransaction.Spec.Steps[0])
		require.Equal("9wtAH2_Vg", oneMoreTransaction.Spec.Steps[1])
		require.Equal("ajksdkasjbd", oneMoreTransaction.Spec.Steps[2])
	})

	t.Run("list with pretty format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list these transactions by a valid field and in pretty format
		// Then it should print a table with 6 lines printed: header, separator, three transactions and empty line
		result := tracetestcli.Exec(t, "list transaction --sortBy name --sortDirection asc --output pretty", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		parsedTable := helpers.UnmarshalTable(t, result.StdOut)
		require.Len(parsedTable, 3)

		firstLine := parsedTable[0]
		require.Equal("asuhfdkj", firstLine["ID"])
		require.Equal("Another Transaction", firstLine["NAME"])
		require.Equal("1", firstLine["VERSION"])
		require.Equal("4", firstLine["STEPS"])
		require.Equal("0", firstLine["RUNS"])
		require.Equal("", firstLine["LAST RUN TIME"])
		require.Equal("0", firstLine["LAST RUN SUCCESSES"])
		require.Equal("0", firstLine["LAST RUN FAILURES"])

		secondLine := parsedTable[1]
		require.Equal("Qti5R3_VR", secondLine["ID"])
		require.Equal("New Transaction", secondLine["NAME"])
		require.Equal("1", secondLine["VERSION"])
		require.Equal("2", secondLine["STEPS"])
		require.Equal("0", secondLine["RUNS"])
		require.Equal("", secondLine["LAST RUN TIME"])
		require.Equal("0", secondLine["LAST RUN SUCCESSES"])
		require.Equal("0", secondLine["LAST RUN FAILURES"])

		thirdLine := parsedTable[2]
		require.Equal("i2ug34j", thirdLine["ID"])
		require.Equal("One More Transaction", thirdLine["NAME"])
		require.Equal("1", thirdLine["VERSION"])
		require.Equal("3", thirdLine["STEPS"])
		require.Equal("0", thirdLine["RUNS"])
		require.Equal("", thirdLine["LAST RUN TIME"])
		require.Equal("0", thirdLine["LAST RUN SUCCESSES"])
		require.Equal("0", thirdLine["LAST RUN FAILURES"])
	})

	t.Run("list with YAML format skipping the first and taking two items", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list these transactions by a valid field, paging options and in YAML format
		// Then I should receive two transactions
		result := tracetestcli.Exec(t, "list transaction --sortBy name --sortDirection asc --skip 1 --take 2 --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		transactions := helpers.UnmarshalYAMLSequence[types.AugmentedTransactionResource](t, result.StdOut)
		require.Len(transactions, 2)

		newTransaction := transactions[0]
		require.Equal("Transaction", newTransaction.Type)
		require.Equal("Qti5R3_VR", newTransaction.Spec.ID)
		require.Equal("New Transaction", newTransaction.Spec.Name)
		require.Equal("a transaction", newTransaction.Spec.Description)
		require.Equal(0, newTransaction.Spec.Summary.Runs)
		require.Equal(0, newTransaction.Spec.Summary.LastRun.Fails)
		require.Equal(0, newTransaction.Spec.Summary.LastRun.Passes)
		require.Len(newTransaction.Spec.Steps, 2)
		require.Equal("9wtAH2_Vg", newTransaction.Spec.Steps[0])
		require.Equal("ajksdkasjbd", newTransaction.Spec.Steps[1])

		oneMoreTransaction := transactions[1]
		require.Equal("Transaction", oneMoreTransaction.Type)
		require.Equal("i2ug34j", oneMoreTransaction.Spec.ID)
		require.Equal("One More Transaction", oneMoreTransaction.Spec.Name)
		require.Equal("one more transaction", oneMoreTransaction.Spec.Description)
		require.Equal(0, oneMoreTransaction.Spec.Summary.Runs)
		require.Equal(0, oneMoreTransaction.Spec.Summary.LastRun.Fails)
		require.Equal(0, oneMoreTransaction.Spec.Summary.LastRun.Passes)
		require.Len(oneMoreTransaction.Spec.Steps, 3)
		require.Equal("9wtAH2_Vg", oneMoreTransaction.Spec.Steps[0])
		require.Equal("9wtAH2_Vg", oneMoreTransaction.Spec.Steps[1])
		require.Equal("ajksdkasjbd", oneMoreTransaction.Spec.Steps[2])
	})
}
