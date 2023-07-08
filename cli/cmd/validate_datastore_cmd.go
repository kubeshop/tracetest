package cmd

import (
	"os"

	"github.com/kubeshop/tracetest/cli/actions"

	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/spf13/cobra"
)

var (
	validateDataStoreParams = &validateDataStoreParameters{}
	validateCmd             *cobra.Command
	validateDatastoreCmd    *cobra.Command
)

func init() {
	validateCmd = &cobra.Command{
		GroupID: cmdGroupResources.ID,
		Use:     "validate",
		Short:   "Validate a resource file",
		Long:    "Validate a resource file",
		PreRun:  setupCommand(),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		PostRun: teardownCommand,
	}

	validateDatastoreCmd = &cobra.Command{
		Use:    "datastore",
		Short:  "Validate a DataStore resource file, checking if it is possible to connect on it",
		Long:   "Validate a DataStore resource file, checking if it is possible to connect on it",
		PreRun: setupCommand(),
		Run: WithResultHandler(WithParamsHandler(validateDataStoreParams)(func(cmd *cobra.Command, _ []string) (string, error) {
			client := utils.GetAPIClient(cliConfig)

			action := actions.NewValidateDataStoreAction(client)
			return action.ValidateDatastore(validateDataStoreParams.DefinitionFile)
		})),
		PostRun: teardownCommand,
	}

	validateDatastoreCmd.PersistentFlags().StringVarP(&validateDataStoreParams.DefinitionFile, "file", "f", "", "definition file with datastore data")

	rootCmd.AddCommand(validateCmd)
	validateCmd.AddCommand(validateDatastoreCmd)
}

type validateDataStoreParameters struct {
	DefinitionFile string
}

func (p validateDataStoreParameters) Validate(cmd *cobra.Command, args []string) []error {
	errors := make([]error, 0)

	if p.DefinitionFile == "" {
		errors = append(errors, paramError{
			Parameter: "file",
			Message:   "Definition file must be provided",
		})
	}

	_, err := os.Stat(p.DefinitionFile)
	if os.IsNotExist(err) {
		errors = append(errors, paramError{
			Parameter: "file",
			Message:   "Definition file does not exist",
		})
	}

	return errors
}
