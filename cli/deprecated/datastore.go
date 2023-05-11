package deprecated

import (
	"github.com/kubeshop/tracetest/cli/global"
	global_middlewares "github.com/kubeshop/tracetest/cli/global/middlewares"
	"github.com/kubeshop/tracetest/cli/resources"
	"github.com/spf13/cobra"
)

var (
	// apply param
	dataStoreApplyFile string
	// export param
	exportOutputFile string
	dataStoreID      string
)

type DatastoreLegacy []*cobra.Command

func NewDatastoreLegacy(root global.Root) DatastoreLegacy {
	postRun := global_middlewares.ComposeNoopRun(global_middlewares.WithTeardownMiddleware(root.Setup))
	run := global_middlewares.ComposeRun(global_middlewares.WithAnalytics("datastore", "cmd"), global_middlewares.WithResultHandler(root.Setup))

	dataStoreCmd := &cobra.Command{
		GroupID:    global.GroupConfig.ID,
		Use:        "datastore",
		Short:      "Manage your tracetest data stores",
		Long:       "Manage your tracetest data stores",
		Deprecated: "Please use `tracetest (apply|delete|export|get) datastore` commands instead.",
		PreRun:     run(root.Setup.PreRun),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		PostRun: postRun(),
	}

	dataStoreApplyCmd := &cobra.Command{
		Use:        "apply",
		Short:      "Apply (create/update) data store configuration to your Tracetest server",
		Long:       "Apply (create/update) data store configuration to your Tracetest server",
		Deprecated: "Please use `tracetest apply datastore --file [path]` command instead.",
		PreRun:     run(root.Setup.PreRun),
		Run: func(cmd *cobra.Command, args []string) {

			// call new apply command
			apply := resources.NewApply(root)
			apply.Parameters.DefinitionFile = dataStoreApplyFile

			apply.Cmd.Run(apply.Cmd, []string{"datastore"})
		},
		PostRun: postRun(),
	}

	dataStoreApplyCmd.Flags().StringVarP(&dataStoreApplyFile, "file", "f", "", "Data store definition file")

	dataStoreExportCmd := &cobra.Command{
		Use:        "export",
		Short:      "Exports a data store configuration into a file",
		Long:       "Exports a data store configuration into a file",
		Deprecated: "Please use `tracetest export datastore --id [id]` command instead.",
		PreRun:     run(root.Setup.PreRun),
		Run: func(cmd *cobra.Command, args []string) {
			// call new export command
			export := resources.NewExport(root)

			export.Parameters.ResourceFile = exportOutputFile
			export.Parameters.ResourceID = "current"
			export.Cmd.Run(export.Cmd, []string{"datastore"})
		},
		PostRun: postRun(),
	}

	dataStoreListCmd := &cobra.Command{
		Use:        "list",
		Short:      "List data store configurations to your tracetest server",
		Long:       "List data store configurations to your tracetest server",
		Deprecated: "Please use `tracetest get datastore --id current` command instead.",
		PreRun:     run(root.Setup.PreRun),
		Run: func(cmd *cobra.Command, args []string) {
			// call new get command
			get := resources.NewGet(root)

			get.Parameters.ResourceID = "current"
			get.Cmd.Run(get.Cmd, []string{"datastore"})
		},
		PostRun: postRun(),
	}

	dataStoreCmd.AddCommand(dataStoreApplyCmd)
	dataStoreCmd.AddCommand(dataStoreExportCmd)

	dataStoreExportCmd.Flags().StringVarP(&exportOutputFile, "output", "o", "", "Output file")
	dataStoreExportCmd.Flags().StringVarP(&dataStoreID, "id", "i", "", "Data store ID")

	dataStoreApplyCmd.PersistentFlags().StringVarP(&dataStoreApplyFile, "file", "f", "", "file containing the data store configuration")
	dataStoreCmd.AddCommand(dataStoreApplyCmd)

	// export
	dataStoreExportCmd.PersistentFlags().StringVarP(&exportOutputFile, "output", "o", "", "file where data store configuration will be saved")
	dataStoreExportCmd.PersistentFlags().StringVarP(&dataStoreID, "id", "", "", "id of the data store that will be exported")
	dataStoreCmd.AddCommand(dataStoreExportCmd)

	// list
	dataStoreCmd.AddCommand(dataStoreListCmd)

	return DatastoreLegacy{dataStoreCmd, dataStoreApplyCmd, dataStoreExportCmd, dataStoreListCmd}
}
