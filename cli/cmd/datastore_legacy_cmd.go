package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// apply param
	dataStoreApplyFile string
	// export param
	exportOutputFile string
	dataStoreID      string
)

var dataStoreCmd = &cobra.Command{
	GroupID:    cmdGroupConfig.ID,
	Use:        "datastore",
	Short:      "Manage your tracetest data stores",
	Long:       "Manage your tracetest data stores",
	Deprecated: "Please use `tracetest (apply|delete|export|get) datastore` commands instead.",
	PreRun:     setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	PostRun: teardownCommand,
}

var dataStoreApplyCmd = &cobra.Command{
	Use:        "apply",
	Short:      "Apply (create/update) data store configuration to your Tracetest server",
	Long:       "Apply (create/update) data store configuration to your Tracetest server",
	Deprecated: "Please use `tracetest apply datastore --file [path]` command instead.",
	PreRun:     setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		// call new apply command
		definitionFile = dataStoreApplyFile
		applyCmd.Run(applyCmd, []string{"datastore"})
	},
	PostRun: teardownCommand,
}

var dataStoreExportCmd = &cobra.Command{
	Use:        "export",
	Short:      "Exports a data store configuration into a file",
	Long:       "Exports a data store configuration into a file",
	Deprecated: "Please use `tracetest export datastore --id [id]` command instead.",
	PreRun:     setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		// call new export command
		exportResourceID = "current"
		exportResourceFile = exportOutputFile
		exportCmd.Run(exportCmd, []string{"datastore"})
	},
	PostRun: teardownCommand,
}

var dataStoreListCmd = &cobra.Command{
	Use:        "list",
	Short:      "List data store configurations to your tracetest server",
	Long:       "List data store configurations to your tracetest server",
	Deprecated: "Please use `tracetest get datastore --id current` command instead.",
	PreRun:     setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		// call new get command
		resourceID = "current"
		getCmd.Run(getCmd, []string{"datastore"})
	},
	PostRun: teardownCommand,
}

func init() {
	rootCmd.AddCommand(dataStoreCmd)

	// apply
	dataStoreApplyCmd.PersistentFlags().StringVarP(&dataStoreApplyFile, "file", "f", "", "file containing the data store configuration")
	dataStoreCmd.AddCommand(dataStoreApplyCmd)

	// export
	dataStoreExportCmd.PersistentFlags().StringVarP(&exportOutputFile, "output", "o", "", "file where data store configuration will be saved")
	dataStoreExportCmd.PersistentFlags().StringVarP(&dataStoreID, "id", "", "", "id of the data store that will be exported")
	dataStoreCmd.AddCommand(dataStoreExportCmd)

	// list
	dataStoreCmd.AddCommand(dataStoreListCmd)
}
