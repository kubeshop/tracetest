package cmd

import (
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
		Short:   "Validate a resource file, checking if it has valid data",
		Long:    "Validate a resource file, checking if it has valid data",
		PreRun:  setupCommand(),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		PostRun: teardownCommand,
	}

	validateDatastoreCmd = &cobra.Command{
		Use:    "datastore",
		Short:  "Validate a DataStore resource file, checking if it has valid data",
		Long:   "Validate a DataStore resource file, checking if it has valid data",
		PreRun: setupCommand(),
		Run: WithResultHandler(func(_ *cobra.Command, _ []string) (string, error) {
			client := utils.GetAPIClient(cliConfig)

			action := actions.NewValidateDataStoreAction(client)
			return action.ValidateDatastore(validateDataStoreParams.DefinitionFile)
		}),
		PostRun: teardownCommand,
	}

	validateDatastoreCmd.PersistentFlags().StringVarP(&validateDataStoreParams.DefinitionFile, "file", "f", "", "id of the environment to be used")

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

	return errors
}
